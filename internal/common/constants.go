package common

import (
	uuid "github.com/satori/go.uuid"
)

// GlobalWorkspaceUUID is the UUID for the global workspace that contains global roles and rules
var GlobalWorkspaceUUID = uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000")
