package iamdatasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	uuid "github.com/satori/go.uuid"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func DataSourceServiceUserDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServiceUserDetailsRead,
		Schema: map[string]*schema.Schema{
			"service_user_id": {Type: schema.TypeString, Required: true},
			"id":              {Type: schema.TypeString, Computed: true},
			"name":            {Type: schema.TypeString, Computed: true},
		},
	}
}

func dataSourceServiceUserDetailsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	idStr := d.Get("service_user_id").(string)
	suID, err := uuid.FromString(idStr)
	if err != nil {
		return diag.FromErr(err)
	}

	det, err := c.GetWorkspaceServiceUserDetail(*c.WorkspaceUUID, suID)
	if err != nil {
		return diag.FromErr(err)
	}
	if det == nil || det.UUID == nil {
		d.SetId("")
		return nil
	}

	d.SetId(det.UUID.String())
	d.Set("name", det.Name)
	return nil
}
