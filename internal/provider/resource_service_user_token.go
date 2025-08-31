package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	uuid "github.com/satori/go.uuid"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func resourceServiceUserToken() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a token for a service user within a Sotoon workspace.",
		CreateContext: resourceServiceUserTokenCreate,
		ReadContext:   resourceServiceUserTokenRead,
		DeleteContext: resourceServiceUserTokenDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Composite stable identifier. Does not affect lifecycle.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Name of the token.",
			},
			"expires_at": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Expiration time of the token in RFC3339 format.",
			},
			"service_user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Service User UUID.",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The newly issued service user token value",
			},
		},
	}
}

func resourceServiceUserTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	serviceUserID := d.Get("service_user_id").(string)
	serviceUserUUID, err := uuid.FromString(serviceUserID)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)

	var expiresAt *time.Time
	if expiresAtString, ok := d.GetOk("expires_at"); ok {
		expAt, err := time.Parse(time.RFC3339, expiresAtString.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		expiresAt = &expAt
	}

	tok, err := c.CreateServiceUserToken(&serviceUserUUID, c.WorkspaceUUID, name, expiresAt)
	if err != nil {
		return diag.FromErr(err)
	}
	if tok == nil || tok.UUID == nil || tok.Secret == "" {
		return diag.Errorf("empty token response")
	}

	if err := d.Set("value", tok.Secret); err != nil {
		return diag.Errorf("error setting token value: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceUserUUID.String(), tok.UUID.String()))
	return resourceServiceUserTokenRead(ctx, d, meta)
}

func resourceServiceUserTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	serviceUserID, tokenID, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	list, err := c.GetWorkspaceServiceUserTokenList(&serviceUserID, c.WorkspaceUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	if list != nil {
		for _, t := range *list {
			if t.UUID != nil && t.UUID.String() == tokenID.String() {
				d.Set("service_user_id", serviceUserID.String())
				d.Set("name", t.Name)
				if t.ExpiresAt != nil {
					d.Set("expires_at", t.ExpiresAt.Format(time.RFC3339))
				}
				return nil
			}
		}
	}

	d.SetId("")
	return nil
}

func resourceServiceUserTokenDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	serviceUserID, tokenID, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.DeleteServiceUserToken(&serviceUserID, c.WorkspaceUUID, &tokenID); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
