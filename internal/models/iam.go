package models

import (
	"fmt"

	iamtypes "git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
	"git.cafebazaar.ir/infrastructure/integration/sib/terraform-provider/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
	uuid "github.com/satori/go.uuid"
)

type Rule struct {
	UUID          types.String   `tfsdk:"id"`
	Name          types.String   `tfsdk:"name"`
	WorkspaceUUID types.String   `tfsdk:"workspace_id"`
	Actions       []types.String `tfsdk:"actions"`
	Path          types.String   `tfsdk:"path"`
	Service       types.String   `tfsdk:"service"`
	IsDenial      types.Bool     `tfsdk:"is_denial"`
}

func (r *Rule) GetParams() (*uuid.UUID, string, *uuid.UUID, []string, string, string, bool, error) {
	ruleUUID, err := uuid.FromString(r.UUID.ValueString())
	if err != nil {
		return nil, "", nil, nil, "", "", false, fmt.Errorf("invalid rule id: %w", err)
	}
	WorkspaceUUID, err := uuid.FromString(r.WorkspaceUUID.ValueString())
	if err != nil {
		return nil, "", nil, nil, "", "", false, fmt.Errorf("invalid rule workspace id: %w", err)
	}
	return &ruleUUID,
		r.Name.ValueString(),
		&WorkspaceUUID,
		utils.GetStringListFromStringValueList(r.Actions),
		r.Path.ValueString(),
		r.Service.ValueString(),
		r.IsDenial.ValueBool(),
		nil
}

func NewRuleModelFromRuleObject(rule *iamtypes.Rule) *Rule {
	workspaceUUID, service, path := utils.ParseRRI(rule.Object)
	return &Rule{
		UUID:          types.StringValue(rule.UUID.String()),
		WorkspaceUUID: types.StringValue(workspaceUUID),
		Name:          types.StringValue(rule.Name),
		IsDenial:      types.BoolValue(rule.Deny),
		Service:       types.StringValue(service),
		Path:          types.StringValue(path),
		Actions:       utils.GetStringValueListFromStringList(rule.Actions),
	}
}

func GetRuleModelListFromRuleObjectList(rules []*iamtypes.Rule) []Rule {
	result := []Rule{}
	for _, rule := range rules {
		result = append(result, *NewRuleModelFromRuleObject(rule))
	}
	return result
}

type Role struct {
	UUID          types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	WorkspaceUUID types.String `tfsdk:"workspace_id"`
	Rules         []Rule       `tfsdk:"rules"`
}

type RoleWithRulesID struct {
	UUID          types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	WorkspaceUUID types.String `tfsdk:"workspace_id"`
	Rules         []UUIDField  `tfsdk:"rules"`
}

func (s *RoleWithRulesID) Equal(object *RoleWithRulesID) bool {
	return s.UUID.Equal(object.UUID) &&
		s.WorkspaceUUID.Equal(object.WorkspaceUUID) &&
		s.Name.Equal(object.Name)
}

func (s *RoleWithRulesID) ToRole() *Role {
	result := &Role{}
	result.UUID = s.UUID
	result.Name = s.Name
	result.WorkspaceUUID = s.WorkspaceUUID
	result.Rules = []Rule{}
	for _, rule := range s.Rules {
		result.Rules = append(result.Rules, Rule{
			UUID: rule.UUID,
		})
	}
	return result
}

func (s *Role) ToRoleWithRulesID() *RoleWithRulesID {
	result := &RoleWithRulesID{}
	result.UUID = s.UUID
	result.Name = s.Name
	result.WorkspaceUUID = s.WorkspaceUUID
	result.Rules = []UUIDField{}
	for _, rule := range s.Rules {
		result.Rules = append(result.Rules, UUIDField{
			UUID: rule.UUID,
		})
	}
	return result
}

func NewRoleModelFromRoleResAndRulesObject(role *iamtypes.RoleRes, roleRules []*iamtypes.Rule) *Role {
	rules := GetRuleModelListFromRuleObjectList(roleRules)
	return &Role{
		UUID:          types.StringValue((*role.UUID).String()),
		WorkspaceUUID: types.StringValue((*role.Workspace.UUID).String()),
		Name:          types.StringValue(role.Name),
		Rules:         rules,
	}
}

type User struct {
	UUID          types.String `tfsdk:"id"`
	WorkspaceUUID types.String `tfsdk:"workspace_id"`
	Name          types.String `tfsdk:"name"`
	Email         types.String `tfsdk:"email"`
}

func NewUserModelFromUserObject(user *iamtypes.User, workspaceUUID string) *User {
	return &User{
		WorkspaceUUID: types.StringValue(workspaceUUID),
		UUID:          types.StringValue((*user.UUID).String()),
		Name:          types.StringValue(user.Name),
		Email:         types.StringValue(user.Email),
	}
}

type Users struct {
	UUID          types.String `tfsdk:"id"`
	WorkspaceUUID types.String `tfsdk:"workspace_id"`
	GroupUUID     types.String `tfsdk:"group_id"`
	Users         []User       `tfsdk:"users"`
}

type ServiceUser struct {
	UUID          types.String `tfsdk:"id"`
	WorkspaceUUID types.String `tfsdk:"workspace_id"`
	Name          types.String `tfsdk:"name"`
}

func NewServiceUserModelFromServiceUserObject(user *iamtypes.ServiceUser, workspaceUUID string) *ServiceUser {
	return &ServiceUser{
		WorkspaceUUID: types.StringValue(workspaceUUID),
		UUID:          types.StringValue((*user.UUID).String()),
		Name:          types.StringValue(user.Name),
	}
}

type ServiceUsers struct {
	UUID          types.String  `tfsdk:"id"`
	WorkspaceUUID types.String  `tfsdk:"workspace_id"`
	GroupUUID     types.String  `tfsdk:"group_id"`
	ServiceUsers  []ServiceUser `tfsdk:"service_users"`
}

type RoleUser struct {
	UUID        types.String `tfsdk:"id"`
	UserID      types.String `tfsdk:"user_id"`
	RoleID      types.String `tfsdk:"role_id"`
	WorkspaceID types.String `tfsdk:"workspace_id"`
	Items       types.Map    `tfsdk:"items"`
}

type RoleServiceUser struct {
	UUID          types.String `tfsdk:"id"`
	ServiceUserID types.String `tfsdk:"service_user_id"`
	RoleID        types.String `tfsdk:"role_id"`
	WorkspaceID   types.String `tfsdk:"workspace_id"`
	Items         types.Map    `tfsdk:"items"`
}

type Group struct {
	UUID          types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	WorkspaceUUID types.String `tfsdk:"workspace_id"`
}

func NewGroupModelFromGroupObject(group *iamtypes.Group) *Group {
	return &Group{
		UUID:          types.StringValue(group.UUID.String()),
		WorkspaceUUID: types.StringValue(group.Workspace.UUID.String()),
		Name:          types.StringValue(group.Name),
	}
}

type UserGroup struct {
	UUID        types.String `tfsdk:"id"`
	WorkspaceID types.String `tfsdk:"workspace_id"`
	GroupID     types.String `tfsdk:"group_id"`
	UserID      types.String `tfsdk:"user_id"`
}

type RoleGroup struct {
	UUID        types.String `tfsdk:"id"`
	WorkspaceID types.String `tfsdk:"workspace_id"`
	GroupID     types.String `tfsdk:"group_id"`
	RoleID      types.String `tfsdk:"role_id"`
	Items       types.Map    `tfsdk:"items"`
}
