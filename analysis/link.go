package analysis

import (
	"context"

	"github.com/SecurityBrewery/catalyst-analysis/plugin"
)

var _ plugin.ResourceType = &Link{}

type Link struct{}

func (u *Link) Info() plugin.ResourceTypeInfo {
	return plugin.ResourceTypeInfo{
		ID:                 "link",
		Name:               "Link",
		Attributes:         []string{},
		EnrichmentPatterns: []string{},
	}
}

func (u *Link) Resource(_ context.Context, id string) (*plugin.Resource, error) {
	return u.toLinkResource(id), nil
}

func (u *Link) toLinkResource(id string) *plugin.Resource {
	return &plugin.Resource{
		Type:       u.Info().ID,
		ID:         id,
		Name:       id,
		Icon:       "Link",
		Attributes: []plugin.Attribute{},
	}
}
