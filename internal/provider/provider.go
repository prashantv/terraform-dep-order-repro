package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{
				"tftest_notifier": notifierResource(),
				"tftest_policy":   policyResource(),
			},
		}

		return p
	}
}
