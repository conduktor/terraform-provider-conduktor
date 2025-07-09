package console_partner_zone_v2

import (
	"context"
	"testing"

	ctlresource "github.com/conduktor/ctl/resource"
	console "github.com/conduktor/terraform-provider-conduktor/internal/model/console"
	"github.com/conduktor/terraform-provider-conduktor/internal/schema"
	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestPartnerZoneV2ModelMapping(t *testing.T) {

	ctx := context.Background()

	jsonPartnerZoneV2Resource := []byte(test.TestAccTestdata(t, "console/partner_zone_v2/api.json"))

	ctlResource := ctlresource.Resource{}
	err := ctlResource.UnmarshalJSON(jsonPartnerZoneV2Resource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "PartnerZone", ctlResource.Kind)
	assert.Equal(t, "v2", ctlResource.Version)
	assert.Equal(t, "partner-zone", ctlResource.Name)
	assert.Equal(t, map[string]any{"name": "partner-zone", "labels": map[string]any{"key": "value"}}, ctlResource.Metadata)
	assert.Equal(t, "mycluster", ctlResource.Spec["cluster"])
	assert.Equal(t, jsonPartnerZoneV2Resource, ctlResource.Json)

	// convert into internal model
	internal, err := console.NewPartnerZoneConsoleResourceFromClientResource(ctlResource)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "PartnerZone", internal.Kind)
	assert.Equal(t, "v2", internal.ApiVersion)
	assert.Equal(t, "partner-zone", internal.Metadata.Name)
	assert.Equal(t, map[string]string{
		"key": "value",
	}, internal.Metadata.Labels)
	assert.Equal(t, "mycluster", internal.Spec.Cluster)
	assert.Equal(t, "Partner Zone", internal.Spec.DisplayName)
	assert.Equal(t, "This is a partner zone", internal.Spec.Description)
	assert.Equal(t, "https://partnerzone.example.com", internal.Spec.Url)
	assert.Equal(t, "my-service-account", internal.Spec.AuthenticationMode.ServiceAccount)
	assert.Equal(t, "OAUTHBEARER", internal.Spec.AuthenticationMode.Type)
	expectedTopics := []console.PartnerZoneTopic{
		{
			Name:         "topic-1",
			BackingTopic: "backing-topic-1",
			Permission:   "READ",
		},
		{
			Name:         "topic-2",
			BackingTopic: "backing-topic-2",
			Permission:   "WRITE",
		},
	}
	assert.Equal(t, expectedTopics, internal.Spec.Topics)
	assert.Equal(t, "John Doe", internal.Spec.Partner.Name)
	assert.Equal(t, "Data analyst", internal.Spec.Partner.Role)
	assert.Equal(t, "johndoe@company.io", internal.Spec.Partner.Email)
	assert.Equal(t, "07827 837 177", internal.Spec.Partner.Phone)
	assert.Equal(t, int64(1000000), internal.Spec.TrafficControlPolicies.MaxProduceRate)
	assert.Equal(t, int64(2000000), internal.Spec.TrafficControlPolicies.MaxConsumeRate)
	assert.Equal(t, int64(30), internal.Spec.TrafficControlPolicies.LimitCommitOffset)
	expectedToAdd := []console.PartnerZoneToAdd{
		{
			Key:              "key-1",
			Value:            "value-1",
			OverrideIfExists: false,
		},
		{
			Key:              "key-2",
			Value:            "value-2",
			OverrideIfExists: true,
		},
	}
	expectedToRemove := []console.PartnerZoneToRemove{
		{
			KeyRegex: "my_org_prefix.*",
		},
	}
	assert.Equal(t, expectedToAdd, internal.Spec.Headers.AddOnProduce)
	assert.Equal(t, expectedToRemove, internal.Spec.Headers.RemoveOnConsume)

	// convert to terraform model
	tfModel, err := InternalModelToTerraform(ctx, &internal)
	if err != nil {
		t.Fatal(err)
		return
	}
	labels, _ := schema.StringMapToMapValue(ctx, map[string]string{"key": "value"})
	assert.Equal(t, types.StringValue("partner-zone"), tfModel.Name)
	assert.Equal(t, labels, tfModel.Labels)
	assert.Equal(t, types.StringValue("mycluster"), tfModel.Spec.Cluster)
	assert.Equal(t, types.StringValue("This is a partner zone"), tfModel.Spec.Description)
	assert.Equal(t, types.StringValue("Partner Zone"), tfModel.Spec.DisplayName)
	assert.Equal(t, false, tfModel.Spec.AuthenticationMode.IsNull())
	assert.Equal(t, false, tfModel.Spec.AuthenticationMode.IsUnknown())
	assert.Equal(t, false, tfModel.Spec.Topics.IsNull())
	assert.Equal(t, false, tfModel.Spec.Topics.IsUnknown())

	// convert back to internal model
	internal2, err := TFToInternalModel(ctx, &tfModel)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.Equal(t, "PartnerZone", internal2.Kind)
	assert.Equal(t, "v2", internal2.ApiVersion)
	assert.Equal(t, "partner-zone", internal2.Metadata.Name)
	assert.Equal(t, map[string]string{
		"key": "value",
	}, internal2.Metadata.Labels)
	assert.Equal(t, "mycluster", internal2.Spec.Cluster)
	assert.Equal(t, "This is a partner zone", internal2.Spec.Description)
	assert.Equal(t, "Partner Zone", internal2.Spec.DisplayName)
	assert.Equal(t, expectedTopics, internal2.Spec.Topics)
	assert.Equal(t, "John Doe", internal2.Spec.Partner.Name)
	assert.Equal(t, "Data analyst", internal2.Spec.Partner.Role)
	assert.Equal(t, "johndoe@company.io", internal2.Spec.Partner.Email)
	assert.Equal(t, "07827 837 177", internal2.Spec.Partner.Phone)
	assert.Equal(t, int64(1000000), internal2.Spec.TrafficControlPolicies.MaxProduceRate)
	assert.Equal(t, int64(2000000), internal2.Spec.TrafficControlPolicies.MaxConsumeRate)
	assert.Equal(t, int64(30), internal2.Spec.TrafficControlPolicies.LimitCommitOffset)
	assert.Equal(t, expectedToAdd, internal2.Spec.Headers.AddOnProduce)
	assert.Equal(t, expectedToRemove, internal2.Spec.Headers.RemoveOnConsume)

	// convert back to ctl model
	ctlResource2, err := internal2.ToClientResource()
	if err != nil {
		t.Fatal(err)
		return
	}
	// compare without json
	if !cmp.Equal(ctlResource, ctlResource2, cmpopts.IgnoreFields(ctlresource.Resource{}, "Json")) {
		t.Errorf("expected %+v, got %+v", ctlResource, ctlResource2)
	}
}
