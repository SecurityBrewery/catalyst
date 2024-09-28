package analysis

import (
	"context"

	"github.com/SecurityBrewery/catalyst-analysis/plugin"
)

var _ plugin.ResourceType = &Other{}

type Other struct{}

func (u *Other) Info() plugin.ResourceTypeInfo {
	return plugin.ResourceTypeInfo{
		ID:                 "other",
		Name:               "Other",
		Attributes:         []string{},
		EnrichmentPatterns: []string{},
	}
}

func (u *Other) Resource(_ context.Context, id string) (*plugin.Resource, error) {
	return u.toOtherResource(id), nil
}

func (u *Other) toOtherResource(id string) *plugin.Resource {
	return &plugin.Resource{
		Type:       u.Info().ID,
		ID:         id,
		Name:       id,
		Icon:       "Box",
		Attributes: []plugin.Attribute{},
	}
}
