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

var defaultSSHKeyType = "ssh-rsa"

func (c *bepaClient) DeleteDefaultUserPublicKey(publicKeyUUID *uuid.UUID) error {
	replaceDict := map[string]string{
		userUUIDPlaceholder:      c.userUUID,
		publicKeyUUIDPlaceholder: publicKeyUUID.String(),
	}
	apiURL := substringReplace(trimURLSlash(routes.RoutePublicKeyDelete), replaceDict)

	err := c.Do(http.MethodDelete, apiURL, 0, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *bepaClient) GetOneDefaultUserPublicKey(publicKeyUUID *uuid.UUID) (*types.PublicKey, error) {
	replaceDict := map[string]string{
		userUUIDPlaceholder:      c.userUUID,
		publicKeyUUIDPlaceholder: publicKeyUUID.String(),
	}
	apiURL := substringReplace(trimURLSlash(routes.RoutePublicKeyGetOne), replaceDict)

	publicKey := &types.PublicKey{}
	err := c.Do(http.MethodGet, apiURL, 0, nil, &publicKey)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

func (c *bepaClient) GetAllDefaultUserPublicKeys() ([]*types.PublicKey, error) {
	replaceDict := map[string]string{
		userUUIDPlaceholder: c.userUUID,
	}
	apiURL := substringReplace(trimURLSlash(routes.RoutePublicKeyGetAll), replaceDict)

	publicKeys := []*types.PublicKey{}
	err := c.Do(http.MethodGet, apiURL, 0, nil, &publicKeys)
	if err != nil {
		return nil, err
	}
	return publicKeys, nil
}

func (c *bepaClient) CreatePublicKeyForDefaultUser(title, keyType, key string) (*types.PublicKey, error) {
	publicKeyReq := &types.PublicKeyReq{
		Title: title,
		Key:   fmt.Sprintf("%s %s", keyType, key),
	}

	replaceDict := map[string]string{
		userUUIDPlaceholder: c.userUUID,
	}
	apiURL := substringReplace(trimURLSlash(routes.RoutePublicKeyCreate), replaceDict)

	createdPublicKey := &types.PublicKey{}
	if err := c.Do(http.MethodPost, apiURL, 0, publicKeyReq, createdPublicKey); err != nil {
		return nil, err
	}
	return createdPublicKey, nil
}

func (c *bepaClient) CreatePublicKeyFromFileForDefaultUser(title, fileAdd string) (*types.PublicKey, error) {
	if fileAdd == "" {
		fileAdd = os.Getenv("HOME") + "/.ssh/id_rsa.pub"
	}
	key, err := ioutil.ReadFile(fileAdd) // #nosec
	if err != nil {
		return nil, err
	}
	return c.CreatePublicKeyForDefaultUser(title, defaultSSHKeyType, string(key))
}

func (c *bepaClient) VerifyPublicKey(keyType string, key string, workspaceUUID string, username string, hostname string) (bool, error) {
	publicKeyVerifyReq := &types.PublicKeyVerifyReq{
		KeyType:        keyType,
		Key:            key,
		Workspace_uuid: workspaceUUID,
		Email:          username,
		Hostname:       hostname,
	}

	replaceDict := map[string]string{
		userUUIDPlaceholder: c.userUUID,
	}
	apiURL := substringReplace(trimURLSlash(routes.RoutePublicKeyVerify), replaceDict)

	if err := c.Do(http.MethodPost, apiURL, 0, publicKeyVerifyReq, nil); err != nil {
		return false, err
	}
	return true, nil
}
