package mapper

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type ConversionOrder string

const (
	FromTerraform ConversionOrder = "from terraform"
	IntoTerraform ConversionOrder = "into terraform"
)

// WrapDiagError wraps a diag.Diagnostics into an error with context.
func WrapDiagError(diag diag.Diagnostics, field string, order ConversionOrder) error {
	if field == "" {
		return fmt.Errorf("Mapping error, unable to convert %s: %v", order, diag)
	} else {
		return fmt.Errorf("Mapping error on %s, unable to convert %s: %v", field, order, diag)
	}
}

// WrapError wraps an error into an error with context.
func WrapError(err error, field string, order ConversionOrder) error {
	if field == "" {
		return fmt.Errorf("Mapping error, unable to convert %s: %v", order, err)
	} else {
		return fmt.Errorf("Mapping error on %s, unable to convert %s: %v", field, order, err)
	}
}
