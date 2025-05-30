package introspection

import (
	"context"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

// ListIntrospectionsOptsBuilder allows extensions to add additional parameters to the
// ListIntrospections request.
type ListIntrospectionsOptsBuilder interface {
	ToIntrospectionsListQuery() (string, error)
}

// ListIntrospectionsOpts allows the filtering and sorting of paginated collections through
// the Introspection API. Filtering is achieved by passing in struct field values that map to
// the node attributes you want to see returned. Marker and Limit are used
// for pagination.
type ListIntrospectionsOpts struct {
	// Requests a page size of items.
	Limit int `q:"limit"`

	// The ID of the last-seen item.
	Marker string `q:"marker"`
}

// ToIntrospectionsListQuery formats a ListIntrospectionsOpts into a query string.
func (opts ListIntrospectionsOpts) ToIntrospectionsListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListIntrospections makes a request against the Inspector API to list the current introspections.
func ListIntrospections(client *gophercloud.ServiceClient, opts ListIntrospectionsOptsBuilder) pagination.Pager {
	url := listIntrospectionsURL(client)
	if opts != nil {
		query, err := opts.ToIntrospectionsListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		var rpage = IntrospectionPage{pagination.LinkedPageBase{PageResult: r}}
		return rpage
	})
}

// GetIntrospectionStatus makes a request against the Inspector API to get the
// status of a single introspection.
func GetIntrospectionStatus(ctx context.Context, client *gophercloud.ServiceClient, nodeID string) (r GetIntrospectionStatusResult) {
	resp, err := client.Get(ctx, introspectionURL(client, nodeID), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// StartOptsBuilder allows extensions to add additional parameters to the
// Start request.
type StartOptsBuilder interface {
	ToStartIntrospectionQuery() (string, error)
}

// StartOpts represents options to start an introspection.
type StartOpts struct {
	// Whether the current installation of ironic-inspector can manage PXE booting of nodes.
	ManageBoot *bool `q:"manage_boot"`
}

// ToStartIntrospectionQuery converts a StartOpts into a request.
func (opts StartOpts) ToStartIntrospectionQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// StartIntrospection initiate hardware introspection for node NodeID .
// All power management configuration for this node needs to be done prior to calling the endpoint.
func StartIntrospection(ctx context.Context, client *gophercloud.ServiceClient, nodeID string, opts StartOptsBuilder) (r StartResult) {
	_, err := opts.ToStartIntrospectionQuery()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, introspectionURL(client, nodeID), nil, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// AbortIntrospection abort running introspection.
func AbortIntrospection(ctx context.Context, client *gophercloud.ServiceClient, nodeID string) (r AbortResult) {
	resp, err := client.Post(ctx, abortIntrospectionURL(client, nodeID), nil, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetIntrospectionData return stored data from successful introspection.
func GetIntrospectionData(ctx context.Context, client *gophercloud.ServiceClient, nodeID string) (r DataResult) {
	resp, err := client.Get(ctx, introspectionDataURL(client, nodeID), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ReApplyIntrospection triggers introspection on stored unprocessed data.
// No data is allowed to be sent along with the request.
func ReApplyIntrospection(ctx context.Context, client *gophercloud.ServiceClient, nodeID string) (r ApplyDataResult) {
	resp, err := client.Post(ctx, introspectionUnprocessedDataURL(client, nodeID), nil, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
