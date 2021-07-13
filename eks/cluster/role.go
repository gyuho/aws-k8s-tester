package cluster

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	aws_iam "github.com/aws/aws-k8s-tester/pkg/aws/iam"
	aws_v2 "github.com/aws/aws-sdk-go-v2/aws"
	aws_iam_v2 "github.com/aws/aws-sdk-go-v2/service/iam"
	smithy "github.com/aws/smithy-go"
	"go.uber.org/zap"
)

// see https://github.com/aws/aws-k8s-tester/blob/v1.6.0/eks/cluster/role.go for CloudFormation based workflow

func (ts *tester) createRole() error {
	fmt.Print(ts.cfg.EKSConfig.Colorize("\n\n[yellow]*********************************\n"))
	fmt.Printf(ts.cfg.EKSConfig.Colorize("[light_green]createRole [default](%q)\n"), ts.cfg.EKSConfig.ConfigPath)

	if !ts.cfg.EKSConfig.Role.Create {
		ts.cfg.Logger.Info("Role.Create false; skipping creation")
		return aws_iam.ValidateV2(
			ts.cfg.Logger,
			ts.cfg.IAMAPIV2,
			ts.cfg.EKSConfig.Role.Name,
			[]string{"eks.amazonaws.com"},
			[]string{
				// Prior to April 16, 2020, AmazonEKSServicePolicy was also required and the suggested name was eksServiceRole. With the AWSServiceRoleForAmazonEKS service-linked role, that policy is no longer required for clusters created on or after April 16, 2020.
				// ref. https://docs.aws.amazon.com/eks/latest/userguide/service_IAM_role.html
				"arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
			},
		)
	}
	if ts.cfg.EKSConfig.Role.ARN != "" {
		ts.cfg.Logger.Info("role already created; no need to create a new one")
		return nil
	}
	if ts.cfg.EKSConfig.Role.Name == "" {
		return errors.New("cannot create a cluster role with an empty Role.Name")
	}

	if err := ts._createRole(); err != nil {
		return err
	}
	if err := ts.createPolicy(); err != nil {
		return err
	}
	if err := ts.attachPolicy(); err != nil {
		return err
	}

	ts.cfg.Logger.Info("created a new role and attached policy",
		zap.String("cluster-role-arn", ts.cfg.EKSConfig.Role.ARN),
		zap.String("cluster-role-name", ts.cfg.EKSConfig.Role.Name),
	)
	return nil
}

func (ts *tester) deleteRole() error {
	fmt.Print(ts.cfg.EKSConfig.Colorize("\n\n[yellow]*********************************\n"))
	fmt.Printf(ts.cfg.EKSConfig.Colorize("[light_blue]deleteRole [default](%q)\n"), ts.cfg.EKSConfig.ConfigPath)

	if !ts.cfg.EKSConfig.Role.Create {
		ts.cfg.Logger.Info("Role.Create false; skipping deletion")
		return nil
	}

	var errs []string
	if err := ts.detachPolicy(); err != nil {
		errs = append(errs, err.Error())
	}
	if err := ts.deletePolicy(); err != nil {
		errs = append(errs, err.Error())
	}
	if err := ts._deleteRole(); err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) == 0 {
		ts.cfg.Logger.Info("deleted a role",
			zap.String("cluster-role-arn", ts.cfg.EKSConfig.Role.ARN),
			zap.String("cluster-role-name", ts.cfg.EKSConfig.Role.Name),
		)
		return nil
	}
	return errors.New(strings.Join(errs, ","))
}

func (ts *tester) _createRole() error {
	ts.cfg.Logger.Info("creating role", zap.String("name", ts.cfg.EKSConfig.Role.Name))
	out, err := ts.cfg.IAMAPIV2.CreateRole(
		context.Background(),
		&aws_iam_v2.CreateRoleInput{
			RoleName:                 aws_v2.String(ts.cfg.EKSConfig.Role.Name),
			Path:                     aws_v2.String("/"),
			AssumeRolePolicyDocument: aws_v2.String(createAssumeRolePolicyDocument(ts.cfg.EKSConfig.Role.ServicePrincipals)),
		},
	)
	if err != nil {
		return err
	}

	ts.cfg.Logger.Info("created role")
	ts.cfg.EKSConfig.Role.ARN = aws_v2.ToString(out.Role.Arn)
	ts.cfg.EKSConfig.Sync()
	return nil
}

