package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	uuid "github.com/satori/go.uuid"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func resourceServiceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceUserCreate,
		ReadContext:   resourceServiceUserRead,
		DeleteContext: resourceServiceUserDelete,
		UpdateContext: resourceServiceUserUpdate,
		Schema: map[string]*schema.Schema{
			"id":          {Type: schema.TypeString, Computed: true},
			"name":        {Type: schema.TypeString, Required: true},
			"description": {Type: schema.TypeString, Optional: true},
		},
	}
}

func resourceServiceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	serviceUser, err := c.CreateServiceUser(name, description, c.WorkspaceUUID)
	if err != nil {
		return diag.FromErr(err)
	}
	if serviceUser == nil || serviceUser.UUID == nil {
		return diag.Errorf("empty service user response")
	}

	d.SetId(serviceUser.UUID.String())
	return resourceServiceUserRead(ctx, d, meta)
}

func resourceServiceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id := d.Id()
	u, err := uuid.FromString(id)
	if err != nil {
		return diag.FromErr(err)
	}
	su, err := c.GetServiceUser(c.WorkspaceUUID, &u)
	if err != nil {
		return diag.FromErr(err)
	}
	if su == nil || su.UUID == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", su.Name)
	return nil
}

func resourceServiceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	if d.HasChange("name") || d.HasChange("description") {
		id := d.Id()
		u, err := uuid.FromString(id)
		if err != nil {
			return diag.FromErr(err)
		}
		name := d.Get("name").(string)
		desc := ""
		if v, ok := d.GetOk("description"); ok {
			desc = v.(string)
		}
		if _, err := c.UpdateServiceUser(*c.WorkspaceUUID, u, name, desc); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceServiceUserRead(ctx, d, meta)
}

func resourceServiceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id := d.Id()
	u, err := uuid.FromString(id)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.DeleteServiceUser(c.WorkspaceUUID, &u); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
