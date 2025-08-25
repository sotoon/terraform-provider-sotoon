package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"

	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func dataSourceUserRoles() *schema.Resource {
	return &schema.Resource{
		Description: "Lists all roles bound to a specific (human) user.",
		ReadContext: dataSourceUserRolesRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User UUID.",
			},
			"roles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Roles bound to the user.",
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

func dataSourceUserRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(client.Client) // interface is fine; no need to cast to *client.Client

	userStr := d.Get("user_id").(string)
	userUUID, err := uuid.FromString(userStr)
	if err != nil {
		return diag.Errorf("invalid user_id %q: %s", userStr, err)
	}

	u, err := c.GetUser(&userUUID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get user detail (roles): %w", err))
	}
	if u == nil || u.UUID == nil {
		return diag.Errorf("user %s not found", userUUID.String())
	}

	roleList := make([]map[string]interface{}, 0, len(u.Roles))
	for _, r := range u.Roles {
		var rid string
		if r.UUID != nil {
			rid = r.UUID.String()
		}
		roleList = append(roleList, map[string]interface{}{
			"id":   rid,
			"name": r.Name,
		})
	}

	if err := d.Set("roles", roleList); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set roles: %w", err))
	}

	d.SetId(userUUID.String())
	return nil
}
