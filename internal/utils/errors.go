package utils

import (
	"errors"

	iamTypes "git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"
)

const (
	ERROR_GROUP_CREATE         string = "Unable to Create Group"
	ERROR_GROUP_WORKSPACE_READ string = "Unable to Retreive Group Workspace"
	ERROR_GROUP_READ           string = "Unable to Retreive Group"
	ERROR_GROUP_DELETE         string = "Error Deleting Group"

	ERROR_ROLE_GROUP_BINDING_CREATE           string = "Error while creating role-group binding"
	ERROR_ROLE_GROUP_BINDING_READ             string = "Error getting role group binding"
	ERROR_ROLE_GROUP_BINDING_DELETE           string = "Error deleting role group binding"
	ERROR_ROLE_GROUP_BINDING_DELETE_TO_UPDATE string = ERROR_ROLE_GROUP_BINDING_DELETE + " (to create new one)"
	ERROR_ROLE_GROUP_BINDING_ITEMS_READ       string = "Error getting current role group binding items"

	ERROR_ROLE_SERVICE_USER_BINDING_CREATE           string = "Error creating role service user binding"
	ERROR_ROLE_SERVICE_USER_BINDING_READ             string = "Error getting role service user binding"
	ERROR_ROLE_SERVICE_USER_BINDING_DELETE           string = "Error deleting role service user binding"
	ERROR_ROLE_SERVICE_USER_BINDING_DELETE_TO_UPDATE string = ERROR_ROLE_SERVICE_USER_BINDING_DELETE + " (to create new one)"

	ERROR_ROLE_USER_BINDING_CREATE           string = "Error creating role user binding"
	ERROR_ROLE_USER_BINDING_READ             string = "Error getting role user binding"
	ERROR_ROLE_USER_BINDING_DELETE           string = "Error deleting role user binding"
	ERROR_ROLE_USER_BINDING_DELETE_TO_UPDATE string = ERROR_ROLE_USER_BINDING_DELETE + " (to create new one)"

	ERROR_ROLE_CREATE              string = "Unable to Create Role"
	ERROR_ROLE_READ                string = "Unable to Retreive Role"
	ERROR_ROLE_UPDATE              string = "Unable to Update Role"
	ERROR_ROLE_DELETE              string = "Error While Deleting Role"
	ERROR_ROLE_RULE_BINDING_CREATE string = "Unable to Bind Rule to Role"
	ERROR_ROLE_RULE_BINDING_READ   string = "Unable to Retreive binded Rules to Role"
	ERROR_ROLE_RULE_BINDING_DELETE string = "Unable to Unbind Rule to Role"

	ERROR_RULE_CREATE string = "Unable to Create Rule"
	ERROR_RULE_READ   string = "Unable to Retreive Rule"
	ERROR_RULE_UPDATE string = "Error Occured While Updating Rule"
	ERROR_RULE_DELETE string = "Error Occured While Updating Rule"

	ERROR_SERVICE_USER_GROUP_CREATE           string = "Error creating service_user group binding"
	ERROR_SERVICE_USER_GROUP_READ             string = "Error getting service_user group binding"
	ERROR_SERVICE_USER_GROUP_DELETE           string = "Error deleting service_user group binding"
	ERROR_SERVICE_USER_GROUP_DELETE_TO_UPDATE string = ERROR_SERVICE_USER_GROUP_DELETE + " (to create new one)"

	ERROR_USER_GROUP_CREATE           string = "Error creating user group binding"
	ERROR_USER_GROUP_READ             string = "Error getting user group binding"
	ERROR_USER_GROUP_DELETE           string = "Error deleting user group binding"
	ERROR_USER_GROUP_DELETE_TO_UPDATE string = ERROR_USER_GROUP_DELETE + " (to create new one)"
)

func GetIAMErrorMessage(err error) string {
	var iamClientError *iamTypes.RequestExecutionError
	errorMessage := err.Error()
	if errors.As(err, &iamClientError) {
		switch status := iamClientError.StatusCode; status {
		case 400:
			errorMessage = "Provider internal error! Please contact with Sotoon."
		case 401:
			errorMessage = "There is a problem about your authorization. Please check your access key."
		case 403:
			errorMessage = "Permision Denied."
		case 404:
			errorMessage = "Resource not found."
		case 409:
			errorMessage = "Conflict. The resource exists."
		case 500:
			errorMessage = "Sotoon internal error! Please contact with Sotoon."
		}
	}
	return errorMessage
}
