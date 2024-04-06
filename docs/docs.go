// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "email": "svenrisse0@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/books": {
            "post": {
                "description": "create book with fields",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "books"
                ],
                "summary": "Create a Book",
                "parameters": [
                    {
                        "description": "Add book",
                        "name": "book",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_svenrisse_bookshelf_internal_data.Book"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/github_com_svenrisse_bookshelf_internal_data.Book"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "422": {
                        "description": "Unprocessable Entity"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_svenrisse_bookshelf_internal_data.Book": {
            "description": "Book information",
            "type": "object",
            "properties": {
                "author": {
                    "type": "string",
                    "example": "J.R.R. Tolkien"
                },
                "genres": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "Fantasy",
                        "Epic",
                        "Children's literature"
                    ]
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "pages": {
                    "type": "integer",
                    "example": 320
                },
                "title": {
                    "type": "string",
                    "example": "The Hobbit"
                },
                "version": {
                    "type": "integer",
                    "example": 1
                },
                "year": {
                    "type": "integer",
                    "example": 1937
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "bookshelf.svenrisse.com",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "Bookshelf API",
	Description:      "This is the API for the Bookshelf application.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
