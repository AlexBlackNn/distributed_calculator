// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
            "name": "API Support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/expression": {
            "post": {
                "description": "Создает новое выражение на сервере",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Calculations"
                ],
                "summary": "Создание нового выражения",
                "parameters": [
                    {
                        "description": "Запрос на создание выражения",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/expression.Request"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/expression.Response"
                        }
                    }
                }
            }
        },
        "/expression/{uuid}": {
            "get": {
                "description": "Получает результат по указанному идентификатору из хранилища",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Calculations"
                ],
                "summary": "Получение результата по идентификатору",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Идентификатор результата",
                        "name": "uid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/result.Response"
                        }
                    }
                }
            }
        },
        "/monitoring/worker": {
            "get": {
                "description": "Получает количество воркеров доступных для выполнения задачи",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Monitoring"
                ],
                "summary": "Получение количества активных воркеров",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/worker.Response"
                        }
                    }
                }
            }
        },
        "/settings/plus-execution-time": {
            "post": {
                "description": "Создает новое выражение на сервере",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Settings"
                ],
                "summary": "Создание нового выражения",
                "parameters": [
                    {
                        "description": "Запрос на создание выражения",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/execution_time_plus.Request"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/execution_time_plus.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "execution_time_plus.Request": {
            "type": "object",
            "required": [
                "execution_time"
            ],
            "properties": {
                "execution_time": {
                    "type": "integer"
                }
            }
        },
        "execution_time_plus.Response": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "response": {
                    "$ref": "#/definitions/response.Response"
                }
            }
        },
        "expression.Request": {
            "type": "object",
            "required": [
                "expression"
            ],
            "properties": {
                "expression": {
                    "type": "string"
                }
            }
        },
        "expression.Response": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "response": {
                    "$ref": "#/definitions/response.Response"
                }
            }
        },
        "response.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "result": {
                    "type": "number"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "result.Response": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "response": {
                    "$ref": "#/definitions/response.Response"
                }
            }
        },
        "worker.Response": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "response": {
                    "$ref": "#/definitions/response.Response"
                }
            }
        }
    },
    "externalDocs": {
        "description": "OpenAPI"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Swagger Example API",
	Description:      "This is a sample server celler server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
