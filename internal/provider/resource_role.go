package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"

	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages an IAM role within a Sotoon workspace.",
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the role.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the role.",
			},
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	name := d.Get("name").(string)

	existing, err := c.GetRoleByName(ctx, name)
	if err == nil {
		d.SetId(existing.UUID.String())
		return resourceRoleRead(ctx, d, meta)
	}

	created, err := c.CreateRole(ctx, name)
	if err != nil {
		return diag.Errorf("failed to create role %q: %s", name, err)
	}
	d.SetId(created.UUID.String())
	return resourceRoleRead(ctx, d, meta)
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	id := d.Id()
	ruleUUID, err := uuid.FromString(id)
	if err != nil {
		d.SetId("")
		return nil
	}

	res, err := c.GetRole(ctx, &ruleUUID)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("name", res.Name)
	return nil
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id := d.Id()
	roleUUID, err := uuid.FromString(id)
	if err != nil {
		return diag.Errorf("invalid role ID %q: %s", id, err)
	}

	if d.HasChange("name") {
		newName := d.Get("name").(string)
		updated, err := c.UpdateRole(ctx, &roleUUID, newName)
		if err != nil {
			return diag.Errorf("failed to update role %q: %s", id, err)
		}
		d.SetId(updated.UUID.String())
	}

	return resourceRoleRead(ctx, d, meta)
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	if err := c.DeleteRole(ctx, d.Id()); err != nil {
		if errors.Is(err, client.ErrNotFound) {
			tflog.Warn(ctx, "Role already deleted or not found", map[string]interface{}{"role_id": d.Id()})
			return nil
		}
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
