package provider

import (
	"context"
	"errors" 

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

// resourceUser defines the schema and CRUD functions for the sotoon_iam_user resource.
func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages an IAM user within a Sotoon workspace.",
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the user, returned by the API.",
			},
			"user_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the user.",
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The email address of the user. Must be unique within the workspace.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the user.",
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	email := d.Get("email").(string)


	_, err := c.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			tflog.Debug(ctx, "User not found, sending invitation", map[string]interface{}{"email": email})

			_, inviteErr := c.InviteUser(ctx, email)
			if inviteErr != nil {
				if strings.Contains(strings.ToLower(inviteErr.Error()), "forbidden") {
					return diag.Errorf("Forbidden: The API token does not have permission to invite new users.")
				}
				return diag.Errorf("Failed to invite user with email '%s': %s", email, inviteErr)
			}
		} else if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return diag.Errorf("Forbidden: The API token does not have permission to read user information. Cannot check if user '%s' exists.", email)
		} else {
			return diag.FromErr(err)
		}
	} else {
		tflog.Debug(ctx, "User already exists, adopting into state", map[string]interface{}{"email": email})
	}

	d.SetId(email)

	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	userEmail := d.Id()

	user, err := c.GetUserByEmail(ctx, userEmail)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := d.Set("email", user.Email); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", user.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("user_uuid", user.UUID.String()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// resourceUserDelete "disowns" the user from the state without deleting it from Sotoon.
func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
