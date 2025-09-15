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
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Composite stable identifier. Does not affect lifecycle.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Group UUID.",
			},
			"service_user_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Set of Service User UUIDs to bind to the group.",
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
	groupID := d.Get("group_id").(string)
	groupUUID, err := uuid.FromString(groupID)
	if err != nil {
		return diag.Errorf("invalid group_id: %s", err)
	}

	sortedServiceUserIds := uniqueSorted(fromSchemaSetToStrings(d.Get("service_user_ids").(*schema.Set)))

	serviceUsersList, err := c.GetAllGroupServiceUserList(context.Background(), c.WorkspaceUUID, &groupUUID)
	if err != nil {
		return diag.Errorf("read group service-users: %s", err)
	}

	remoteServiceUsersID := make([]string, 0, len(serviceUsersList))

	for _, u := range serviceUsersList {
		remoteServiceUsersID = append(remoteServiceUsersID, u.UUID.String())
	}

	remoteServiceUsersID = uniqueSorted(remoteServiceUsersID)
	toAddList := diff(toSet(sortedServiceUserIds), toSet(remoteServiceUsersID))

	if len(toAddList) > 0 {
		uuids := make([]uuid.UUID, 0, len(toAddList))
		for _, id := range toAddList {
			uuid, _ := uuid.FromString(id)
			uuids = append(uuids, uuid)
		}

		if _, err := c.BulkAddServiceUsersToGroup(*c.WorkspaceUUID, groupUUID, uuids); err != nil {
			return diag.Errorf("add service users to group %s: %s", groupID, err)
		}
	}

	bindHash := hashOfIDs(sortedServiceUserIds)
	d.Set("bindings_hash", bindHash)
	d.SetId(groupUUID.String() + ":" + bindHash)

	return resourceGroupServiceUserRead(ctx, d, meta)
}

func resourceGroupServiceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	groupID := d.Get("group_id").(string)
	groupUUID, err := uuid.FromString(groupID)
	if err != nil {
		d.SetId("")
		return nil
	}

	sortedServiceUserIds := uniqueSorted(fromSchemaSetToStrings(d.Get("service_user_ids").(*schema.Set)))

	serviceUsersList, err := c.GetAllGroupServiceUserList(context.Background(), c.WorkspaceUUID, &groupUUID)
	if err != nil {
		return diag.Errorf("read group service-users: %s", err)
	}

	remoteServiceUsersID := make([]string, 0, len(serviceUsersList))
	for _, u := range serviceUsersList {
		remoteServiceUsersID = append(remoteServiceUsersID, u.UUID.String())
	}
	remoteServiceUsersID = uniqueSorted(remoteServiceUsersID)

	eff := intersect(toSet(sortedServiceUserIds), toSet(remoteServiceUsersID))
	effective := uniqueSorted(setKeys(eff))

	d.Set("service_user_ids", effective)

	bindHash, _ := d.Get("bindings_hash").(string)
	if bindHash == "" {
		bindHash = hashOfIDs(effective)
		d.Set("bindings_hash", bindHash)
	}

	d.SetId(groupUUID.String() + ":" + bindHash)

	tflog.Info(ctx, "Read service-user group", map[string]interface{}{"group_id": groupID, "have": len(effective)})

	return nil
}

func resourceGroupServiceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	groupID := d.Get("group_id").(string)
	groupUUID, err := uuid.FromString(groupID)
	if err != nil {
		return diag.Errorf("invalid group_id %q: %s", groupID, err)
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
