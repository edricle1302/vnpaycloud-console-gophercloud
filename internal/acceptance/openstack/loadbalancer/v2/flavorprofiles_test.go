//go:build acceptance || networking || loadbalancer || flavorprofiles

package v2

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/clients"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/ptr"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/loadbalancer/v2/flavorprofiles"

	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestFlavorProfilesList(t *testing.T) {
	client, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	allPages, err := flavorprofiles.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allFlavorProfiles, err := flavorprofiles.ExtractFlavorProfiles(allPages)
	th.AssertNoErr(t, err)

	for _, flavorprofile := range allFlavorProfiles {
		tools.PrintResource(t, flavorprofile)
	}
}

func TestFlavorProfilesCRUD(t *testing.T) {
	lbClient, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	flavorProfile, err := CreateFlavorProfile(t, lbClient)
	th.AssertNoErr(t, err)
	defer DeleteFlavorProfile(t, lbClient, flavorProfile)

	tools.PrintResource(t, flavorProfile)

	th.AssertEquals(t, "amphora", flavorProfile.ProviderName)

	flavorProfileUpdateOpts := flavorprofiles.UpdateOpts{
		Name: ptr.To(tools.RandomString("TESTACCTUP-", 8)),
	}

	flavorProfileUpdated, err := flavorprofiles.Update(context.TODO(), lbClient, flavorProfile.ID, flavorProfileUpdateOpts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, *flavorProfileUpdateOpts.Name, flavorProfileUpdated.Name)

	t.Logf("Successfully updated flavorprofile %s", flavorProfileUpdated.Name)
}
