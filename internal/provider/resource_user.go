// internal/provider/resource_user.go

package provider

import (
	"context"
	"errors" // Import the errors package

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sotoon/terraform-provider-sotoon/internal/client" // Ensure this module path is correct
	"github.com/sotoon/iam-client/pkg/types"                      // Import types to inspect the custom error
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
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The email address of the user. Must be unique within the workspace.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The display name of the user.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "The password for the new user.",
			},
		},
	}
}

// resourceUserCreate creates a new user based on the Terraform plan.
func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	userRequest := client.CreateUserRequest{
		Email:    d.Get("email").(string),
		Name:     d.Get("name").(string),
		Password: d.Get("password").(string),
	}

	createdUser, err := c.CreateUser(ctx, userRequest)
	if err != nil {
		// --- ADVANCED ERROR HANDLING ---
		// Check if the error is the specific type from the iam-client library.
		var reqErr *types.RequestExecutionError
		if errors.As(err, &reqErr) {
			// If it is, create a more detailed diagnostic for the user.
			// This prevents the provider from crashing on the unmarshal error
			// by showing the raw API response instead.
			return diag.Errorf(
				"API Error: Failed to create user with status code %d. Raw response: %s",
				reqErr.StatusCode,
				string(reqErr.Data),
			)
		}
		// For other types of errors, use the default handler.
		return diag.FromErr(err)
	}

	d.SetId(createdUser.UUID.String())

	return resourceUserRead(ctx, d, meta)
}

// resourceUserRead retrieves the user's information from the API.
func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	userID := d.Id()

	user, err := c.GetUserByID(ctx, userID)
	if err != nil {
		// Check if the error is our custom "not found" error from the client wrapper.
		if errors.Is(err, client.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if user == nil {
		d.SetId("")
		return nil
	}

	d.Set("email", user.Email)
	d.Set("name", user.Name)

	return nil
}

// resourceUserDelete removes the user from the platform.
func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	userID := d.Id()

	err := c.DeleteUser(ctx, userID)
	if err != nil {
		// Also ignore "not found" errors on delete, as the resource is already gone.
		if errors.Is(err, client.ErrNotFound) {
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
