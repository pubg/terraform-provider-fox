package fox

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-fox/pkg/common"
	"terraform-provider-fox/pkg/ip"
)

func resourceEnv() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvCreate,
		ReadContext:   resourceEnvRead,
		UpdateContext: resourceEnvUpdate,
		DeleteContext: resourceEnvDelete,
		Schema: map[string]*schema.Schema{
			"env": {
				Type:     schema.TypeString,
				Required: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cidrs": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
				// ValidateFunc: validation.ValidateRFC3339TimeString,	// ValidateRFC3339TimeString is deprecated
			},
			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
				// ValidateFunc: validation.ValidateRFC3339TimeString,	// ValidateRFC3339TimeString is deprecated
			},
		},
	}
}

func resourceEnvCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(common.Config)

	// get info
	env := d.Get("env").(string)
	groups := d.Get("groups").([]interface{})
	cidrs := d.Get("cidrs").([]interface{})

	// create resource
	err := ip.CreateIpInfo(config, env, groups, cidrs, &diags)
	if err != nil {
		return diags
	}

	// set id
	d.SetId(env)

	resourceEnvRead(ctx, d, m)

	return diags
}

func resourceEnvRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(common.Config)

	// get id
	id := d.Id()
	if id == "" {
		err := errors.New("envID is empty")
		return diag.FromErr(err)
	}

	// get data
	ipInfo, err := ip.GetIpInfo(config, id, &diags)
	if err != nil {
		return diags
	}

	// set data
	err, subErrMsg := bindDataSrcEnv(d, ipInfo)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("bind IpInfo to data resource fail: %s", subErrMsg),
			Detail:   fmt.Sprintf("ipInfo: %+v\nerror: %s", ipInfo, err.Error()),
		})
		return diags
	}

	return diags
}

func resourceEnvUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(common.Config)

	// get id
	id := d.Id()

	if d.HasChanges("groups", "cidrs") {
		// get info
		groups := d.Get("groups").([]interface{})
		cidrs := d.Get("cidrs").([]interface{})

		// update resource
		err := ip.UpdateIpInfo(config, id, groups, cidrs, &diags)
		if err != nil {
			return diags
		}
	}

	return resourceEnvRead(ctx, d, m)
}

func resourceEnvDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(common.Config)

	// get id
	id := d.Id()

	// delete resource
	err := ip.DeleteIpInfo(config, id, &diags)
	if err != nil {
		return diags
	}

	return diags
}
