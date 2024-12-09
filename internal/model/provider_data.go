package model

import "github.com/conduktor/terraform-provider-conduktor/internal/client"

type ProviderData struct {
	Mode   *client.Mode
	Client *client.Client
}
