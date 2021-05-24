package fox

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"terraform-provider-fox/pkg/common"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"fox_ip_env": resourceEnv(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"fox_ping":     dataSourcePing(),
			"fox_ip_env":   dataSourceEnv(),
			"fox_ip_group": dataSourceGroup(),
			"fox_ip_all":   dataSourceAll(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	address := d.Get("address").(string)
	config := common.Config{
		Address: address,
	}
	log.Printf("[DEBUG] Use Fox address (%s)\n", address)
	return config, diags
}
