package utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moonlight8978/terraform-provider-rauthy/pkg/rauthy"
)

type Resource interface {
	SetClient(client *rauthy.Client)
}

func ConfigureProvider[TResource Resource](ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse, resource TResource) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*rauthy.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *rauthy.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	resource.SetClient(client)
}

// StringPtrToFramework converts a nullable API string pointer to a Terraform Framework types.String.
// Returns types.StringNull() if the pointer is nil, otherwise types.StringValue(*ptr).
func StringPtrToFramework(ptr *string) types.String {
	if ptr == nil {
		return types.StringNull()
	}
	return types.StringValue(*ptr)
}

// FrameworkToStringPtr converts a Terraform Framework types.String to a nullable API string pointer.
// Returns nil if the value is null or unknown, otherwise returns a pointer to the string value.
func FrameworkToStringPtr(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	s := v.ValueString()
	return &s
}
