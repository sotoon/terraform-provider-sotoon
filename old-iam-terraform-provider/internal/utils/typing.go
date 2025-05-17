package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetStringListFromStringValueList(stringValueList []types.String) []string {
	result := []string{}
	for _, element := range stringValueList {
		value := element.ValueString()
		if value != "" {
			result = append(result, element.ValueString())
		}
	}
	return result
}

func GetStringValueListFromStringList(stringValueList []string) []types.String {
	result := []types.String{}
	for _, element := range stringValueList {
		if element != "" {
			result = append(result, types.StringValue(element))
		}
	}
	return result
}
