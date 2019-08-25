package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/eiladin/go-simple-startpage/model"
	"github.com/gin-gonic/gin"
)

func GetNetworkHandler(c *gin.Context) {
	c.JSON(http.StatusOK, model.LoadNetwork())
}

func AddNetworkHandler(c *gin.Context) {
	networkItem, statusCode, err := convertHTTPBodyToNetwork(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	c.JSON(statusCode, gin.H{"id": model.SaveNetwork(networkItem)})
}

func GetStatusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, model.GetStatus())
}

func UpdateStatusHandler(c *gin.Context) {
	siteItem, statusCode, err := convertHTTPBodyToStatusSite(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	c.JSON(http.StatusOK, model.UpdateStatus(siteItem))
}

func convertHTTPBodyToStatusSite(httpBody io.ReadCloser) (model.StatusSite, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return model.StatusSite{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	return convertJSONBodyToStatusSite(body)
}

func convertJSONBodyToStatusSite(jsonBody []byte) (model.StatusSite, int, error) {
	var statusItem model.StatusSite
	err := json.Unmarshal(jsonBody, &statusItem)
	if err != nil {
		return model.StatusSite{}, http.StatusBadRequest, err
	}
	return statusItem, http.StatusOK, nil
}

func convertHTTPBodyToNetwork(httpBody io.ReadCloser) (model.Network, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return model.Network{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	return convertJSONBodyToNetwork(body)
}

func convertJSONBodyToNetwork(jsonBody []byte) (model.Network, int, error) {
	var networkItem model.Network
	err := json.Unmarshal(jsonBody, &networkItem)
	if err != nil {
		return model.Network{}, http.StatusBadRequest, err
	}
	return networkItem, http.StatusOK, nil
}
