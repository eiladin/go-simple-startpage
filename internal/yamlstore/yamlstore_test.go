package yamlstore

import (
	"os"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/models"
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
	net := models.Network{}
	err = f.GetNetwork(&net)
	suite.Error(err)
	site := models.Site{ID: 1}
	err = f.GetSite(&site)
	suite.Error(err)
}

func (suite *YamlStoreSuite) TestFunctions() {
	f, err := New("./testfile.yaml")
	suite.NoError(err)

	net := models.Network{
		Network: "test",
		Links: []models.Link{
			{Name: "test-link-1"},
			{Name: "test-link-2"},
		},
		Sites: []models.Site{
			{FriendlyName: "test-site-1"},
			{FriendlyName: "test-site-2"},
		},
	}
	suite.NoError(f.CreateNetwork(&net))
	// CreateNetwork assertions
	suite.Equal(uint(1), net.ID, "Network ID should be '1'")
	suite.Equal(uint(1), net.Sites[0].ID, "Site ID should be '1'")
	suite.Equal(uint(2), net.Sites[1].ID, "Site ID should be '2'")
	suite.Equal(uint(1), net.Links[0].ID, "Link ID should be '1'")
	suite.Equal(uint(2), net.Links[1].ID, "Link ID should be '2'")

	findNet := models.Network{ID: 1}
	suite.NoError(f.GetNetwork(&findNet))
	// GetNetwork assertions
	suite.Equal("test", findNet.Network, "Network should be 'test'")
	suite.Equal("test-site-1", findNet.Sites[0].FriendlyName, "Site FriendlyName should be 'test-site-1'")
	suite.Equal("test-link-1", findNet.Links[0].Name, "Link Name should be 'test-link-1'")

	findSite := models.Site{ID: 1}
	suite.NoError(f.GetSite(&findSite))
	// GetSite assertions
	suite.Equal("test-site-1", findSite.FriendlyName, "Site FriendlyName should be 'test-site-1'")

	missingSite := models.Site{ID: 3}
	err = f.GetSite(&missingSite)
	suite.EqualError(err, store.ErrNotFound.Error())
	os.Remove("./testfile.yaml")
}

func TestYamlStoreSuite(t *testing.T) {
	suite.Run(t, new(YamlStoreSuite))
}
