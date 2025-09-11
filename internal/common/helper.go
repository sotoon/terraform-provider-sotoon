package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

// create sorted and unique array of uuids
func UniqueSorted(in []string) []string {
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
func HashOfIDs(ids []string) string {
	h := sha256.Sum256([]byte(strings.Join(ids, ",")))
	return hex.EncodeToString(h[:])
}

// read array of strings from schema
func FromSchemaSetToStrings(s *schema.Set) []string {
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

// check if there is two ids, check string for uuid and then returns two values with possible errors
func ParseTwoPartID(id string) (uuid.UUID, uuid.UUID, error) {
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

func ToSet(xs []string) map[string]struct{} {
	m := make(map[string]struct{}, len(xs))
	for _, x := range xs {
		m[x] = struct{}{}
	}
	return m
}

func SetKeys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func Intersect(a, b map[string]struct{}) map[string]struct{} {
	out := make(map[string]struct{})
	for k := range a {
		if _, ok := b[k]; ok {
			out[k] = struct{}{}
		}
	}
	return out
}

func Diff(a, b map[string]struct{}) []string {
	out := []string{}
	for k := range a {
		if _, ok := b[k]; !ok {
			out = append(out, k)
		}
	}
	return UniqueSorted(out)
}
