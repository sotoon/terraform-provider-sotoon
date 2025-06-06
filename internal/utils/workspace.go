package utils

import (
	"fmt"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	uuid "github.com/satori/go.uuid"
)

func GetWokspaceNameFromUUID(workspaceUUID types.String, iamClient *client.Client, diag diag.Diagnostics) (string, error) {
	var wokrspaceName string
	if workspaceUUID.IsNull() {
		var ok bool
		wokrspaceName, ok = KNOWN_WORKSPACES_NAME[GLOBAL_WORKSPACE]
		if !ok {
			diag.AddError(
				"Cannot find gloabl workspace name",
				"It's a provider internal error. Please share this with maintainers.",
			)
			return "", fmt.Errorf("cannot find gloabl workspace name")
		}
	} else {
		workspaceUUIDString := workspaceUUID.ValueString()

		workspaceUUIDObject, _ := uuid.FromString(workspaceUUIDString)
		workspace, err := (*iamClient).GetWorkspace(&workspaceUUIDObject)
		if err != nil {
			diag.AddError(
				"Unable to Read Workspace",
				err.Error(),
			)
			return "", fmt.Errorf("unable to Read Workspace, error: %s", err.Error())
		}
		wokrspaceName = workspace.Name
	}
	return wokrspaceName, nil
}
