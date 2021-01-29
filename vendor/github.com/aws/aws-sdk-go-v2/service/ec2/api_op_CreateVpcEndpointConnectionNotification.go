// Code generated by smithy-go-codegen DO NOT EDIT.

package ec2

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Creates a connection notification for a specified VPC endpoint or VPC endpoint
// service. A connection notification notifies you of specific endpoint events. You
// must create an SNS topic to receive notifications. For more information, see
// Create a Topic (https://docs.aws.amazon.com/sns/latest/dg/CreateTopic.html) in
// the Amazon Simple Notification Service Developer Guide. You can create a
// connection notification for interface endpoints only.
func (c *Client) CreateVpcEndpointConnectionNotification(ctx context.Context, params *CreateVpcEndpointConnectionNotificationInput, optFns ...func(*Options)) (*CreateVpcEndpointConnectionNotificationOutput, error) {
	if params == nil {
		params = &CreateVpcEndpointConnectionNotificationInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CreateVpcEndpointConnectionNotification", params, optFns, addOperationCreateVpcEndpointConnectionNotificationMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CreateVpcEndpointConnectionNotificationOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type CreateVpcEndpointConnectionNotificationInput struct {

	// One or more endpoint events for which to receive notifications. Valid values are
	// Accept, Connect, Delete, and Reject.
	//
	// This member is required.
	ConnectionEvents []string

	// The ARN of the SNS topic for the notifications.
	//
	// This member is required.
	ConnectionNotificationArn *string

	// Unique, case-sensitive identifier that you provide to ensure the idempotency of
	// the request. For more information, see How to Ensure Idempotency
	// (https://docs.aws.amazon.com/AWSEC2/latest/APIReference/Run_Instance_Idempotency.html).
	ClientToken *string

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation. Otherwise, it is
	// UnauthorizedOperation.
	DryRun bool

	// The ID of the endpoint service.
	ServiceId *string

	// The ID of the endpoint.
	VpcEndpointId *string
}

type CreateVpcEndpointConnectionNotificationOutput struct {

	// Unique, case-sensitive identifier that you provide to ensure the idempotency of
	// the request.
	ClientToken *string

	// Information about the notification.
	ConnectionNotification *types.ConnectionNotification

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata
}

func addOperationCreateVpcEndpointConnectionNotificationMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsEc2query_serializeOpCreateVpcEndpointConnectionNotification{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpCreateVpcEndpointConnectionNotification{}, middleware.After)
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
	if err = addOpCreateVpcEndpointConnectionNotificationValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCreateVpcEndpointConnectionNotification(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opCreateVpcEndpointConnectionNotification(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "ec2",
		OperationName: "CreateVpcEndpointConnectionNotification",
	}
}
