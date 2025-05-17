package utils

import "strings"

func CreateRRI(workspaceUUID, serviceName, path string) string {
	return "rri:v1:cafebazaar.cloud:" + workspaceUUID + ":" + serviceName + ":" + path
}

func ParseRRI(object string) (workspaceUUID, serviceName, path string) {
	parts := strings.Split(object, ":")
	return parts[3], parts[4], parts[5]
}