func (ts *tester) _deleteRole() error {
	ts.cfg.Logger.Info("deleting role", zap.String("name", ts.cfg.EKSConfig.Role.Name))
	if _, ok := ts.cfg.EKSConfig.Status.DeletedResources[ts.cfg.EKSConfig.Role.Name]; ok {
		return nil
	}

	_, err := ts.cfg.IAMAPIV2.DeleteRole(
		context.Background(),
		&aws_iam_v2.DeleteRoleInput{
			RoleName: aws_v2.String(ts.cfg.EKSConfig.Role.Name),
		},
	)
	if err != nil {
		ts.cfg.Logger.Warn("failed to delete cluster role", zap.Error(err))
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if strings.Contains(apiErr.ErrorCode(), "NotFound") {
				ts.cfg.EKSConfig.Status.DeletedResources[ts.cfg.EKSConfig.Role.Name] = "Role.Name"
				ts.cfg.EKSConfig.Sync()
				return nil
			}
		}
		return err
	}

	ts.cfg.Logger.Info("deleted role")
	ts.cfg.EKSConfig.Status.DeletedResources[ts.cfg.EKSConfig.Role.Name] = "Role.Name"
	ts.cfg.EKSConfig.Sync()
	return nil
}

// https://docs.aws.amazon.com/eks/latest/userguide/ebs-csi.html
// https://github.com/kubernetes-sigs/aws-alb-ingress-controller/blob/master/docs/examples/iam-policy.json
// https://github.com/aws/eks-charts/tree/master/stable/appmesh-controller
func (ts *tester) createPolicy() error {
	ts.cfg.Logger.Info("creating policy", zap.String("name", ts.cfg.EKSConfig.Role.PolicyName))
	pout, err := ts.cfg.IAMAPIV2.CreatePolicy(
		context.Background(),
		&aws_iam_v2.CreatePolicyInput{
			PolicyName:     aws_v2.String(ts.cfg.EKSConfig.Role.PolicyName),
			PolicyDocument: aws_v2.String(createRolePolicyDocument()),
		},
	)
	if err != nil {
		return err
	}

	ts.cfg.Logger.Info("created policy")
	ts.cfg.EKSConfig.Role.PolicyARN = aws_v2.ToString(pout.Policy.Arn)
	ts.cfg.EKSConfig.Sync()
	return nil
}

func (ts *tester) deletePolicy() error {
	ts.cfg.Logger.Info("deleting policy")
	if _, ok := ts.cfg.EKSConfig.Status.DeletedResources[ts.cfg.EKSConfig.Role.PolicyARN]; ok {
		return nil
	}

	_, err := ts.cfg.IAMAPIV2.DeletePolicy(
		context.Background(),
		&aws_iam_v2.DeletePolicyInput{
			PolicyArn: aws_v2.String(ts.cfg.EKSConfig.Role.PolicyARN),
		},
	)
	if err != nil {
		ts.cfg.Logger.Warn("failed to delete policy", zap.Error(err))
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if strings.Contains(apiErr.ErrorCode(), "NotFound") {
				ts.cfg.EKSConfig.Status.DeletedResources[ts.cfg.EKSConfig.Role.PolicyARN] = "Role.PolicyARN"
				ts.cfg.EKSConfig.Sync()
				return nil
			}
		}
		return err
	}

	ts.cfg.Logger.Info("deleted policy")
	ts.cfg.EKSConfig.Status.DeletedResources[ts.cfg.EKSConfig.Role.PolicyARN] = "Role.PolicyARN"
	ts.cfg.EKSConfig.Sync()
	return nil
}

