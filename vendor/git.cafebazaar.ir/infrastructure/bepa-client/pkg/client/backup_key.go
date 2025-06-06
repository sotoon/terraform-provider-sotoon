package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/routes"
	"git.cafebazaar.ir/infrastructure/bepa-client/pkg/types"

	uuid "github.com/satori/go.uuid"
)

func (c *bepaClient) DeleteDefaultWorkspaceBackupKey(backupKeyUUID *uuid.UUID) error {
	replaceDict := map[string]string{
		workspaceUUIDPlaceholder: c.defaultWorkspace,
		backupKeyUUIDPlaceholder: backupKeyUUID.String(),
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteBackupKeyDelete), replaceDict)

	err := c.Do(http.MethodDelete, apiURL, 0, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *bepaClient) GetOneDefaultBackupKey(BackupKeyUUID *uuid.UUID) (*types.BackupKey, error) {
	replaceDict := map[string]string{
		workspaceUUIDPlaceholder: c.defaultWorkspace,
		backupKeyUUIDPlaceholder: BackupKeyUUID.String(),
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteBackupKeyGetOne), replaceDict)

	backupKey := &types.BackupKey{}
	err := c.Do(http.MethodGet, apiURL, 0, nil, &backupKey)
	if err != nil {
		return nil, err
	}
	return backupKey, nil
}

func (c *bepaClient) GetAllDefaultBackupKeys() ([]*types.BackupKey, error) {
	replaceDict := map[string]string{
		workspaceUUIDPlaceholder: c.defaultWorkspace,
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteBackupKeyGetAll), replaceDict)

	backupKeys := []*types.BackupKey{}
	err := c.Do(http.MethodGet, apiURL, 0, nil, &backupKeys)
	if err != nil {
		return nil, err
	}
	return backupKeys, nil
}

func (c *bepaClient) CreateBackupKeyForDefaultWorkspace(title, keyType, key string) (*types.BackupKey, error) {
	backupKeyReq := &types.BackupKeyReq{
		Title: title,
		Key:   fmt.Sprintf("%s %s", keyType, key),
	}

	replaceDict := map[string]string{
		workspaceUUIDPlaceholder: c.defaultWorkspace,
	}
	apiURL := substringReplace(trimURLSlash(routes.RouteBackupKeyCreate), replaceDict)

	createdBackupKey := &types.BackupKey{}
	if err := c.Do(http.MethodPost, apiURL, 0, backupKeyReq, createdBackupKey); err != nil {
		return nil, err
	}
	return createdBackupKey, nil
}

func (c *bepaClient) CreateBackupKeyFromFileForDefaultUser(title, fileAdd string) (*types.BackupKey, error) {
	if fileAdd == "" {
		fileAdd = os.Getenv("HOME") + "/.ssh/id_rsa.pub"
	}
	key, err := ioutil.ReadFile(fileAdd) // #nosec
	if err != nil {
		return nil, err
	}
	return c.CreateBackupKeyForDefaultWorkspace(title, defaultSSHKeyType, string(key))
}
