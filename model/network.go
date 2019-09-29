package model

import (
	"github.com/eiladin/go-simple-startpage/db"
)

type Network struct {
	Network string `json:"n"`
	Links   []Link `json:"ls"`
	Sites   []Site `json:"ss"`
}

type Link struct {
	Name      string `json:"name"`
	Uri       string `json:"uri"`
	SortOrder int    `json:"sortOrder"`
}

type Site struct {
	FriendlyName   string `json:"friendlyName"`
	Uri            string `json:"uri"`
	Icon           string `json:"icon"`
	IsSupportedApp bool   `json:"isSupportedApp"`
	SortOrder      int    `json:"sortOrder"`
	Tags           []Tag  `json:"ts"`
}

type Tag struct {
	Value string `json:"value"`
}

func LoadNetwork() Network {
	n := db.ReadNetwork()
	return convertNetworkFromDbModel(n)
}

func SaveNetwork(n Network) uint {
	d := convertNetworkToDbModel(n)
	db.SaveNetwork(&d)
	return d.ID
}

func convertNetworkToDbModel(n Network) db.Network {
	d := db.Network{
		Network: n.Network,
	}
	for _, l := range n.Links {
		dl := db.Link{
			Name:      l.Name,
			Uri:       l.Uri,
			SortOrder: l.SortOrder,
		}
		d.Links = append(d.Links, dl)
	}
	for _, s := range n.Sites {
		ds := db.Site{
			FriendlyName:   s.FriendlyName,
			Uri:            s.Uri,
			Icon:           s.Icon,
			IsSupportedApp: s.IsSupportedApp,
			SortOrder:      s.SortOrder,
		}
		for _, t := range s.Tags {
			dt := db.Tag{
				Value: t.Value,
			}
			ds.Tags = append(ds.Tags, dt)
		}
		d.Sites = append(d.Sites, ds)
	}
	return d
}

func convertNetworkFromDbModel(d db.Network) Network {
	n := Network{
		Network: d.Network,
	}
	for _, dl := range d.Links {
		l := Link{
			Name:      dl.Name,
			Uri:       dl.Uri,
			SortOrder: dl.SortOrder,
		}
		n.Links = append(n.Links, l)
	}
	for _, ds := range d.Sites {
		s := Site{
			FriendlyName:   ds.FriendlyName,
			Uri:            ds.Uri,
			Icon:           ds.Icon,
			IsSupportedApp: ds.IsSupportedApp,
			SortOrder:      ds.SortOrder,
		}
		for _, dt := range ds.Tags {
			t := Tag{
				Value: dt.Value,
			}
			s.Tags = append(s.Tags, t)
		}
		n.Sites = append(n.Sites, s)
	}
	return n
}
