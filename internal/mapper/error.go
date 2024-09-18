package errors

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type ConversionOrder string

const (
	FromTerraform ConversionOrder = "from terraform"
	IntoTerraform ConversionOrder = "into terraform"
)

func WrapDiagError(diag diag.Diagnostics, field string, order ConversionOrder) error {
	if field == "" {
		return fmt.Errorf("Mapping error, unable to convert %s: %v", order, diag)
	} else {
		return fmt.Errorf("Mapping error on %s, unable to convert %s: %v", field, order, diag)
	}
}
