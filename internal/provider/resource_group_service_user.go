package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"

	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func resourceGroupServiceUser() *schema.Resource {
	return &schema.Resource{
		Description:   "Binds service users to a group within a Sotoon workspace.",
		CreateContext: resourceGroupServiceUserCreate,
		ReadContext:   resourceGroupServiceUserRead,
		DeleteContext: resourceGroupServiceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_user_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceGroupServiceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	groupStr := d.Get("group_id").(string)
	groupUUID, err := uuid.FromString(groupStr)
	if err != nil {
		return diag.Errorf("invalid group_id format: %s", err)
	}

	raw := d.Get("service_user_ids").(*schema.Set).List()

	suUUIDs := make([]uuid.UUID, 0, len(raw))
	for _, v := range raw {
		s := v.(string)
		u, err := uuid.FromString(s)
		if err != nil {
			return diag.Errorf("invalid service_user_id in list: %s", err)
		}
		suUUIDs = append(suUUIDs, u)
	}

	if len(suUUIDs) > 0 {
		if _, err := c.IAMClient.BulkAddServiceUsersToGroup(*c.WorkspaceUUID, groupUUID, suUUIDs); err != nil {
			return diag.Errorf("failed to add service users to group %s: %s", groupUUID, err)
		}
	}

	d.SetId(groupUUID.String() + ":" + uuid.NewV4().String())
	return resourceGroupServiceUserRead(ctx, d, meta)
}

func resourceGroupServiceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Reading service user group", map[string]interface{}{"id": d.Id()})
	return nil
}

func resourceGroupServiceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	groupStr := d.Get("group_id").(string)
	groupUUID, err := uuid.FromString(groupStr)
	if err != nil {
		return diag.Errorf("invalid group_id %q: %s", groupStr, err)
	}

	serviceUserIds := d.Get("service_user_ids").(*schema.Set).List()

	for _, v := range serviceUserIds {
		s := v.(string)
		u, err := uuid.FromString(s)
		if err != nil {
			return diag.Errorf("invalid service_user_id in list: %s", err)
		}

		if err := c.UnbindServiceUserFromGroup(c.WorkspaceUUID, &groupUUID, &u); err != nil {
			return diag.Errorf("unbind service user %s from group %s failed: %s", u.String(), groupUUID.String(), err)
		}
	}

	d.SetId("")
	return nil
}

func uuidsToStringSlice(in []uuid.UUID) []string {
	out := make([]string, len(in))
	for i, u := range in {
		out[i] = u.String()
	}
	return out
}
