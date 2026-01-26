package acctest

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/moonlight8978/terraform-provider-rauthy/internal/provider"
)

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"rauthy": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func TestAccPreCheck(t *testing.T) {
	apiKey := os.Getenv("RAUTHY_API_KEY")
	if apiKey == "" {
		t.Fatalf("RAUTHY_API_KEY environment variable is not set")
	}
}
