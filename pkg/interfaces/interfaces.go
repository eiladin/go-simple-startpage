package interfaces

import "github.com/eiladin/go-simple-startpage/pkg/model"

// NetworkService interface
type NetworkService interface {
	CreateNetwork(net *model.Network)
	FindNetwork(net *model.Network)
}

// SiteService interface
type SiteService interface {
	FindSite(site *model.Site)
}
