package model

import "github.com/conduktor/terraform-provider-conduktor/internal/client"

type ProviderData struct {
	ConsoleClient *client.ConsoleClient
	GatewayClient *client.GatewayClient
}
