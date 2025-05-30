//go:build acceptance || objectstorage || containers

package v1

import (
	"context"
	"strings"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/clients"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/objectstorage/v1/containers"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

// numContainers is the number of containers to create for testing.
var numContainers = 2

func TestContainers(t *testing.T) {
	client, err := clients.NewObjectStorageV1Client()
	if err != nil {
		t.Fatalf("Unable to create client: %v", err)
	}

	// Create a slice of random container names.
	cNames := make([]string, numContainers)
	for i := 0; i < numContainers; i++ {
		cNames[i] = "gophercloud-test-container-" + tools.RandomFunnyStringNoSlash(8)
	}

	// Create numContainers containers.
	for i := 0; i < len(cNames); i++ {
		res := containers.Create(context.TODO(), client, cNames[i], nil)
		th.AssertNoErr(t, res.Err)
	}
	// Delete the numContainers containers after function completion.
	defer func() {
		for i := 0; i < len(cNames); i++ {
			res := containers.Delete(context.TODO(), client, cNames[i])
			th.AssertNoErr(t, res.Err)
		}
	}()

	// List the numContainer names that were just created. To just list those,
	// the 'prefix' parameter is used.
	err = containers.List(client, &containers.ListOpts{Prefix: "gophercloud-test-container-"}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		containerList, err := containers.ExtractInfo(page)
		th.AssertNoErr(t, err)

		for _, n := range containerList {
			t.Logf("Container: Name [%s] Count [%d] Bytes [%d]",
				n.Name, n.Count, n.Bytes)
		}

		return true, nil
	})
	th.AssertNoErr(t, err)

	// List the info for the numContainer containers that were created.
	err = containers.List(client, &containers.ListOpts{Prefix: "gophercloud-test-container-"}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		containerList, err := containers.ExtractNames(page)
		th.AssertNoErr(t, err)
		for _, n := range containerList {
			t.Logf("Container: Name [%s]", n)
		}

		return true, nil
	})
	th.AssertNoErr(t, err)

	// Update one of the numContainer container metadata.
	metadata := map[string]string{
		"Gophercloud-Test": "containers",
	}
	read := ".r:*,.rlistings"
	write := "*:*"
	iTrue := true
	empty := ""
	opts := &containers.UpdateOpts{
		Metadata:          metadata,
		ContainerRead:     &read,
		ContainerWrite:    &write,
		DetectContentType: new(bool),
		ContainerSyncTo:   &empty,
		ContainerSyncKey:  &empty,
	}

	updateres := containers.Update(context.TODO(), client, cNames[0], opts)
	th.AssertNoErr(t, updateres.Err)
	// After the tests are done, delete the metadata that was set.
	defer func() {
		temp := []string{}
		for k := range metadata {
			temp = append(temp, k)
		}
		empty := ""
		opts = &containers.UpdateOpts{
			RemoveMetadata:    temp,
			ContainerRead:     &empty,
			ContainerWrite:    &empty,
			DetectContentType: &iTrue,
		}
		res := containers.Update(context.TODO(), client, cNames[0], opts)
		th.AssertNoErr(t, res.Err)

		// confirm the metadata was removed
		getOpts := containers.GetOpts{
			Newest: true,
		}

		resp := containers.Get(context.TODO(), client, cNames[0], getOpts)
		cm, err := resp.ExtractMetadata()
		th.AssertNoErr(t, err)
		for k := range metadata {
			if _, ok := cm[k]; ok {
				t.Errorf("Unexpected custom metadata with key: %s", k)
			}
		}
		container, err := resp.Extract()
		th.AssertNoErr(t, err)
		th.AssertEquals(t, empty, strings.Join(container.Read, ","))
		th.AssertEquals(t, empty, strings.Join(container.Write, ","))
	}()

	// Retrieve a container's metadata.
	getOpts := containers.GetOpts{
		Newest: true,
	}

	resp := containers.Get(context.TODO(), client, cNames[0], getOpts)
	cm, err := resp.ExtractMetadata()
	th.AssertNoErr(t, err)
	for k := range metadata {
		if cm[k] != metadata[strings.Title(k)] {
			t.Errorf("Expected custom metadata with key: %s", k)
		}
	}
	container, err := resp.Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, read, strings.Join(container.Read, ","))
	th.AssertEquals(t, write, strings.Join(container.Write, ","))

	// Retrieve a container's timestamp
	cHeaders, err := containers.Get(context.TODO(), client, cNames[0], getOpts).Extract()
	th.AssertNoErr(t, err)
	t.Logf("Container: Name [%s] Timestamp: [%f]\n", cNames[0], cHeaders.Timestamp)
}

func TestListAllContainers(t *testing.T) {
	client, err := clients.NewObjectStorageV1Client()
	if err != nil {
		t.Fatalf("Unable to create client: %v", err)
	}

	numContainers := 20

	// Create a slice of random container names.
	cNames := make([]string, numContainers)
	for i := 0; i < numContainers; i++ {
		cNames[i] = "gophercloud-test-container-" + tools.RandomFunnyStringNoSlash(8)
	}

	// Create numContainers containers.
	for i := 0; i < len(cNames); i++ {
		res := containers.Create(context.TODO(), client, cNames[i], nil)
		th.AssertNoErr(t, res.Err)
	}
	// Delete the numContainers containers after function completion.
	defer func() {
		for i := 0; i < len(cNames); i++ {
			res := containers.Delete(context.TODO(), client, cNames[i])
			th.AssertNoErr(t, res.Err)
		}
	}()

	// List all the numContainer names that were just created. To just list those,
	// the 'prefix' parameter is used.
	allPages, err := containers.List(client, &containers.ListOpts{Limit: 5, Prefix: "gophercloud-test-container-"}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	containerInfoList, err := containers.ExtractInfo(allPages)
	th.AssertNoErr(t, err)
	for _, n := range containerInfoList {
		t.Logf("Container: Name [%s] Count [%d] Bytes [%d]",
			n.Name, n.Count, n.Bytes)
	}
	th.AssertEquals(t, numContainers, len(containerInfoList))

	// List the info for all the numContainer containers that were created.
	allPages, err = containers.List(client, &containers.ListOpts{Limit: 2, Prefix: "gophercloud-test-container-"}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	containerNamesList, err := containers.ExtractNames(allPages)
	th.AssertNoErr(t, err)
	for _, n := range containerNamesList {
		t.Logf("Container: Name [%s]", n)
	}
	th.AssertEquals(t, numContainers, len(containerNamesList))
}

func TestBulkDeleteContainers(t *testing.T) {
	client, err := clients.NewObjectStorageV1Client()
	if err != nil {
		t.Fatalf("Unable to create client: %v", err)
	}

	numContainers := 20

	// Create a slice of random container names.
	cNames := make([]string, numContainers)
	for i := 0; i < numContainers; i++ {
		cNames[i] = "gophercloud-test-container-" + tools.RandomFunnyStringNoSlash(8)
	}

	// Create numContainers containers.
	for i := 0; i < len(cNames); i++ {
		res := containers.Create(context.TODO(), client, cNames[i], nil)
		th.AssertNoErr(t, res.Err)
	}

	expectedResp := containers.BulkDeleteResponse{
		ResponseStatus: "200 OK",
		Errors:         [][]string{},
		NumberDeleted:  numContainers,
	}

	resp, err := containers.BulkDelete(context.TODO(), client, cNames).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, *resp)
	th.AssertDeepEquals(t, *resp, expectedResp)

	for _, c := range cNames {
		_, err = containers.Get(context.TODO(), client, c, nil).Extract()
		th.AssertErr(t, err)
	}
}
