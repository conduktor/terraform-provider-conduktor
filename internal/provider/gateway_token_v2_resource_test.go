package provider

import (
	"testing"
	"time"

	"github.com/conduktor/terraform-provider-conduktor/internal/test"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGatewayTokenV2Resource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resourceRef := "conduktor_gateway_token_v2.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/token_v2/resource_create.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "vcluster", "vcluster_sa"),
					resource.TestCheckResourceAttr(resourceRef, "username", "user10"),
					resource.TestCheckResourceAttr(resourceRef, "lifetime_seconds", "3600"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/token_v2/resource_update.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "vcluster", "vcluster_sa"),
					resource.TestCheckResourceAttr(resourceRef, "username", "user10"),
					resource.TestCheckResourceAttr(resourceRef, "lifetime_seconds", "3000"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccGatewayTokenV2Minimal(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read from minimal example
			{
				Config: providerConfigGateway + test.TestAccTestdata(t, "gateway/token_v2/resource_minimal.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.minimal", "vcluster", "passthrough"),
					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.minimal", "username", "user_passthrough"),
					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.minimal", "lifetime_seconds", "3600"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccGatewayTokenV2ExampleResource(t *testing.T) {
	test.CheckEnterpriseEnabled(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Create and Read from simple example
			{
				Config: providerConfigGateway + test.TestAccExample(t, "resources", "conduktor_gateway_token_v2", "simple.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.simple", "vcluster", "passthrough"),
					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.simple", "username", "user_passthrough"),
					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.simple", "lifetime_seconds", "3600"),
				),
			},
			// Create and Read from complex example
			{
				Config: providerConfigGateway + test.TestAccExample(t, "resources", "conduktor_gateway_token_v2", "complex.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.complex", "vcluster", "vcluster_sa"),
					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.complex", "username", "user10"),
					resource.TestCheckResourceAttr("conduktor_gateway_token_v2.complex", "lifetime_seconds", "3600"),
				),
			},
		},
	})
}

func TestIsTokenExpired(t *testing.T) {
	// Helper function to create a token string
	createToken := func(expirationTime time.Time) string {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": expirationTime.Unix(),
		})
		tokenString, _ := token.SignedString([]byte("secret"))
		return tokenString
	}

	tests := []struct {
		name          string
		tokenString   string
		expected      bool
		expectedError bool
	}{
		{
			name:          "Valid token",
			tokenString:   createToken(time.Now().Add(time.Hour)),
			expected:      false,
			expectedError: false,
		},
		{
			name:          "Expired token",
			tokenString:   createToken(time.Now().Add(-time.Hour)),
			expected:      true,
			expectedError: false,
		},
		{
			name:          "Invalid token",
			tokenString:   "invalid.token.string",
			expected:      false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expired, err := isTokenExpired(tt.tokenString)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, expired)
		})
	}
}
