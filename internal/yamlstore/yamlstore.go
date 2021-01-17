package yamlstore

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/eiladin/go-simple-startpage/pkg/store"
	"gopkg.in/yaml.v3"
)

var _ store.Store = (*YamlStore)(nil)

type YamlStore struct {
	filepath string
}

func New(filepath string) (store.Store, error) {
	d := YamlStore{filepath: filepath}
	return &d, nil
}

func (d *YamlStore) CreateNetwork(net *models.Network) error {
	net.ID = 1
	for i := range net.Sites {
		net.Sites[i].ID = uint(i + 1)
	}
	for i := range net.Links {
		net.Links[i].ID = uint(i + 1)
	}
	b, err := yaml.Marshal(net)
	if err != nil {
		return err
	}
	_, err = os.OpenFile(d.filepath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(d.filepath, b, 0644)
	return err
}

func (d *YamlStore) GetNetwork(net *models.Network) error {
	yamlFile, err := ioutil.ReadFile(d.filepath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yamlFile, net)
}

func (d *YamlStore) GetSite(site *models.Site) error {
	net := models.Network{}
	d.GetNetwork(&net)

	found := false
	for _, s := range net.Sites {
		if site.ID == s.ID {
			site.ID = s.ID
			site.FriendlyName = s.FriendlyName
			site.CreatedAt = s.CreatedAt
			site.Icon = s.Icon
			site.IsSupportedApp = s.IsSupportedApp
			site.NetworkID = s.NetworkID
			site.Tags = s.Tags
			site.URI = s.URI
			found = true
			break
		}
	}

	if !found {
		return errors.New("record not found")
	}
	return nil
}

func (d *YamlStore) Ping() error {
	return nil
}
