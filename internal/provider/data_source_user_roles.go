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
	c := meta.(*client.Client)

	userStr := d.Get("user_id").(string)
	userUUID, err := uuid.FromString(userStr)
	if err != nil {
		return diag.Errorf("invalid user_id %q: %s", userStr, err)
	}

	u, err := c.GetUserDetailed(ctx, &userUUID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get user detail (roles): %w", err))
	}
	if u.Uuid == nil {
		return diag.Errorf("user %s not found", userUUID.String())
	}

	roleList := make([]map[string]interface{}, 0, len(u.Roles))
	for _, r := range u.Roles {
		roleList = append(roleList, map[string]interface{}{
			"id":   r.Uuid,
			"name": r.Name,
		})
	}
	if err := d.Set("roles", roleList); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set roles: %w", err))
	}

	anchor := userUUID.String()
	if c.WorkspaceUUID != nil {
		anchor = fmt.Sprintf("%s:%s", c.WorkspaceUUID.String(), userUUID.String())
	}

	d.SetId("user-roles:" + anchor)
	return nil
}
