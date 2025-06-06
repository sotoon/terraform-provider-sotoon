package types

import (
	uuid "github.com/satori/go.uuid"
)

type RoleBindingReq struct {
	Items map[string]string `json:"items,omitempty"`
}
type RoleBindingRes struct {
	Items []map[string]string `json:"items,omitempty"`
}
type RuleReq struct {
	Name    string   `json:"name" validate:"required"`
	Actions []string `json:"actions" validate:"required,gte=1"`
	Object  string   `json:"object" validate:"required"`
	Deny    bool     `json:"deny"`
}
type RoleReq struct {
	Name      string `json:"name" validate:"required"`
	Workspace string `json:"workspace" validate:"required"`
}
type RoleRes struct {
	UUID      *uuid.UUID `json:"uuid" faker:"uuidObject"`
	Name      string     `json:"name" validate:"required"`
	Workspace *Workspace `json:"workspace"`
}
type Role struct {
	UUID      *uuid.UUID `json:"uuid" faker:"uuidObject"`
	Workspace *uuid.UUID `json:"workspace" faker:"uuidObject"`
	Name      string     `json:"name"`
}

type RoleBinding struct {
	RoleName  string      `json:"name"`
	UserUUID  *uuid.UUID  `json:"user_uuid" faker:"uuidObject"`
	Workspace *Workspace  `json:"workspace"`
	Items     interface{} `json:"items,omitempty"`
	// todo map[string]string fails on empty map like []
}
type Rule struct {
	UUID          *uuid.UUID `json:"uuid" faker:"uuidObject"`
	WorkspaceUUID *uuid.UUID `json:"workspace" faker:"uuidObject"`
	Name          string     `json:"name"`
	Actions       []string   `json:"actions" validate:"required,gte=1"`
	Object        string     `json:"object" validate:"required"`
	Deny          bool       `json:"deny"`
}
