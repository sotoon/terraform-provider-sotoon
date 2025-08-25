package provider

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

// parse strings to array of uuids and validate if each one is uuid
func convertStringsToUUIDArray(in []string) ([]string, error) {
	out := make([]string, 0, len(in))
	for _, s := range in {
		u, err := uuid.FromString(s)
		if err != nil {
			return nil, err
		}
		out = append(out, u.String())
	}
	return out, nil
}

// create sorted and unique array of uuids
func uniqueSorted(in []string) []string {
	if len(in) == 0 {
		return []string{}
	}

	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, v := range in {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

// creates a hash from array of uuids
func hashOfIDs(ids []string) string {
	h := sha256.Sum256([]byte(strings.Join(ids, ",")))
	return hex.EncodeToString(h[:])
}

// read array of strings from schema
func fromSchemaSetToStrings(s *schema.Set) []string {
	if s == nil {
		return nil
	}
	raw := s.List()
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		out = append(out, v.(string))
	}
	return out
}

// parse array of uuids to array of strings
func uuidsToStringSlice(in []uuid.UUID) []string {
	out := make([]string, len(in))
	for i, u := range in {
		out[i] = u.String()
	}
	return out
}

// check if there is two ids, check string for uuid and then returns two values with possible errors
func parseTwoPartID(id string) (uuid.UUID, uuid.UUID, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return uuid.UUID{}, uuid.UUID{}, fmt.Errorf("unexpected id: %s", id)
	}
	a, err := uuid.FromString(parts[0])
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, err
	}
	b, err := uuid.FromString(parts[1])
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, err
	}
	return a, b, nil
}
