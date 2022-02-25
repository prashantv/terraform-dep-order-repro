package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/prashantv/tf-test/internal/filestore"
)

type Policy = filestore.Policy

func policyResource() *schema.Resource {
	r := &genericResource[Policy]{
		name: "Policy",
		build: func(d *schema.ResourceData) Policy {
			return Policy{
				NotifierIDs: toStringSlice(d.Get("notifier_ids").([]interface{})),
			}
		},
		set: func(obj Policy, d *schema.ResourceData) error {
			return d.Set("notifier_ids", obj.NotifierIDs)
		},
		writeStore:  filestore.WritePolicy,
		readStore:   filestore.ReadPolicy,
		deleteStore: filestore.DeletePolicy,
	}

	return &schema.Resource{
		Description: "Policy resource",

		CreateContext: r.Create,
		ReadContext:   r.Read,
		UpdateContext: r.Update,
		DeleteContext: r.Delete,

		Schema: map[string]*schema.Schema{
			"notifier_ids": {
				Description: "tier identifier",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
			},
		},
	}
}

func toStringSlice(l []interface{}) []string {
	var strs []string
	for _, v := range l {
		strs = append(strs, v.(string))
	}
	return strs
}
