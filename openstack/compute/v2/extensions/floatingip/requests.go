package floatingip

import (
	"errors"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// List returns a Pager that allows you to iterate over a collection of FloatingIPs.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, listURL(client), func(r pagination.PageResult) pagination.Page {
		return FloatingIPsPage{pagination.SinglePageBase(r)}
	})
}

// CreateOptsBuilder describes struct types that can be accepted by the Create call. Notable, the
// CreateOpts struct in this package does.
type CreateOptsBuilder interface {
	ToFloatingIPCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies a Floating IP allocation request
type CreateOpts struct {
	// Pool is the pool of floating IPs to allocate one from
	Pool string
}

// ToFloatingIPCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToFloatingIPCreateMap() (map[string]interface{}, error) {
	if opts.Pool == "" {
		return nil, errors.New("Missing field required for floating IP creation: Pool")
	}

	return map[string]interface{}{"pool": opts.Pool}, nil
}

// Create requests the creation of a new floating IP
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToFloatingIPCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = perigee.Request("POST", createURL(client), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     reqBody,
		Results:     &res.Body,
		OkCodes:     []int{200},
	})
	return res
}

// Get returns data about a previously created FloatingIP.
func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = perigee.Request("GET", getURL(client, id), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})
	return res
}

// Delete requests the deletion of a previous allocated FloatingIP.
func Delete(client *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, res.Err = perigee.Request("DELETE", deleteURL(client, id), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})
	return res
}

// association / disassociation

// Associate pairs an allocated floating IP with an instance
func Associate(client *gophercloud.ServiceClient, serverId, fip string) AssociateResult {
	var res AssociateResult

	addFloatingIp := make(map[string]interface{})
	addFloatingIp["address"] = fip
	reqBody := map[string]interface{}{"addFloatingIp": addFloatingIp}

	_, res.Err = perigee.Request("POST", associateURL(client, serverId), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     reqBody,
		OkCodes:     []int{202},
	})
	return res
}

// Disassociate decouples an allocated floating IP from an instance
func Disassociate(client *gophercloud.ServiceClient, serverId, fip string) DisassociateResult {
	var res DisassociateResult

	removeFloatingIp := make(map[string]interface{})
	removeFloatingIp["address"] = fip
	reqBody := map[string]interface{}{"removeFloatingIp": removeFloatingIp}

	_, res.Err = perigee.Request("POST", disassociateURL(client, serverId), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     reqBody,
		OkCodes:     []int{202},
	})
	return res
}
