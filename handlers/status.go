package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/eiladin/go-simple-startpage/model"
	"github.com/gin-gonic/gin"
)

func GetStatusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, model.GetStatus())
}

func UpdateStatusHandler(c *gin.Context) {
	s, status, err := convertHTTPBodyToStatusSite(c.Request.Body)
	if err != nil {
		c.JSON(status, err)
		return
	}
	c.JSON(http.StatusOK, model.UpdateStatus(s))
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
	var s model.StatusSite
	err := json.Unmarshal(jsonBody, &s)
	if err != nil {
		return model.StatusSite{}, http.StatusBadRequest, err
	}
	return s, http.StatusOK, nil
}