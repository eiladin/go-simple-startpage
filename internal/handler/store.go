package handler

import "github.com/eiladin/go-simple-startpage/pkg/model"

// Store interface
type Store interface {
	CreateNetwork(net *model.Network)
	GetNetwork(net *model.Network)
	GetSite(site *model.Site)
}
