package provider

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

var roleBindingItemsMu sync.Mutex

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

func hashOfBinding(ids []string, items map[string]string) string {
	payload, _ := json.Marshal(struct {
		IDs   []string          `json:"ids"`
		Items map[string]string `json:"items"`
	}{
		IDs:   uniqueSorted(ids),
		Items: items,
	})
	h := sha256.Sum256(payload)
	return hex.EncodeToString(h[:])
}

func roleItemsFromResourceData(d *schema.ResourceData) map[string]string {
	items := map[string]string{}
	if raw, ok := d.GetOk("items"); ok {
		for key, value := range raw.(map[string]interface{}) {
			items[key] = fmt.Sprintf("%v", value)
		}
	}
	return items
}

func roleItemsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for key, value := range a {
		if b[key] != value {
			return false
		}
	}
	return true
}

func roleBindingHasItems(bindingExists bool, remoteItems []map[string]string, wanted map[string]string) bool {
	if !bindingExists {
		return false
	}
	if len(wanted) == 0 && len(remoteItems) == 0 {
		return true
	}
	for _, items := range remoteItems {
		if roleItemsEqual(items, wanted) {
			return true
		}
	}
	return false
}

func removeRoleBindingItems(bindingExists bool, remoteItems []map[string]string, wanted map[string]string) ([]map[string]string, bool) {
	if !bindingExists {
		return remoteItems, false
	}
	if len(wanted) == 0 && len(remoteItems) == 0 {
		return nil, true
	}

	remaining := make([]map[string]string, 0, len(remoteItems))
	found := false
	for _, items := range remoteItems {
		if roleItemsEqual(items, wanted) {
			found = true
			continue
		}
		remaining = append(remaining, items)
	}
	return remaining, found
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

func toSet(xs []string) map[string]struct{} {
	m := make(map[string]struct{}, len(xs))
	for _, x := range xs {
		m[x] = struct{}{}
	}
	return m
}

func setKeys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func intersect(a, b map[string]struct{}) map[string]struct{} {
	out := make(map[string]struct{})
	for k := range a {
		if _, ok := b[k]; ok {
			out[k] = struct{}{}
		}
	}
	return out
}

func diff(a, b map[string]struct{}) []string {
	out := []string{}
	for k := range a {
		if _, ok := b[k]; !ok {
			out = append(out, k)
		}
	}
	return uniqueSorted(out)
}
