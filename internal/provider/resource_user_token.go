package provider

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

// resourceUserToken defines the schema and CRUD functions for the sotoon_iam_user_token resource.
func resourceUserToken() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a user token for the current IAM user.",
		CreateContext: resourceUserTokenCreate,
		ReadContext:   resourceUserTokenRead,
		DeleteContext: resourceUserTokenDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the token.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name/label for the newly minted user token.",
			},
			"expire_at": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Expiration timestamp in RFC3339 format (e.g. 2025-09-30T00:00:00Z).",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   false,
				Description: "The newly issued token value.",
			},
		},
	}
}

func resourceUserTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	name := d.Get("name").(string)

	var expiresAt *time.Time
	if v, ok := d.GetOk("expire_at"); ok {
		t, err := time.Parse(time.RFC3339, v.(string))
		if err != nil {
			return diag.Errorf("invalid expire_at format: %s", err)
		}
		expiresAt = &t
	}

	tflog.Debug(ctx, "Creating user token", map[string]interface{}{
		"name":      name,
		"expire_at": expiresAt,
	})

	created, err := c.CreateMyUserToken(name, expiresAt)
	if err != nil {
		return diag.Errorf("Failed to create user token %q: %s", name, err.Error())
	}

	d.SetId(created.UUID)

	if err := d.Set("value", created.Secret); err != nil {
		return diag.Errorf("error setting token value: %s", err)
	}

	return resourceUserTokenRead(ctx, d, meta)
}

func resourceUserTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	id := d.Id()

	tflog.Debug(ctx, "Reading user token", map[string]interface{}{"token_id": id})

	tid, err := uuid.FromString(id)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = c.GetMyUserToken(&tid)
	if err != nil {
		return diag.Errorf("error reading token %s: %s", id, err)
	}

	return nil
}

func resourceUserTokenDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	id := d.Id()

	tflog.Debug(ctx, "Deleting user token", map[string]interface{}{"token_id": id})

	tid, err := uuid.FromString(id)
	if err != nil {
		d.SetId("")
		return nil
	}

	err = c.DeleteMyUserToken(&tid)
	if err != nil {
		return diag.Errorf("error deleting token %q: %s", id, err)
	}

	d.SetId("")
	return nil
}
