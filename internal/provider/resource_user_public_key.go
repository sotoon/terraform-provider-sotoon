package provider

import (
	"context"
	"fmt"
	"strings"

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
	key := d.Get("public_key").(string)

	keysList, err := c.GetAllMyUserPublicKeyList(ctx)
	if err != nil {
		return diag.Errorf("failed to list all public keys: %s", err)
	}
	for _, keyInfo := range keysList {
		if strings.Contains(key, keyInfo.Key) {
			if keyInfo.Title == title {
				d.SetId(keyInfo.Uuid)
				if err := d.Set("title", keyInfo.Title); err != nil {
					return diag.FromErr(fmt.Errorf("failed to set title: %w", err))
				}
				if err := d.Set("key_type", keyInfo.Type); err != nil {
					return diag.FromErr(fmt.Errorf("failed to set key_type: %w", err))
				}
				return nil
			}
			return diag.Errorf("this key has already been registered by another title")
		}
	}

	created, err := c.CreateMyUserPublicKey(ctx, title, key)
	if err != nil {
		return diag.Errorf("Failed to create user public key %q: %s", title, err.Error())
	}

	d.SetId(created.Uuid)
	if err := d.Set("key_type", created.Type); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set key_type: %w", err))
	}
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

	key, err := c.GetUserPublicKey(ctx, &uid)
	if err != nil {
		return diag.Errorf("error reading public-key %s: %s", id, err)
	}

	if err := d.Set("title", key.Title); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set title: %w", err))
	}
	if err := d.Set("key_type", key.Type); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set key_type: %w", err))
	}
	if err := d.Set("public_key", key.PublicKey); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set public_key: %w", err))
	}
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

	err = c.DeleteUserPublicKey(ctx, &uid)
	if err != nil {
		return diag.Errorf("error deleting public-key %q: %s", id, err)
	}

	d.SetId("")
	return nil
}
