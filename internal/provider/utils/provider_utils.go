package utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
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
