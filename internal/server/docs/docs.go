// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Sami Khan",
            "url": "https://github.com/eiladin/go-simple-startpage"
        },
        "license": {
            "name": "MIT",
            "url": "https://github.com/eiladin/go-simple-startpage/blob/master/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/appconfig": {
            "get": {
                "description": "get application configuration",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "AppConfig"
                ],
                "summary": "Get AppConfig",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/config.Config"
                        }
                    }
                }
            }
        },
        "/api/healthz": {
            "get": {
                "description": "run healthcheck",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "HealthCheck"
                ],
                "summary": "Get Health",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Healthcheck"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/models.Healthcheck"
                        }
                    }
                }
            }
        },
        "/api/network": {
            "get": {
                "description": "get network",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Network"
                ],
                "summary": "Get Network",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Network"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/httperror.HTTPError"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/httperror.HTTPError"
                        }
                    }
                }
            },
            "post": {
                "description": "add or update network",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Network"
                ],
                "summary": "Add Network",
                "parameters": [
                    {
                        "description": "Add Network",
                        "name": "network",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Network"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.NetworkID"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httperror.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httperror.HTTPError"
                        }
                    }
                }
            }
        },
        "/api/status/{id}": {
            "get": {
                "description": "get status given a site id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Status"
                ],
                "summary": "Get Status",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Site ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Status"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httperror.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/httperror.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httperror.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "config.Config": {
            "type": "object",
            "properties": {
                "version": {
                    "type": "string"
                }
            }
        },
        "httperror.HTTPError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "object"
                }
            }
        },
        "models.Healthcheck": {
            "type": "object",
            "properties": {
                "errors": {
                    "type": "object",
                    "$ref": "#/definitions/models.HealthcheckErrors"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "models.HealthcheckErrors": {
            "type": "object",
            "properties": {
                "database": {
                    "type": "string"
                }
            }
        },
        "models.Link": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "uri": {
                    "type": "string"
                }
            }
        },
        "models.Network": {
            "type": "object",
            "properties": {
                "links": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Link"
                    }
                },
                "network": {
                    "type": "string"
                },
                "sites": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Site"
                    }
                }
            }
        },
        "models.NetworkID": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "models.Site": {
            "type": "object",
            "properties": {
                "friendlyName": {
                    "type": "string"
                },
                "icon": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "ip": {
                    "type": "string"
                },
                "isSupportedApp": {
                    "type": "boolean"
                },
                "isUp": {
                    "type": "boolean"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Tag"
                    }
                },
                "uri": {
                    "type": "string"
                }
            }
        },
        "models.Status": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "ip": {
                    "type": "string"
                },
                "isUp": {
                    "type": "boolean"
                }
            }
        },
        "models.Tag": {
            "type": "object",
            "properties": {
                "value": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Go Simple Startpage API",
	Description: "This is the API for the Go Simple Startpage App",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}