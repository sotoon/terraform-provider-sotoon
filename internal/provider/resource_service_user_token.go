package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	uuid "github.com/satori/go.uuid"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func resourceServiceUserToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceUserTokenCreate,
		ReadContext:   resourceServiceUserTokenRead,
		DeleteContext: resourceServiceUserTokenDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   false,
				Description: "The newly issued service user token value.",
			},
		},
	}
}

func resourceServiceUserTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	idStr := d.Get("service_user_id").(string)
	suID, err := uuid.FromString(idStr)
	if err != nil {
		return diag.FromErr(err)
	}

	tok, err := c.CreateServiceUserToken(&suID, c.WorkspaceUUID)
	if err != nil {
		return diag.FromErr(err)
	}
	if tok == nil || tok.UUID == nil {
		return diag.Errorf("empty token response")
	}

	if tok.Secret == "" {
		return diag.Errorf("service user token secret not returned by API at creation time")
	}
	if err := d.Set("value", tok.Secret); err != nil {
		return diag.Errorf("error setting token value: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", suID.String(), tok.UUID.String()))
	return resourceServiceUserTokenRead(ctx, d, meta)
}

func resourceServiceUserTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	suID, tokID, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	list, err := c.GetWorkspaceServiceUserTokenList(&suID, c.WorkspaceUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	found := false
	if list != nil {
		for _, t := range *list {
			if t.UUID != nil && t.UUID.String() == tokID.String() {
				found = true
				break
			}
		}
	}
	if !found {
		d.SetId("")
		return nil
	}

	_ = d.Set("service_user_id", suID.String())

	return nil
}

func resourceServiceUserTokenDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	suID, tokID, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.DeleteServiceUserToken(&suID, c.WorkspaceUUID, &tokID); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
