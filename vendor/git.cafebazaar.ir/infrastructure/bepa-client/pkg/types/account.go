package types

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	UserTypeUser        = "user"
	UserTypeServiceUser = "service-user"
)

type UserRes struct {
	UUID            string              `json:"uuid"`
	Name            string              `json:"name"`
	Email           string              `json:"email"`
	UserType        string              `json:"user_type"`
	IsSuspended     bool                `json:"is_suspended"`
	CreatedAt       time.Time           `json:"created_at,omitempty"`
	UpdatedAt       time.Time           `json:"updated_at,omitempty"`
	InvitationToken string              `json:"invitation_token,omitempty"`
	Items           []map[string]string `json:"items,omitempty"`
}
type UserTokenReq struct {
	Secret   string `json:"secret" validate:"required"`
	UserType string `json:"user_type"`
}
type UserUpdateReq struct {
	Name     string `json:"name"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=8"`
}

type UserReq struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserAcceptInvitationReq struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserCanReq struct {
	Path string `json:"path" validate:"required"`
}

type UserSecretRes struct {
	Secret string `json:"secret"`
}

type UserTokenByCredsReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type InviteUserReq struct {
	Email string `json:"email" validate:"required,email"`
}

type WorkspaceReq struct {
	Name string `json:"name" validate:"required,rfc1123_label"`
}

type PublicKeyReq struct {
	Title string `json:"title"`
	Key   string `json:"key" validate:"required"`
}

type VerifRes struct {
	Message string `json:"message"`
}

type PublicKeyVerifyReq struct {
	Key            string
	KeyType        string
	Workspace_uuid string `json:"workspace_uuid"`
	Hostname       string `json:"hostname"`
	Email          string `json:"email"`
}

type Workspace struct {
	UUID         *uuid.UUID `json:"uuid" faker:"uuidObject"`
	Name         string     `json:"name"`
	Namespace    string     `json:"namespace"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Organization *uuid.UUID `json:"organization" faker:"uuidObject"`
}

type WorkspaceWithOrganization struct {
	UUID         *uuid.UUID    `json:"uuid" faker:"uuidObject"`
	Name         string        `json:"name"`
	Namespace    string        `json:"namespace"`
	Organization *Organization `json:"organization"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type Organization struct {
	UUID           *uuid.UUID `json:"uuid" faker:"uuidObject"`
	Name           string     `json:"name_en"`
	EnterpriseName string     `json:"enterprise_name"`
	EconomicCode   string     `json:"economic_code"`
	NationalId     string     `json:"national_id"`
}

type AuthnChallengeRequiredResponse struct {
	ChallengeToken string `json:"challenge_token"`
	ChallengeType  string `json:"challenge_type"`
}

type AuthnChallengeRequest struct {
	ChallengeToken  string `json:"challenge_token"`
	ChallengeAnswer string `json:"challenge_answer"`
}

func (r *AuthnChallengeRequiredResponse) Error() string {
	return fmt.Sprintf("challenge of type '%s' required", r.ChallengeType)
}

type UserToken struct {
	UUID         string     `json:"uuid"`
	User         string     `json:"user"`
	Secret       string     `json:"secret"`
	Active       bool       `json:"active"`
	LastAccessAt *time.Time `json:"last_access_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

type User struct {
	UUID            *uuid.UUID `json:"uuid" faker:"uuidObject"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	InvitationToken string     `json:"invitation_token,omitempty"`
}
type Group struct {
	UUID      *uuid.UUID `json:"uuid" faker:"uuidObject"`
	Name      string     `json:"name"`
	Workspace Workspace  `json:"workspace"`
}
type GroupReq struct {
	Name      string `json:"name"`
	Workspace string `json:"workspace"`
}
type GroupRes struct {
	UUID          *uuid.UUID `json:"uuid" faker:"uuidObject"`
	Name          string     `json:"name"`
	WorkspaceUUID string     `json:"workspace"`
	Descriotion   string     `json:"description"`
}
type GroupUserRes struct {
	Group string `json:"group"`
	User  string `json:"user"`
}
type ServiceUser struct {
	UUID      *uuid.UUID `json:"uuid" faker:"uuidObject"`
	Name      string     `json:"name"`
	Workspace string     `json:"workspace"`
}
type ServiceUserReq struct {
	Name      string `json:"name"`
	Workspace string `json:"workspace"`
}
type ServiceUserToken struct {
	UUID        *uuid.UUID `json:"uuid" faker:"uuidObject"`
	ServiceUser string     `json:"service_user"`
	Secret      string     `json:"secret"`
}

type InvitationInfo struct {
	Token string `json:"invitation_token"`
}

type UserSecret struct {
	Secret string `json:"secret"`
}
type Service struct {
	Name    string   `json:"name"`
	Actions []string `json:"actions"`
}

type PublicKey struct {
	UUID  string `json:"uuid"`
	Title string `json:"title"`
	Key   string `json:"key"`
	User  string `json:"user"`
}

type KiseSecret struct {
	UUID   string `json:"uuid"`
	Title  string `json:"title"`
	Secret string `json:"secret"`
	User   string `json:"user"`
}
type BackupKey struct {
	UUID      string `json:"uuid"`
	Title     string `json:"title"`
	Key       string `json:"key"`
	Type      string `json:"type"`
	Workspace string `json:"workspace"`
}
type BackupKeyReq struct {
	Title string `json:"title"`
	Key   string `json:"key" validate:"required"`
}

type WebhookWorkspaceOrganization struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

type WebhookWorkspace struct {
	UUID         uuid.UUID    `json:"uuid"`
	Name         string       `json:"name"`
	IsSuspended  bool         `json:"is_suspended"`
	Organization Organization `json:"organization"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type WebhookUser struct {
	UUID                  uuid.UUID `json:"uuid"`
	Name                  string    `json:"name"`
	Email                 string    `json:"email"`
	IsEmailVerified       bool      `json:"email_verified"`
	PhoneNumber           string    `json:"phone_number"`
	IsPhoneNumberVerified bool      `json:"phone_number_verified"`
	FirstName             string    `json:"first_name"`
	LastName              string    `json:"last_name"`
	Birthday              string    `json:"birthday"`
	IsSuspended           bool      `json:"is_suspended"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type WebhookUserWorkspaceRelation struct {
	UserUUID      uuid.UUID `json:"user_uuid"`
	WorkspaceUUID uuid.UUID `json:"workspace_uuid"`
	IsSuspended   bool      `json:"is_suspended"`
}
