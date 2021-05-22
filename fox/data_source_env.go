package fox

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-fox/pkg/common"
	"terraform-provider-fox/pkg/ip"
	"time"
)

func dataSourceEnv() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnvRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"env": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cidrs": {
				Type:     schema.TypeList,
				Computed: true,
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

func dataSourceEnvRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(common.Config)

	// get id
	id := d.Get("id").(string)
	if id == "" {
		err := errors.New("id is empty")
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

	// always run
	d.SetId(id)

	return diags
}

func bindDataSrcEnv(d *schema.ResourceData, ipInfo *ip.IpInfo) (error, string) {
	const subErrMsgFormat = "%s set fail"
	if err := d.Set("env", ipInfo.Env); err != nil {
		return err, fmt.Sprintf(subErrMsgFormat, "env")
	}
	if err := d.Set("groups", ipInfo.GroupArr); err != nil {
		return err, fmt.Sprintf(subErrMsgFormat, "groups")
	}
	if err := d.Set("cidrs", ipInfo.CidrArr); err != nil {
		return err, fmt.Sprintf(subErrMsgFormat, "cidrs")
	}
	if err := d.Set("created", ipInfo.Created.Format(time.RFC3339)); err != nil {
		return err, fmt.Sprintf(subErrMsgFormat, "created")
	}
	if err := d.Set("last_modified", ipInfo.LastModified.Format(time.RFC3339)); err != nil {
		return err, fmt.Sprintf(subErrMsgFormat, "last_modified")
	}
	return nil, ""
}
