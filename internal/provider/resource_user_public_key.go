package provider

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func resourceUserPublicKey() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages an SSH public-key for the current IAM user.",
		CreateContext: resourceUserPublicKeyCreate,
		ReadContext:   resourceUserPublicKeyRead,
		DeleteContext: resourceUserPublicKeyDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the public‐key record.",
			},
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A friendly title for this key.",
			},
			"key_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The key type (e.g. ssh-rsa, ssh-ed25519).",
			},
			"public_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The SSH public key material.",
			},
		},
	}
}

func resourceUserPublicKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	title := d.Get("title").(string)
	keyType := d.Get("key_type").(string)
	key := d.Get("public_key").(string)

	tflog.Debug(ctx, "Creating public key", map[string]interface{}{"title": title})

	created, err := c.CreateMyUserPublicKey(title, keyType, key)
	if err != nil {
		return diag.Errorf("Failed to create user public key %q: %s", title, err.Error())
	}

	d.SetId(created.UUID)
	return resourceUserPublicKeyRead(ctx, d, meta)
}

func resourceUserPublicKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	id := d.Id()

	tflog.Debug(ctx, "Reading public key", map[string]interface{}{"id": id})

	uid, err := uuid.FromString(id)
	if err != nil {
		d.SetId("")
		return nil
	}

	key, err := c.GetOneDefaultUserPublicKey(&uid)
	if err != nil {
		return diag.Errorf("error reading public-key %s: %s", id, err)
	}

	d.Set("title", key.Title)
	d.Set("key_type", key.Type)
	d.Set("public_key", key.PublicKey)
	return nil
}

func resourceUserPublicKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	id := d.Id()

	tflog.Debug(ctx, "Deleting public key", map[string]interface{}{"id": id})

	uid, err := uuid.FromString(id)
	if err != nil {
		d.SetId("")
		return nil
	}

	err = c.DeleteDefaultUserPublicKey(&uid)
	if err != nil {
		return diag.Errorf("error deleting public-key %q: %s", id, err)
	}

	d.SetId("")
	return nil
}
