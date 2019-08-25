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
