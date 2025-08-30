package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	uuid "github.com/satori/go.uuid"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func resourceServiceUserPublicKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceUserPublicKeyCreateCreate,
		ReadContext:   resourceServiceUserPublicKeyCreateRead,
		DeleteContext: resourceServiceUserPublicKeyCreateDelete,
		Schema: map[string]*schema.Schema{
			"id":              {Type: schema.TypeString, Computed: true},
			"service_user_id": {Type: schema.TypeString, Required: true, ForceNew: true},
			"title":           {Type: schema.TypeString, Required: true, ForceNew: true},
			"public_key":      {Type: schema.TypeString, Required: true, Sensitive: true, ForceNew: true},
		},
	}
}

func resourceServiceUserPublicKeyCreateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	serviceUserID := d.Get("service_user_id").(string)
	title := d.Get("title").(string)
	key := d.Get("public_key").(string)

	serviceUserUUID, err := uuid.FromString(serviceUserID)
	if err != nil {
		return diag.FromErr(err)
	}
	pk, err := c.CreateServiceUserPublicKey(*c.WorkspaceUUID, serviceUserUUID, title, key)
	if err != nil {
		return diag.FromErr(err)
	}
	if pk == nil || pk.UUID == nil {
		return diag.Errorf("empty public key response")
	}
	d.SetId(fmt.Sprintf("%s/%s", serviceUserUUID.String(), pk.UUID.String()))
	return resourceServiceUserPublicKeyCreateRead(ctx, d, meta)
}

func resourceServiceUserPublicKeyCreateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	suID, pkID, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	list, err := c.GetWorkspaceServiceUserPublicKeyList(*c.WorkspaceUUID, suID)
	if err != nil {
		return diag.FromErr(err)
	}
	found := false
	for _, pk := range list {
		if pk != nil && pk.UUID != nil && pk.UUID.String() == pkID.String() {
			found = true
			break
		}
	}
	if !found {
		d.SetId("")
		return nil
	}
	d.Set("service_user_id", suID.String())
	return nil
}

func resourceServiceUserPublicKeyCreateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	suID, pkID, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.DeleteServiceUserPublicKey(*c.WorkspaceUUID, suID, pkID); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
