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

func dataSourceAll() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAllRead,
		Schema: map[string]*schema.Schema{
			"ip_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceAllRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(common.Config)

	// get data
	ipInfoArr, err := ip.GetIpInfoAll(config, &diags)
	if err != nil {
		return diags
	}

	// set data
	err, subErrMsg := bindDataSrcAll(d, ipInfoArr)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("bind ipInfoArr to data resource fail: %s", subErrMsg),
			Detail:   fmt.Sprintf("ipInfoArr: %+v\nerror: %s", ipInfoArr, err.Error()),
		})
		return diags
	}

	// always run
	d.SetId("all")

	return diags
}

func bindDataSrcAll(d *schema.ResourceData, ipInfoArr *[]ip.IpInfo) (error, string) {
	if ipInfoArr == nil {
		err := errors.New("bind data fail")
		return err, "ipInfoArr is null"
	}

	const subErrMsgFormat = "%s set fail"
	ois := make([]interface{}, len(*ipInfoArr), len(*ipInfoArr))
	for i, ipInfo := range *ipInfoArr {
		oi := make(map[string]interface{})
		oi["id"] = ipInfo.Env
		oi["env"] = ipInfo.Env
		oi["groups"] = ipInfo.GroupArr
		oi["cidrs"] = ipInfo.CidrArr
		oi["created"] = ipInfo.Created.Format(time.RFC3339)
		oi["last_modified"] = ipInfo.LastModified.Format(time.RFC3339)
		ois[i] = oi
	}

	if err := d.Set("ip_infos", ois); err != nil {
		return err, fmt.Sprintf(subErrMsgFormat, "ip_infos")
	}

	return nil, ""
}
