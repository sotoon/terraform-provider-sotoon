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
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"bindings_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SHA-256 of sorted, canonical service_user_ids.",
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

	input := fromSchemaSetToStrings(d.Get("service_user_ids").(*schema.Set))
	converted, err := convertStringsToUUIDArray(input)
	if err != nil {
		return diag.Errorf("invalid service_user_ids: %s", err)
	}
	ids := uniqueSorted(converted)
	_ = d.Set("bindings_hash", hashOfIDs(ids))

	suUUIDs := make([]uuid.UUID, 0, len(ids))
	for _, s := range ids {
		u, _ := uuid.FromString(s)
		suUUIDs = append(suUUIDs, u)
	}

	if len(suUUIDs) > 0 {
		if _, err := c.IAMClient.BulkAddServiceUsersToGroup(*c.WorkspaceUUID, groupUUID, suUUIDs); err != nil {
			return diag.Errorf("failed to add service users to group %s: %s", groupUUID, err)
		}
	}

	d.SetId(groupUUID.String() + ":" + hashOfIDs(ids)[:16])
	return resourceGroupServiceUserRead(ctx, d, meta)
}

func resourceGroupServiceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	input := fromSchemaSetToStrings(d.Get("service_user_ids").(*schema.Set))
	converted, err := convertStringsToUUIDArray(input)
	if err == nil {
		_ = d.Set("bindings_hash", hashOfIDs(uniqueSorted(converted)))
	}
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
