// Code generated by smithy-go-codegen DO NOT EDIT.

package eks

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Lists the available add-ons.
func (c *Client) ListAddons(ctx context.Context, params *ListAddonsInput, optFns ...func(*Options)) (*ListAddonsOutput, error) {
	if params == nil {
		params = &ListAddonsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ListAddons", params, optFns, addOperationListAddonsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ListAddonsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ListAddonsInput struct {

	// The name of the cluster.
	//
	// This member is required.
	ClusterName *string

	// The maximum number of add-on results returned by ListAddonsRequest in paginated
	// output. When you use this parameter, ListAddonsRequest returns only maxResults
	// results in a single page along with a nextToken response element. You can see
	// the remaining results of the initial request by sending another
	// ListAddonsRequest request with the returned nextToken value. This value can be
	// between 1 and 100. If you don't use this parameter, ListAddonsRequest returns up
	// to 100 results and a nextToken value, if applicable.
	MaxResults *int32

	// The nextToken value returned from a previous paginated ListAddonsRequest where
	// maxResults was used and the results exceeded the value of that parameter.
	// Pagination continues from the end of the previous results that returned the
	// nextToken value. This token should be treated as an opaque identifier that is
	// used only to retrieve the next items in a list and not for other programmatic
	// purposes.
	NextToken *string
}

type ListAddonsOutput struct {

	// A list of available add-ons.
	Addons []string

	// The nextToken value returned from a previous paginated ListAddonsResponse where
	// maxResults was used and the results exceeded the value of that parameter.
	// Pagination continues from the end of the previous results that returned the
	// nextToken value. This token should be treated as an opaque identifier that is
	// used only to retrieve the next items in a list and not for other programmatic
	// purposes.
	NextToken *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata
}

func addOperationListAddonsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsRestjson1_serializeOpListAddons{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestjson1_deserializeOpListAddons{}, middleware.After)
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
	if err = addOpListAddonsValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opListAddons(options.Region), middleware.Before); err != nil {
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

// ListAddonsAPIClient is a client that implements the ListAddons operation.
type ListAddonsAPIClient interface {
	ListAddons(context.Context, *ListAddonsInput, ...func(*Options)) (*ListAddonsOutput, error)
}

var _ ListAddonsAPIClient = (*Client)(nil)

// ListAddonsPaginatorOptions is the paginator options for ListAddons
type ListAddonsPaginatorOptions struct {
	// The maximum number of add-on results returned by ListAddonsRequest in paginated
	// output. When you use this parameter, ListAddonsRequest returns only maxResults
	// results in a single page along with a nextToken response element. You can see
	// the remaining results of the initial request by sending another
	// ListAddonsRequest request with the returned nextToken value. This value can be
	// between 1 and 100. If you don't use this parameter, ListAddonsRequest returns up
	// to 100 results and a nextToken value, if applicable.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// ListAddonsPaginator is a paginator for ListAddons
type ListAddonsPaginator struct {
	options   ListAddonsPaginatorOptions
	client    ListAddonsAPIClient
	params    *ListAddonsInput
	nextToken *string
	firstPage bool
}

// NewListAddonsPaginator returns a new ListAddonsPaginator
func NewListAddonsPaginator(client ListAddonsAPIClient, params *ListAddonsInput, optFns ...func(*ListAddonsPaginatorOptions)) *ListAddonsPaginator {
	options := ListAddonsPaginatorOptions{}
	if params.MaxResults != nil {
		options.Limit = *params.MaxResults
	}

	for _, fn := range optFns {
		fn(&options)
	}

	if params == nil {
		params = &ListAddonsInput{}
	}

	return &ListAddonsPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *ListAddonsPaginator) HasMorePages() bool {
	return p.firstPage || p.nextToken != nil
}

// NextPage retrieves the next ListAddons page.
func (p *ListAddonsPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*ListAddonsOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxResults = limit

	result, err := p.client.ListAddons(ctx, &params, optFns...)
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

func newServiceMetadataMiddleware_opListAddons(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "eks",
		OperationName: "ListAddons",
	}
}
