package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "Fetches a single IAM user within a specific Sotoon workspace.",
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the workspace to fetch users from.",
			},
			"uuid": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The UUID of the user to fetch.",
				ExactlyOneOf: []string{"uuid", "email"},
				Computed:     true,
			},
			"email": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The email of the user to fetch.",
				ExactlyOneOf: []string{"uuid", "email"},
				Computed:     true,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the user.",
			},
			"is_suspended": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the user account is suspended.",
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	workspaceID := d.Get("workspace_id").(string)

	workspaceUUID, err := uuid.FromString(workspaceID)
	if err != nil {
		return diag.Errorf("Invalid workspace_id format: not a valid UUID")
	}

	tflog.Debug(ctx, "Reading users for workspace", map[string]interface{}{"workspace_id": workspaceID})

	userData := make(map[string]interface{})

	if uuid, uuidFound := d.GetOk("uuid"); uuidFound {
		user, err := c.GetWorkspaceUserByUUID(ctx, &workspaceUUID, uuid.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		userData = map[string]interface{}{
			"uuid":         *user.Uuid,
			"email":        *user.Email,
			"name":         *user.Name,
			"is_suspended": *user.IsSuspended,
		}

	}

	if email, emailFound := d.GetOk("email"); emailFound {
		user, err := c.GetWorkspaceUserByEmail(ctx, &workspaceUUID, email.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		userData = map[string]interface{}{
			"uuid":         user.Uuid,
			"email":        user.Email,
			"name":         user.Name,
			"is_suspended": user.IsSuspended,
		}
	}

	if err := d.Set("uuid", userData["uuid"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email", userData["email"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", userData["name"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_suspended", userData["is_suspended"]); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("user_" + userData["uuid"].(string)) 

	return nil
}
