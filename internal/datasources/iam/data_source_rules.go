package iamdatasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"

	"github.com/sotoon/terraform-provider-sotoon/internal/client"
)

func DataSourceRules() *schema.Resource {
	return &schema.Resource{
		Description: "Fetches a list of IAM rules within a specific Sotoon workspace.",
		ReadContext: dataSourceRulesRead,
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the workspace to fetch rules from.",
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of rules found in the workspace.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for the rule.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the rule.",
						},
						"actions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of actions this rule grants or denies.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"object": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The object this rule applies to.",
						},
						"deny": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this rule denies access (true) or allows (false).",
						},
					},
				},
			},
		},
	}
}

func dataSourceRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	workspaceID := d.Get("workspace_id").(string)

	tflog.Debug(ctx, "Reading rules for workspace", map[string]interface{}{"workspace_id": workspaceID})

	if _, err := uuid.FromString(workspaceID); err != nil {
		return diag.Errorf("Invalid workspace_id format: not a valid UUID")
	}

	rules, err := c.GetWorkspaceRules(ctx)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return diag.Errorf(
				"Forbidden: cannot list rules in this workspace. Check the token's IAM roles. Original error: %s",
				err,
			)
		}
		return diag.FromErr(err)
	}

	ruleList := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		actions := make([]interface{}, len(rule.Actions))
		for i, a := range rule.Actions {
			actions[i] = a
		}

		ruleList = append(ruleList, map[string]interface{}{
			"id":      rule.UUID.String(),
			"name":    rule.Name,
			"actions": actions,
			"object":  rule.Object,
			"deny":    rule.Deny,
		})
	}

	if err := d.Set("rules", ruleList); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set rules list: %w", err))
	}

	d.SetId(workspaceID)
	return nil
}
