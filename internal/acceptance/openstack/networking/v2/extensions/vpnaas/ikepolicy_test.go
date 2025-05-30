//go:build acceptance || networking || vpnaas

package vpnaas

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/clients"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/networking/v2/extensions/vpnaas/ikepolicies"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestIKEPolicyList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	allPages, err := ikepolicies.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allPolicies, err := ikepolicies.ExtractPolicies(allPages)
	th.AssertNoErr(t, err)

	for _, policy := range allPolicies {
		tools.PrintResource(t, policy)
	}
}

func TestIKEPolicyCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	policy, err := CreateIKEPolicy(t, client)
	th.AssertNoErr(t, err)
	defer DeleteIKEPolicy(t, client, policy.ID)

	tools.PrintResource(t, policy)

	newPolicy, err := ikepolicies.Get(context.TODO(), client, policy.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, newPolicy)

	updatedName := "updatedname"
	updatedDescription := "updated policy"
	updateOpts := ikepolicies.UpdateOpts{
		Name:        &updatedName,
		Description: &updatedDescription,
		Lifetime: &ikepolicies.LifetimeUpdateOpts{
			Value: 7000,
		},
	}
	updatedPolicy, err := ikepolicies.Update(context.TODO(), client, policy.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, updatedPolicy)
}
