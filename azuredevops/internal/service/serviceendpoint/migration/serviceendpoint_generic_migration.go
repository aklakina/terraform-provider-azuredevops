package migration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ServiceEndpointGenericV2Migrator creates a migrator resource for any type-specific service endpoint to generic_v2
func ServiceEndpointGenericV2Migrator(
	schemaMap map[string]*schema.Schema,
) *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource will be deprecated in favor of the new resource azuredevops_serviceendpoint_generic_v2",
		Schema:             schemaMap,
	}
}

// GenericV2UpgradeConfig contains configuration for upgrading to generic_v2
type GenericV2UpgradeConfig struct {
	EndpointType      string
	AuthScheme        string
	AuthParamsMapping map[string]string
	DataMapping       map[string]string
	ServerUrlKey      string
}

// ServiceEndpointGenericV2StateUpgrade creates a state upgrade function for any type-specific service endpoint to generic_v2
func ServiceEndpointGenericV2StateUpgrade(config GenericV2UpgradeConfig) schema.StateUpgradeFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// Set the type
		rawState["type"] = config.EndpointType

		// Set the authorization scheme
		rawState["authorization_scheme"] = config.AuthScheme

		// Map authorization parameters
		authParams := make(map[string]interface{})
		for oldKey, newKey := range config.AuthParamsMapping {
			if val, ok := rawState[oldKey]; ok {
				authParams[newKey] = val
				// Remove the old key from the state
				delete(rawState, oldKey)
			}
		}
		rawState["authorization_parameters"] = authParams

		// Map data parameters
		dataParams := make(map[string]interface{})
		for oldKey, newKey := range config.DataMapping {
			if val, ok := rawState[oldKey]; ok {
				dataParams[newKey] = val
				// Remove the old key from the state
				delete(rawState, oldKey)
			}
		}
		rawState["parameters"] = dataParams

		// Rename service_endpoint_name to name
		if val, ok := rawState["service_endpoint_name"]; ok {
			rawState["name"] = val
			delete(rawState, "service_endpoint_name")
		}

		// Handle server_url
		if config.ServerUrlKey != "" {
			if val, ok := rawState[config.ServerUrlKey]; ok {
				rawState["server_url"] = val
				delete(rawState, config.ServerUrlKey)
			}
		}

		return rawState, nil
	}
}
