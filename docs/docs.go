// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/products": {
            "get": {
                "description": "get all products",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Get products with paging",
                "parameters": [
                    {
                        "type": "string",
                        "default": "1",
                        "description": "Number of page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "10",
                        "description": "Books count in a page",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.ProductsPaginatedResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controllers.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.InternalErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "create product",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Create a new product",
                "parameters": [
                    {
                        "description": "Request's body",
                        "name": "params",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.ProductCreateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controllers.InvalidRequestResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controllers.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.InternalErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/products/{id}": {
            "get": {
                "description": "get product by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Get product",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Product Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.ProductData"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controllers.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.InternalErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "update a product by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Update product",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Product Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Request's body",
                        "name": "params",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.ProductUpdateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controllers.InvalidRequestResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controllers.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.InternalErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "delete product by id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Delete product",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Product Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controllers.InvalidRequestResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controllers.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.InternalErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/products/{id}/sale": {
            "put": {
                "description": "product sale",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Product sale",
                "parameters": [
                    {
                        "description": "Request's body",
                        "name": "params",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.ProductSale"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controllers.InvalidRequestResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controllers.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.InternalErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.InternalErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "controllers.InvalidRequestResponse": {
            "type": "object",
            "properties": {
                "errors": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "controllers.ProductCreateReq": {
            "type": "object",
            "required": [
                "description",
                "name",
                "price",
                "stock"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "stock": {
                    "type": "integer"
                }
            }
        },
        "controllers.ProductData": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "stock": {
                    "type": "integer"
                }
            }
        },
        "controllers.ProductSale": {
            "type": "object",
            "required": [
                "count",
                "id"
            ],
            "properties": {
                "count": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                }
            }
        },
        "controllers.ProductUpdateReq": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "stockLevel": {
                    "type": "integer"
                }
            }
        },
        "controllers.ProductsPaginatedResponse": {
            "type": "object",
            "properties": {
                "meta": {
                    "$ref": "#/definitions/controllers.RequestMeta"
                },
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controllers.ProductData"
                    }
                }
            }
        },
        "controllers.RequestMeta": {
            "type": "object",
            "properties": {
                "current_page": {
                    "type": "integer"
                },
                "limit": {
                    "type": "integer"
                },
                "total_products": {
                    "type": "integer"
                }
            }
        },
        "controllers.Response": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
