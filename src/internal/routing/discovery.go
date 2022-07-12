package routing

import (
	"github.com/futurehomeno/cliffhanger/discovery"
)

// GetDiscoveryResource returns a service discovery configuration.
func GetDiscoveryResource() *discovery.Resource {
	return &discovery.Resource{
		ResourceName:           ResourceName,
		ResourceType:           discovery.ResourceTypeApp,
		ResourceFullName:       "Iqcontrols",
		Author:                 "support@futurehome.no",
		IsInstanceConfigurable: false,
		InstanceID:             "1",
		Version:                "1",
	}
}
