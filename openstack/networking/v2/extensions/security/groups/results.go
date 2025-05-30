package groups

import (
	"encoding/json"
	"time"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

// SecGroup represents a container for security group rules.
type SecGroup struct {
	// The UUID for the security group.
	ID string

	// Human-readable name for the security group. Might not be unique.
	// Cannot be named "default" as that is automatically created for a tenant.
	Name string

	// The security group description.
	Description string

	// A slice of security group rules that dictate the permitted behaviour for
	// traffic entering and leaving the group.
	Rules []rules.SecGroupRule `json:"security_group_rules"`

	// Indicates if the security group is stateful or stateless.
	Stateful bool `json:"stateful"`

	// TenantID is the project owner of the security group.
	TenantID string `json:"tenant_id"`

	// UpdatedAt and CreatedAt contain ISO-8601 timestamps of when the state of the
	// security group last changed, and when it was created.
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`

	// ProjectID is the project owner of the security group.
	ProjectID string `json:"project_id"`

	// Tags optionally set via extensions/attributestags
	Tags []string `json:"tags"`
}

func (r *SecGroup) UnmarshalJSON(b []byte) error {
	type tmp SecGroup

	// Support for older neutron time format
	var s1 struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339NoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339NoZ `json:"updated_at"`
	}

	err := json.Unmarshal(b, &s1)
	if err == nil {
		*r = SecGroup(s1.tmp)
		r.CreatedAt = time.Time(s1.CreatedAt)
		r.UpdatedAt = time.Time(s1.UpdatedAt)

		return nil
	}

	// Support for newer neutron time format
	var s2 struct {
		tmp
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	err = json.Unmarshal(b, &s2)
	if err != nil {
		return err
	}

	*r = SecGroup(s2.tmp)
	r.CreatedAt = time.Time(s2.CreatedAt)
	r.UpdatedAt = time.Time(s2.UpdatedAt)

	return nil
}

// SecGroupPage is the page returned by a pager when traversing over a
// collection of security groups.
type SecGroupPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of security groups has
// reached the end of a page and the pager seeks to traverse over a new one. In
// order to do this, it needs to construct the next page's URL.
func (r SecGroupPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"security_groups_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}

	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a SecGroupPage struct is empty.
func (r SecGroupPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractGroups(r)
	return len(is) == 0, err
}

// ExtractGroups accepts a Page struct, specifically a SecGroupPage struct,
// and extracts the elements into a slice of SecGroup structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractGroups(r pagination.Page) ([]SecGroup, error) {
	var s struct {
		SecGroups []SecGroup `json:"security_groups"`
	}
	err := (r.(SecGroupPage)).ExtractInto(&s)
	return s.SecGroups, err
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a security group.
func (r commonResult) Extract() (*SecGroup, error) {
	var s struct {
		SecGroup *SecGroup `json:"security_group"`
	}
	err := r.ExtractInto(&s)
	return s.SecGroup, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a SecGroup.
type CreateResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a SecGroup.
type UpdateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a SecGroup.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}
