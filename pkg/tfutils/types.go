package tfutils

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ListToStringSlice(l types.List) []string {
	if l.IsNull() || l.IsUnknown() {
		return []string{}
	}

	var result []string
	for _, val := range l.Elements() {
		result = append(result, val.(types.String).ValueString())
	}

	return result
}

func StringSliceToList(slice []string) types.List {
	var result []attr.Value
	for _, val := range slice {
		result = append(result, types.StringValue(val))
	}

	return types.ListValueMust(types.StringType, result)
}
