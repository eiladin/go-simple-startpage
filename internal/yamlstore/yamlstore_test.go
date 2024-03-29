package yamlstore

import (
	"os"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/network"
	"github.com/eiladin/go-simple-startpage/pkg/store"
	"github.com/stretchr/testify/suite"
)

type YamlStoreSuite struct {
	suite.Suite
}

func (suite *YamlStoreSuite) TestPing() {
	f, err := New("")
	suite.NoError(err)
	suite.NoError(f.Ping())
}

func (suite *YamlStoreSuite) TestFileNotFound() {
	f, err := New("./not-found.yaml")
	suite.NoError(err)
	net := network.Network{}
	err = f.GetNetwork(&net)
	suite.Error(err)
	site := network.Site{ID: 1}
	err = f.GetSite(&site)
	suite.Error(err)
}

func (suite *YamlStoreSuite) TestCreateFail() {
	suite.NoError(os.Mkdir("tmp", 0000))
	defer os.RemoveAll("tmp")
	filename := "./tmp/test-file"
	f, err := New(filename)
	suite.NoError(err)
	net := network.Network{
		Network: "test",
	}
	suite.Error(f.CreateNetwork(&net))
}

func (suite *YamlStoreSuite) TestFunctions() {
	f, err := New("./testfile.yaml")
	suite.NoError(err)

	net := network.Network{
		Network: "test",
		Links: []network.Link{
			{Name: "test-link-1"},
			{Name: "test-link-2"},
		},
		Sites: []network.Site{
			{
				Name: "test-site-1",
				Tags: []string{"tag-1"},
			},
			{
				Name: "test-site-2",
				Tags: []string{"tag-2"},
			},
		},
	}
	suite.NoError(f.CreateNetwork(&net))
	// CreateNetwork assertions
	suite.Equal("test", net.Network, "Network should be 'test'")
	suite.Equal("test-site-1", net.Sites[0].Name, "Site Name should be 'test-site-1'")
	suite.Equal("test-site-2", net.Sites[1].Name, "Site Name should be 'test-site-2'")
	suite.Equal("test-link-1", net.Links[0].Name, "Link Name should be 'test-link-1'")
	suite.Equal("test-link-2", net.Links[1].Name, "Link Name should be 'test-link-2'")
	suite.Equal("tag-1", net.Sites[0].Tags[0], "Tag Value should be 'tag-1'")
	suite.Equal("tag-2", net.Sites[1].Tags[0], "Tag Value should be 'tag-2'")

	findNet := network.Network{ID: 1}
	suite.NoError(f.GetNetwork(&findNet))
	// GetNetwork assertions
	suite.Equal("test", findNet.Network, "Network should be 'test'")
	suite.Equal("test-site-1", findNet.Sites[0].Name, "Site Name should be 'test-site-1'")
	suite.Equal("test-link-1", findNet.Links[0].Name, "Link Name should be 'test-link-1'")

	findSite := network.Site{Name: "test-site-1"}
	suite.NoError(f.GetSite(&findSite))
	// GetSite assertions
	suite.Equal("test-site-1", findSite.Name, "Site Name should be 'test-site-1'")

	missingSite := network.Site{ID: 3}
	err = f.GetSite(&missingSite)
	suite.EqualError(err, store.ErrNotFound.Error())
	os.Remove("./testfile.yaml")
}

func TestYamlStoreSuite(t *testing.T) {
	suite.Run(t, new(YamlStoreSuite))
}
