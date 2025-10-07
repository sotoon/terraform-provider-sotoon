package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	uuid "github.com/satori/go.uuid"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func resourceServiceUser() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a service user within a Sotoon workspace.",
		CreateContext: resourceServiceUserCreate,
		ReadContext:   resourceServiceUserRead,
		DeleteContext: resourceServiceUserDelete,
		UpdateContext: resourceServiceUserUpdate,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Composite stable identifier. Does not affect lifecycle.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the service user.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the service user.",
			},
		},
	}
}

func resourceServiceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	name := d.Get("name").(string)
	description := d.Get("description").(string)

	currentServiceUsers, err := c.GetServiceUsers(ctx)
	if err != nil {
		return diag.Errorf("failed to get services in resourceServiceUserCreate %s", err.Error())
	}

	for _, su := range currentServiceUsers {
		if su.Name == name {
			d.Set("description", su.Description)
			d.SetId(su.Uuid)
			return nil
		}
	}

	serviceUser, err := c.CreateServiceUser(ctx, name, description)
	if err != nil {
		return diag.FromErr(err)
	}
	if serviceUser == nil || serviceUser.Uuid == "" {
		return diag.Errorf("empty service user response")
	}

	d.SetId(serviceUser.Uuid)
	return resourceServiceUserRead(ctx, d, meta)
}

func resourceServiceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id := d.Id()
	u, err := uuid.FromString(id)
	if err != nil {
		return diag.FromErr(err)
	}
	su, err := c.GetServiceUser(ctx, &u)
	if err != nil {
		return diag.FromErr(err)
	}
	if su == nil || su.Uuid == "" {
		d.SetId("")
		return nil
	}
	if err := d.Set("name", su.Name); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set name: %w", err))
	}
	if err := d.Set("description", su.Description); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set description: %w", err))
	}
	return nil
}

func resourceServiceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id := d.Id()
	u, err := uuid.FromString(id)
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)
	desc := d.Get("description").(string)
	res, err := c.UpdateServiceUser(ctx, u, name, desc)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", res.Name); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set name: %w", err))
	}
	if err := d.Set("description", res.Description); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set description: %w", err))
	}
	return nil
}

func resourceServiceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id := d.Id()
	u, err := uuid.FromString(id)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.DeleteServiceUser(ctx, &u); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
