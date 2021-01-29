// Code generated by smithy-go-codegen DO NOT EDIT.

package cloudformation

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Lists all exported output values in the account and Region in which you call
// this action. Use this action to see the exported output values that you can
// import into other stacks. To import values, use the Fn::ImportValue
// (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-importvalue.html)
// function. For more information, see  AWS CloudFormation Export Stack Output
// Values
// (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-stack-exports.html).
func (c *Client) ListExports(ctx context.Context, params *ListExportsInput, optFns ...func(*Options)) (*ListExportsOutput, error) {
	if params == nil {
		params = &ListExportsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ListExports", params, optFns, addOperationListExportsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ListExportsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ListExportsInput struct {

	// A string (provided by the ListExports response output) that identifies the next
	// page of exported output values that you asked to retrieve.
	NextToken *string
}

type ListExportsOutput struct {

	// The output for the ListExports action.
	Exports []types.Export

	// If the output exceeds 100 exported output values, a string that identifies the
	// next page of exports. If there is no additional page, this value is null.
	NextToken *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata
}

func addOperationListExportsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsquery_serializeOpListExports{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpListExports{}, middleware.After)
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
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opListExports(options.Region), middleware.Before); err != nil {
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

// ListExportsAPIClient is a client that implements the ListExports operation.
type ListExportsAPIClient interface {
	ListExports(context.Context, *ListExportsInput, ...func(*Options)) (*ListExportsOutput, error)
}

var _ ListExportsAPIClient = (*Client)(nil)

// ListExportsPaginatorOptions is the paginator options for ListExports
type ListExportsPaginatorOptions struct {
	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// ListExportsPaginator is a paginator for ListExports
type ListExportsPaginator struct {
	options   ListExportsPaginatorOptions
	client    ListExportsAPIClient
	params    *ListExportsInput
	nextToken *string
	firstPage bool
}

// NewListExportsPaginator returns a new ListExportsPaginator
func NewListExportsPaginator(client ListExportsAPIClient, params *ListExportsInput, optFns ...func(*ListExportsPaginatorOptions)) *ListExportsPaginator {
	options := ListExportsPaginatorOptions{}

	for _, fn := range optFns {
		fn(&options)
	}

	if params == nil {
		params = &ListExportsInput{}
	}

	return &ListExportsPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *ListExportsPaginator) HasMorePages() bool {
	return p.firstPage || p.nextToken != nil
}

// NextPage retrieves the next ListExports page.
func (p *ListExportsPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*ListExportsOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	result, err := p.client.ListExports(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextToken

	if p.options.StopOnDuplicateToken && prevToken != nil && p.nextToken != nil && *prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

func newServiceMetadataMiddleware_opListExports(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "cloudformation",
		OperationName: "ListExports",
	}
}
