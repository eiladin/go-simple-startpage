package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/eiladin/go-simple-startpage/config"
	"github.com/eiladin/go-simple-startpage/status"
	"github.com/gin-gonic/gin"
)

func GetConfigHandler(c *gin.Context) {
	c.JSON(http.StatusOK, config.Get())
}

func AddConfigHandler(c *gin.Context) {
	configItem, statusCode, err := convertHTTPBodyToConfig(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	c.JSON(statusCode, gin.H{"id": config.Add(configItem)})
}

func GetStatusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, status.Get())
}

func UpdateStatusHandler(c *gin.Context) {
	siteItem, statusCode, err := convertHTTPBodyToStatusSite(c.Request.Body)
	if err != nil {
		c.JSON(statusCode, err)
		return
	}
	c.JSON(http.StatusOK, status.UpdateStatus(siteItem))
}

func convertHTTPBodyToStatusSite(httpBody io.ReadCloser) (status.Site, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return status.Site{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	return convertJSONBodyToStatusSite(body)
}

func convertJSONBodyToStatusSite(jsonBody []byte) (status.Site, int, error) {
	var statusItem status.Site
	err := json.Unmarshal(jsonBody, &statusItem)
	if err != nil {
		return status.Site{}, http.StatusBadRequest, err
	}
	return statusItem, http.StatusOK, nil
}

func convertHTTPBodyToConfig(httpBody io.ReadCloser) (config.Config, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return config.Config{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	return convertJSONBodyToConfig(body)
}

func convertJSONBodyToConfig(jsonBody []byte) (config.Config, int, error) {
	var configItem config.Config
	err := json.Unmarshal(jsonBody, &configItem)
	if err != nil {
		return config.Config{}, http.StatusBadRequest, err
	}
	return configItem, http.StatusOK, nil
}
