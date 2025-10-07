package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"

	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages an IAM role within a Sotoon workspace.",
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the role.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the role.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the role.",
			},
			"rules": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of rule UUIDs to attach to this role.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	name := d.Get("name").(string)

	existing, err := c.GetRoleByName(ctx, name)
	if err == nil {
		d.SetId(existing.Uuid)
		return resourceRoleRead(ctx, d, meta)
	}

	created, err := c.CreateRole(ctx, name, d.Get("description").(string))
	if err != nil {
		return diag.Errorf("failed to create role %q: %s", name, err)
	}
	d.SetId(created.Uuid)

	// Attach rules if specified
	if v, ok := d.GetOk("rules"); ok && v.(*schema.Set).Len() > 0 {
		roleUUID, err := uuid.FromString(created.Uuid)
		if err != nil {
			return diag.Errorf("invalid role UUID format: %s", err)
		}

		ruleIDs := v.(*schema.Set).List()
		ruleUUIDs := make([]string, 0, len(ruleIDs))

		for _, ruleID := range ruleIDs {
			ruleUUIDs = append(ruleUUIDs, ruleID.(string))
		}

		if len(ruleUUIDs) > 0 {
			if err := c.BulkAddRulesToRole(ctx, roleUUID, ruleUUIDs); err != nil {
				return diag.Errorf("failed to attach rules to role %q: %s", name, err)
			}
		}
	}

	return resourceRoleRead(ctx, d, meta)
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	id := d.Id()
	roleUUID, err := uuid.FromString(id)
	if err != nil {
		d.SetId("")
		return nil
	}

	res, err := c.GetRole(ctx, &roleUUID)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := d.Set("name", res.Name); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set name: %w", err))
	}

	// Get rules attached to this role
	rules, err := c.GetRoleRules(ctx, &roleUUID)
	if err != nil {
		return diag.Errorf("failed to load rulls %s", err.Error())
		// Don't fail the whole read operation if we can't get the rules
	} else {
		ruleIDs := make([]string, 0, len(rules))
		for _, rule := range rules {
			ruleIDs = append(ruleIDs, rule.Uuid)
		}
		if err := d.Set("rules", ruleIDs); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set rules: %w", err))
		}
	}

	return nil
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id := d.Id()
	roleUUID, err := uuid.FromString(id)
	if err != nil {
		return diag.Errorf("invalid role ID %q: %s", id, err)
	}

	if d.HasChange("name") {
		return diag.Errorf("name of role cannot be edited")
	}

	// Handle rule changes if the rules field has been changed
	if d.HasChange("rules") {
		old, new := d.GetChange("rules")
		oldSet := old.(*schema.Set)
		newSet := new.(*schema.Set)

		if newSet.Len() == 0 {
			// to prevent detach all rules by setting empty list as the "rules" is an optional field
			return nil
		}

		// Rules to add (in new but not in old)
		rulesToAdd := newSet.Difference(oldSet)
		if rulesToAdd.Len() > 0 {
			ruleIDs := rulesToAdd.List()
			ruleUUIDs := make([]string, 0, len(ruleIDs))

			for _, ruleID := range ruleIDs {
				ruleUUIDs = append(ruleUUIDs, ruleID.(string))
			}

			if len(ruleUUIDs) > 0 {
				if err := c.BulkAddRulesToRole(ctx, roleUUID, ruleUUIDs); err != nil {
					return diag.Errorf("failed to attach rules to role %q: %s", id, err)
				}
			}
		}

		// Rules to remove (in old but not in new)
		rulesToRemove := oldSet.Difference(newSet)
		if rulesToRemove.Len() > 0 {
			for _, ruleID := range rulesToRemove.List() {
				ruleUUID, err := uuid.FromString(ruleID.(string))
				if err != nil {
					return diag.Errorf("invalid rule UUID format for rule %s: %s", ruleID, err)
				}

				if err := c.UnbindRuleFromRole(ctx, &roleUUID, &ruleUUID); err != nil {
					return diag.Errorf("failed to detach rule %s from role %q: %s", ruleID, id, err)
				}
			}
		}
	}

	return resourceRoleRead(ctx, d, meta)
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	if err := c.DeleteRole(ctx, d.Id()); err != nil {
		if errors.Is(err, client.ErrNotFound) {
			d.SetId("")
			tflog.Warn(ctx, "Role already deleted or not found", map[string]interface{}{"role_id": d.Id()})
			return nil
		}
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
