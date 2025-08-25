package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"

	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func dataSourceServiceUserRoles() *schema.Resource {
	return &schema.Resource{
		Description: "Lists all roles bound to a specific service user within a Sotoon workspace.",
		ReadContext: dataSourceServiceUserRolesRead,
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Workspace UUID.",
			},
			"service_user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Service user UUID.",
			},
			"roles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Roles bound to the service user in the workspace.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Role UUID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Role name.",
						},
					},
				},
			},
		},
	}
}

func dataSourceServiceUserRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	wsStr := d.Get("workspace_id").(string)
	suStr := d.Get("service_user_id").(string)

	wsUUID, err := uuid.FromString(wsStr)
	if err != nil {
		return diag.Errorf("invalid workspace_id %q: %s", wsStr, err)
	}
	suUUID, err := uuid.FromString(suStr)
	if err != nil {
		return diag.Errorf("invalid service_user_id %q: %s", suStr, err)
	}

	detail, err := c.GetWorkspaceServiceUserDetail(wsUUID, suUUID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get service user detail (roles): %w", err))
	}
	if detail == nil {
		return diag.Errorf("service user %s not found in workspace %s", suUUID.String(), wsUUID.String())
	}

	roleList := make([]map[string]interface{}, 0)
	for _, r := range detail.Roles {
		roleList = append(roleList, map[string]interface{}{
			"id":   r.UUID.String(),
			"name": r.Name,
		})
	}

	if err := d.Set("roles", roleList); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set roles: %w", err))
	}

	d.SetId(wsUUID.String() + ":" + suUUID.String())
	return nil
}
