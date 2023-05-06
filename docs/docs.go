// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Tsvetan Dimitrov",
            "email": "tsvetan.dimitrov23@gmail.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://www.mit.edu/~amini/LICENSE.md"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/block/current": {
            "get": {
                "description": "Get current Ethereum block.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "blocks"
                ],
                "summary": "Get current Ethereum block.",
                "responses": {}
            }
        },
        "/api/v1/{address}/transactions": {
            "get": {
                "description": "Get all transactions for a fixed block range given an address.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "blocks"
                ],
                "summary": "Get all transactions for a fixed block range given an address.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Address",
                        "name": "address",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "0.0.0.0:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Ethereum Block Scanner API",
	Description:      "API for exploring Ethereum blocks.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}