func (ts *tester) attachPolicy() error {
	ts.cfg.Logger.Info("attaching policy")
	_, err := ts.cfg.IAMAPIV2.AttachRolePolicy(
		context.Background(),
		&aws_iam_v2.AttachRolePolicyInput{
			RoleName:  aws_v2.String(ts.cfg.EKSConfig.Role.Name),
			PolicyArn: aws_v2.String(ts.cfg.EKSConfig.Role.PolicyARN),
		},
	)
	if err != nil {
		ts.cfg.Logger.Warn("failed to attach policy", zap.Error(err))
		return err
	}
	ts.cfg.Logger.Info("attached policy")
	return nil
}

func (ts *tester) detachPolicy() error {
	ts.cfg.Logger.Info("detaching policy")
	_, err := ts.cfg.IAMAPIV2.DetachRolePolicy(
		context.Background(),
		&aws_iam_v2.DetachRolePolicyInput{
			RoleName:  aws_v2.String(ts.cfg.EKSConfig.Role.Name),
			PolicyArn: aws_v2.String(ts.cfg.EKSConfig.Role.PolicyARN),
		},
	)
	if err != nil {
		ts.cfg.Logger.Warn("failed to detach policy", zap.Error(err))
		return err
	}
	ts.cfg.Logger.Info("detached policy")
	return nil
}

