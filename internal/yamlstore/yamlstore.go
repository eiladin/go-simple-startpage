package yamlstore

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/eiladin/go-simple-startpage/pkg/network"
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

func (d *YamlStore) CreateNetwork(net *network.Network) error {
	b, _ := yaml.Marshal(net)
	_, err := os.OpenFile(d.filepath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(d.filepath, b, 0644)
	return err
}

func (d *YamlStore) GetNetwork(net *network.Network) error {
	yamlFile, err := ioutil.ReadFile(d.filepath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yamlFile, net)
}

func (d *YamlStore) GetSite(site *network.Site) error {
	net := network.Network{}
	err := d.GetNetwork(&net)
	if err != nil {
		return err
	}

	found := false
	for _, s := range net.Sites {
		if site.Name == s.Name {
			site.Name = s.Name
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
