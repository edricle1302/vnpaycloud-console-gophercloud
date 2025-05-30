package ec2credentials

import "github.com/vnpaycloud-console/gophercloud/v2"

func listURL(client *gophercloud.ServiceClient, userID string) string {
	return client.ServiceURL("users", userID, "credentials", "OS-EC2")
}

func getURL(client *gophercloud.ServiceClient, userID string, id string) string {
	return client.ServiceURL("users", userID, "credentials", "OS-EC2", id)
}

func createURL(client *gophercloud.ServiceClient, userID string) string {
	return client.ServiceURL("users", userID, "credentials", "OS-EC2")
}

func deleteURL(client *gophercloud.ServiceClient, userID string, id string) string {
	return client.ServiceURL("users", userID, "credentials", "OS-EC2", id)
}
