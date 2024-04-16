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
            "name": "API Support"
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
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
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
                        "name": "uuid",
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
        "/operations": {
            "get": {
                "description": "Переход с 1 страницы на случайную не предусмотрен! Пагинация быстрая с поиском по индексу. В качестве курсора пустое значение для начала, потом скопировать ПОСЛЕДНЮЮ дату ПОЛЯ CreatedAt , например 2024-02-18T16:27:05.271813Z",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Operations"
                ],
                "summary": "Получение операций с пагинацией",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 2,
                        "description": "Размер страницы",
                        "name": "page_size",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Курсор для пагинации",
                        "name": "cursor",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/operations.Response"
                            }
                        }
                    }
                }
            }
        },
        "/settings/execution-time": {
            "post": {
                "description": "operation_type: minus, plus, mult, div. execution_time \u003e 0",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Settings"
                ],
                "summary": "Установка нового времени выполнения",
                "parameters": [
                    {
                        "description": "Установка времени выполнения",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/execution_time.Request"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/execution_time.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "execution_time.Request": {
            "type": "object",
            "required": [
                "execution_time",
                "operation_type"
            ],
            "properties": {
                "execution_time": {
                    "type": "integer",
                    "example": 1
                },
                "operation_type": {
                    "type": "string",
                    "example": "plus"
                }
            }
        },
        "execution_time.Response": {
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
        "operations.Response": {
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
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Swagger API",
	Description:      "This is a distributed calculation server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
