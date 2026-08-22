package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iam "github.com/sotoon/sotoon-sdk-go/sdk/core/iam_v1"
	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

type userWorkspaceRemover interface {
	GetUserByEmail(context.Context, string) (*iam.IamUser, error)
	RemoveUserFromWorkspace(context.Context, string) error
}

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
			_, inviteErr := c.InviteUser(ctx, email)
			if inviteErr != nil {
				return diag.Errorf("Failed to invite user with email '%s': %s", email, inviteErr)
			}
		} else {
			return diag.FromErr(err)
		}
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
	if err := d.Set("user_uuid", user.Uuid); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceUserDeleteWithClient(ctx, d, meta.(*client.Client))
}

func resourceUserDeleteWithClient(ctx context.Context, d *schema.ResourceData, c userWorkspaceRemover) diag.Diagnostics {
	user, err := c.GetUserByEmail(ctx, d.Id())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to look up user %q before removal: %s", d.Id(), err)
	}

	if user.Uuid == "" {
		return diag.Errorf("cannot remove user %q from workspace: API returned an empty UUID", d.Id())
	}
	if err := c.RemoveUserFromWorkspace(ctx, user.Uuid); err != nil {
		return diag.Errorf("failed to remove user %q from workspace: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
