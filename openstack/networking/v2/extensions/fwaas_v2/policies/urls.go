package policies

import "github.com/vnpaycloud-console/gophercloud/v2"

const (
	rootPath     = "fwaas"
	resourcePath = "firewall_policies"
	insertPath   = "insert_rule"
	removePath   = "remove_rule"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func insertURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id, insertPath)
}

func removeURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id, removePath)
}
