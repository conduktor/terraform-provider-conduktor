# Examples

This directory contains examples that are mostly used for documentation, but can also be run/tested manually via the Terraform CLI.

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. This is useful for creating examples that can run and/or ar testable even if some parts are not relevant for the documentation.

* **provider/console_provider.tf** example file for the provider index page in console mode
* **provider/gateway_provider.tf** example file for the provider index page in gateway mode
* **provider/multi_provider.tf** example file for the provider index page in multi client configuration
* **data-sources/`full data source name`/data-source.tf** example file for the named data source page
* **resources/`full resource name`/resource.tf** example file for the named data source page

In addition to individual resources a working project structure is provided to see all resources working together.
