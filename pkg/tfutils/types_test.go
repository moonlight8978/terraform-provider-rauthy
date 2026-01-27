package tfutils_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moonlight8978/terraform-provider-rauthy/pkg/tfutils"
	"github.com/stretchr/testify/assert"
)

func TestListToStringSlice(t *testing.T) {
	t.Parallel()

	list, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("a"),
		types.StringValue("b"),
		types.StringValue("c"),
	})

	result := tfutils.ListToStringSlice(list)

	assert.Equal(t, []string{"a", "b", "c"}, result)
}

func TestListToStringSlice_Empty(t *testing.T) {
	t.Parallel()

	list, _ := types.ListValue(types.StringType, []attr.Value{})

	result := tfutils.ListToStringSlice(list)

	assert.Equal(t, []string(nil), result)
}

func TestStringSliceToList(t *testing.T) {
	t.Parallel()

	result := tfutils.StringSliceToList([]string{"a", "b", "c"})

	assert.Equal(t, types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("a"),
		types.StringValue("b"),
		types.StringValue("c"),
	}), result)
}

func TestStringSliceToList_Empty(t *testing.T) {
	t.Parallel()

	result := tfutils.StringSliceToList([]string{})

	assert.Equal(t, types.ListValueMust(types.StringType, []attr.Value(nil)), result)
}
