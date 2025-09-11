package iamresources

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sotoon/iam-client/pkg/types"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

// ResourceGroup defines the schema and CRUD functions for the sotoon_iam_group resource.
func ResourceGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages an IAM group within a Sotoon workspace.",
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		DeleteContext: resourceGroupDelete,
		UpdateContext: resourceGroupUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the group, returned by the API.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the group.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the group.",
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	groups, err := c.GetWorkspaceGroups(ctx, c.WorkspaceUUID)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, g := range groups {
		if g.Name == name {
			d.SetId(g.UUID.String())
			return resourceGroupRead(ctx, d, meta)
		}
	}

	created, err := c.CreateGroup(ctx, name, description)
	if err != nil {
		return diag.Errorf("Failed to create group %q: %s", name, err.Error())
	}

	d.SetId(created.UUID.String())

	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	groups, err := c.GetWorkspaceGroups(ctx, c.WorkspaceUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	var found *types.Group
	for _, g := range groups {
		if g.UUID.String() == d.Id() {
			found = g
			break
		}
	}
	if found == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", found.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", found.Description); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	if d.HasChange("description") {
		name := d.Get("name").(string)
		description := d.Get("description").(string)
		err := c.UpdateGroup(ctx, d.Id(), name, description)
		if err != nil {
			return diag.Errorf("failed to update group description for %q: %s", d.Id(), err)
		}
	}

	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	groupID := d.Id()

	if err := c.DeleteGroup(ctx, groupID); err != nil {
		if errors.Is(err, client.ErrNotFound) {
			tflog.Warn(ctx, "Group already deleted or not found", map[string]interface{}{"group_id": groupID})
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
