// Code generated by smithy-go-codegen DO NOT EDIT.

package autoscaling

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Updates the configuration for the specified Auto Scaling group. To update an
// Auto Scaling group, specify the name of the group and the parameter that you
// want to change. Any parameters that you don't specify are not changed by this
// update request. The new settings take effect on any scaling activities after
// this call returns. If you associate a new launch configuration or template with
// an Auto Scaling group, all new instances will get the updated configuration.
// Existing instances continue to run with the configuration that they were
// originally launched with. When you update a group to specify a mixed instances
// policy instead of a launch configuration or template, existing instances may be
// replaced to match the new purchasing options that you specified in the policy.
// For example, if the group currently has 100% On-Demand capacity and the policy
// specifies 50% Spot capacity, this means that half of your instances will be
// gradually terminated and relaunched as Spot Instances. When replacing instances,
// Amazon EC2 Auto Scaling launches new instances before terminating the old ones,
// so that updating your group does not compromise the performance or availability
// of your application. Note the following about changing DesiredCapacity, MaxSize,
// or MinSize:
//
// * If a scale-in activity occurs as a result of a new
// DesiredCapacity value that is lower than the current size of the group, the Auto
// Scaling group uses its termination policy to determine which instances to
// terminate.
//
// * If you specify a new value for MinSize without specifying a value
// for DesiredCapacity, and the new MinSize is larger than the current size of the
// group, this sets the group's DesiredCapacity to the new MinSize value.
//
// * If you
// specify a new value for MaxSize without specifying a value for DesiredCapacity,
// and the new MaxSize is smaller than the current size of the group, this sets the
// group's DesiredCapacity to the new MaxSize value.
//
// To see which parameters have
// been set, call the DescribeAutoScalingGroups API. To view the scaling policies
// for an Auto Scaling group, call the DescribePolicies API. If the group has
// scaling policies, you can update them by calling the PutScalingPolicy API.
func (c *Client) UpdateAutoScalingGroup(ctx context.Context, params *UpdateAutoScalingGroupInput, optFns ...func(*Options)) (*UpdateAutoScalingGroupOutput, error) {
	if params == nil {
		params = &UpdateAutoScalingGroupInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "UpdateAutoScalingGroup", params, optFns, addOperationUpdateAutoScalingGroupMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*UpdateAutoScalingGroupOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type UpdateAutoScalingGroupInput struct {

	// The name of the Auto Scaling group.
	//
	// This member is required.
	AutoScalingGroupName *string

	// One or more Availability Zones for the group.
	AvailabilityZones []string

	// Enables or disables Capacity Rebalancing. For more information, see Amazon EC2
	// Auto Scaling Capacity Rebalancing
	// (https://docs.aws.amazon.com/autoscaling/ec2/userguide/capacity-rebalance.html)
	// in the Amazon EC2 Auto Scaling User Guide.
	CapacityRebalance *bool

	// The amount of time, in seconds, after a scaling activity completes before
	// another scaling activity can start. The default value is 300. This setting
	// applies when using simple scaling policies, but not when using other scaling
	// policies or scheduled scaling. For more information, see Scaling cooldowns for
	// Amazon EC2 Auto Scaling
	// (https://docs.aws.amazon.com/autoscaling/ec2/userguide/Cooldown.html) in the
	// Amazon EC2 Auto Scaling User Guide.
	DefaultCooldown *int32

	// The desired capacity is the initial capacity of the Auto Scaling group after
	// this operation completes and the capacity it attempts to maintain. This number
	// must be greater than or equal to the minimum size of the group and less than or
	// equal to the maximum size of the group.
	DesiredCapacity *int32

	// The amount of time, in seconds, that Amazon EC2 Auto Scaling waits before
	// checking the health status of an EC2 instance that has come into service. The
	// default value is 0. For more information, see Health check grace period
	// (https://docs.aws.amazon.com/autoscaling/ec2/userguide/healthcheck.html#health-check-grace-period)
	// in the Amazon EC2 Auto Scaling User Guide. Conditional: Required if you are
	// adding an ELB health check.
	HealthCheckGracePeriod *int32

	// The service to use for the health checks. The valid values are EC2 and ELB. If
	// you configure an Auto Scaling group to use ELB health checks, it considers the
	// instance unhealthy if it fails either the EC2 status checks or the load balancer
	// health checks.
	HealthCheckType *string

	// The name of the launch configuration. If you specify LaunchConfigurationName in
	// your update request, you can't specify LaunchTemplate or MixedInstancesPolicy.
	LaunchConfigurationName *string

	// The launch template and version to use to specify the updates. If you specify
	// LaunchTemplate in your update request, you can't specify LaunchConfigurationName
	// or MixedInstancesPolicy.
	LaunchTemplate *types.LaunchTemplateSpecification

	// The maximum amount of time, in seconds, that an instance can be in service. The
	// default is null. If specified, the value must be either 0 or a number equal to
	// or greater than 86,400 seconds (1 day). To clear a previously set value, specify
	// a new value of 0. For more information, see Replacing Auto Scaling instances
	// based on maximum instance lifetime
	// (https://docs.aws.amazon.com/autoscaling/ec2/userguide/asg-max-instance-lifetime.html)
	// in the Amazon EC2 Auto Scaling User Guide.
	MaxInstanceLifetime *int32

	// The maximum size of the Auto Scaling group. With a mixed instances policy that
	// uses instance weighting, Amazon EC2 Auto Scaling may need to go above MaxSize to
	// meet your capacity requirements. In this event, Amazon EC2 Auto Scaling will
	// never go above MaxSize by more than your largest instance weight (weights that
	// define how many units each instance contributes to the desired capacity of the
	// group).
	MaxSize *int32

	// The minimum size of the Auto Scaling group.
	MinSize *int32

	// An embedded object that specifies a mixed instances policy. When you make
	// changes to an existing policy, all optional parameters are left unchanged if not
	// specified. For more information, see Auto Scaling groups with multiple instance
	// types and purchase options
	// (https://docs.aws.amazon.com/autoscaling/ec2/userguide/asg-purchase-options.html)
	// in the Amazon EC2 Auto Scaling User Guide.
	MixedInstancesPolicy *types.MixedInstancesPolicy

	// Indicates whether newly launched instances are protected from termination by
	// Amazon EC2 Auto Scaling when scaling in. For more information about preventing
	// instances from terminating on scale in, see Instance scale-in protection
	// (https://docs.aws.amazon.com/autoscaling/ec2/userguide/as-instance-termination.html#instance-protection)
	// in the Amazon EC2 Auto Scaling User Guide.
	NewInstancesProtectedFromScaleIn *bool

	// The name of an existing placement group into which to launch your instances, if
	// any. A placement group is a logical grouping of instances within a single
	// Availability Zone. You cannot specify multiple Availability Zones and a
	// placement group. For more information, see Placement Groups
	// (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/placement-groups.html) in
	// the Amazon EC2 User Guide for Linux Instances.
	PlacementGroup *string

	// The Amazon Resource Name (ARN) of the service-linked role that the Auto Scaling
	// group uses to call other AWS services on your behalf. For more information, see
	// Service-linked roles
	// (https://docs.aws.amazon.com/autoscaling/ec2/userguide/autoscaling-service-linked-role.html)
	// in the Amazon EC2 Auto Scaling User Guide.
	ServiceLinkedRoleARN *string

	// A policy or a list of policies that are used to select the instances to
	// terminate. The policies are executed in the order that you list them. For more
	// information, see Controlling which Auto Scaling instances terminate during scale
	// in
	// (https://docs.aws.amazon.com/autoscaling/ec2/userguide/as-instance-termination.html)
	// in the Amazon EC2 Auto Scaling User Guide.
	TerminationPolicies []string

	// A comma-separated list of subnet IDs for a virtual private cloud (VPC). If you
	// specify VPCZoneIdentifier with AvailabilityZones, the subnets that you specify
	// for this parameter must reside in those Availability Zones.
	VPCZoneIdentifier *string
}

type UpdateAutoScalingGroupOutput struct {
	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata
}

func addOperationUpdateAutoScalingGroupMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsquery_serializeOpUpdateAutoScalingGroup{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpUpdateAutoScalingGroup{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addOpUpdateAutoScalingGroupValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opUpdateAutoScalingGroup(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opUpdateAutoScalingGroup(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "autoscaling",
		OperationName: "UpdateAutoScalingGroup",
	}
}
