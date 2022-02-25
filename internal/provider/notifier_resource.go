package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/prashantv/tf-test/internal/filestore"
)

type Notifier = filestore.Notifier

func notifierResource() *schema.Resource {
	r := &genericResource[Notifier]{
		name: "Notifier",
		build: func(d *schema.ResourceData) Notifier {
			return Notifier{
				Email: d.Get("email").(string),
			}
		},
		set: func(obj Notifier, d *schema.ResourceData) error {
			return d.Set("email", obj.Email)
		},
		writeStore:  filestore.WriteNotifier,
		readStore:   filestore.ReadNotifier,
		deleteStore: filestore.DeleteNotifier,
	}

	return &schema.Resource{
		Description: "Notifier resource",

		CreateContext: r.Create,
		ReadContext:   r.Read,
		UpdateContext: r.Update,
		DeleteContext: r.Delete,

		Schema: map[string]*schema.Schema{
			"email": {
				Description: "Email",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}