func createAssumeRolePolicyDocument(sps []string) string {
	p := aws_iam.PolicyDocument{
		Version:   "2012-10-17",
		Statement: createStatementEntriesForAssumeRole(sps),
	}
	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func createRolePolicyDocument() string {
	p := aws_iam.PolicyDocument{
		Version:   "2012-10-17",
		Statement: createStatementEntriesForRolePolicyDocument(),
	}
	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func createStatementEntriesForAssumeRole(sps []string) []aws_iam.StatementEntry {
	return []aws_iam.StatementEntry{
		{
			Effect: "Allow",
			Principal: &aws_iam.PrincipalEntry{
				Service: sps,
			},
			Action: []string{
				"sts:AssumeRole",
			},
		},
	}
}

// TODO: update based on add-on setups
// https://docs.aws.amazon.com/eks/latest/userguide/ebs-csi.html
// https://github.com/kubernetes-sigs/aws-alb-ingress-controller/blob/master/docs/examples/iam-policy.json
// https://github.com/aws/eks-charts/tree/master/stable/appmesh-controller
func createStatementEntriesForRolePolicyDocument() []aws_iam.StatementEntry {
	return []aws_iam.StatementEntry{
		{
			Effect:   "Allow",
			Resource: "*",
			Action: []string{
				"acm:DescribeCertificate",
				"acm:ListCertificates",
				"acm:GetCertificate",
			},
		},
		{
			Effect:   "Allow",
			Resource: "*",
			Action: []string{
				"ec2:AuthorizeSecurityGroupIngress",
				"ec2:CreateSecurityGroup",
				"ec2:CreateTags",
				"ec2:DeleteTags",
				"ec2:DeleteSecurityGroup",
				"ec2:DescribeAccountAttributes",
				"ec2:DescribeAddresses",
				"ec2:DescribeInstances",
				"ec2:DescribeInstanceStatus",
				"ec2:DescribeInternetGateways",
				"ec2:DescribeNetworkInterfaces",
				"ec2:DescribeSecurityGroups",
				"ec2:DescribeSubnets",
				"ec2:DescribeTags",
				"ec2:DescribeVpcs",
				"ec2:ModifyInstanceAttribute",
				"ec2:ModifyNetworkInterfaceAttribute",
				"ec2:RevokeSecurityGroupIngress",
				"ec2:AttachVolume",
				"ec2:CreateSnapshot",
				"ec2:CreateVolume",
				"ec2:DeleteSnapshot",
				"ec2:DeleteVolume",
				"ec2:DescribeSnapshots",
				"ec2:DescribeVolumes",
				"ec2:DetachVolume",
			},
		},
		{
			Effect:   "Allow",
			Resource: "*",
			Action: []string{
				"elasticloadbalancing:AddListenerCertificates",
				"elasticloadbalancing:AddTags",
				"elasticloadbalancing:CreateListener",
				"elasticloadbalancing:CreateLoadBalancer",
				"elasticloadbalancing:CreateRule",
				"elasticloadbalancing:CreateTargetGroup",
				"elasticloadbalancing:DeleteListener",
				"elasticloadbalancing:DeleteLoadBalancer",
				"elasticloadbalancing:DeleteRule",
				"elasticloadbalancing:DeleteTargetGroup",
				"elasticloadbalancing:DeregisterTargets",
				"elasticloadbalancing:DescribeListenerCertificates",
				"elasticloadbalancing:DescribeListeners",
				"elasticloadbalancing:DescribeLoadBalancers",
				"elasticloadbalancing:DescribeLoadBalancerAttributes",
				"elasticloadbalancing:DescribeRules",
				"elasticloadbalancing:DescribeSSLPolicies",
				"elasticloadbalancing:DescribeTags",
				"elasticloadbalancing:DescribeTargetGroups",
				"elasticloadbalancing:DescribeTargetGroupAttributes",
				"elasticloadbalancing:DescribeTargetHealth",
				"elasticloadbalancing:ModifyListener",
				"elasticloadbalancing:ModifyLoadBalancerAttributes",
				"elasticloadbalancing:ModifyRule",
				"elasticloadbalancing:ModifyTargetGroup",
				"elasticloadbalancing:ModifyTargetGroupAttributes",
				"elasticloadbalancing:RegisterTargets",
				"elasticloadbalancing:RemoveListenerCertificates",
				"elasticloadbalancing:RemoveTags",
				"elasticloadbalancing:SetIpAddressType",
				"elasticloadbalancing:SetSecurityGroups",
				"elasticloadbalancing:SetSubnets",
				"elasticloadbalancing:SetWebACL",
			},
		},
		{
			Effect:   "Allow",
			Resource: "*",
			Action: []string{
				"iam:CreateServiceLinkedRole",
				"iam:GetServerCertificate",
				"iam:ListServerCertificates",
			},
		},
		{
			Effect:   "Allow",
			Resource: "*",
			Action: []string{
				"cognito-idp:DescribeUserPoolClient",
			},
		},
		{
			Effect:   "Allow",
			Resource: "*",
			Action: []string{
				"waf-regional:GetWebACLForResource",
				"waf-regional:GetWebACL",
				"waf-regional:AssociateWebACL",
				"waf-regional:DisassociateWebACL",
			},
		},
		{
			Effect:   "Allow",
			Resource: "*",
			Action: []string{
				"tag:GetResources",
				"tag:TagResources",
			},
		},
		{
			Effect:   "Allow",
			Resource: "*",
			Action: []string{
				"waf:GetWebACL",
			},
		},
		{
			Effect:   "Allow",
			Resource: "*",
			Action: []string{
				"wafv2:GetWebACL",
				"wafv2:GetWebACLForResource",
				"wafv2:AssociateWebACL",
				"wafv2:DisassociateWebACL",
			},
		},
		{
			Effect:   "Allow",
			Resource: "*",
			Action: []string{
				"shield:DescribeProtection",
				"shield:GetSubscriptionState",
				"shield:DeleteProtection",
				"shield:CreateProtection",
				"shield:DescribeSubscription",
				"shield:ListProtections",
			},
		},
		{
			Effect:   "Allow",
			Resource: "*",
			Action: []string{
				"appmesh:*",
				"servicediscovery:CreateService",
				"servicediscovery:GetService",
				"servicediscovery:RegisterInstance",
				"servicediscovery:DeregisterInstance",
				"servicediscovery:ListInstances",
				"servicediscovery:ListNamespaces",
				"servicediscovery:ListServices",
				"route53:GetHealthCheck",
				"route53:CreateHealthCheck",
				"route53:UpdateHealthCheck",
				"route53:ChangeResourceRecordSets",
				"route53:DeleteHealthCheck",
			},
		},
	}
}
