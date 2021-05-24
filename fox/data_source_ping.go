package fox

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"strconv"
	"terraform-provider-fox/pkg/common"
	"time"
)

func dataSourcePing() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePingRead,
		Schema: map[string]*schema.Schema{
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(common.Config)

	url, err := common.GetApiUrl(config.Address, "ping")
	if err != nil {
		return diag.FromErr(err)
	}
	status, respBody, err := common.HttpRequest(&common.HttpRequestArgs{
		Method:     http.MethodGet,
		Url:        url,
		TimeoutSec: 10,
	})
	if err != nil {
		return diag.FromErr(err)
	} else if status != http.StatusOK {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "http response not ok",
			Detail:   fmt.Sprintf("status: %d", status),
		})
		return diags
	}

	pong := make(map[string]interface{})
	err = json.Unmarshal(respBody, &pong)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "json decode fail",
			Detail:   fmt.Sprintf("%s\n%s", err.Error(), string(respBody)),
		})
		return diags
	}

	if err := d.Set("message", pong["message"]); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "schema set fail",
			Detail:   err.Error(),
		})
		return diags
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
