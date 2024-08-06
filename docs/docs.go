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
        "/telegram/{botplatform}/{omnichannel}/{tenantId}/{account}": {
            "post": {
                "description": "Incoming message for channel telegram",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "telegram"
                ],
                "summary": "Incoming telegram",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bot Platform",
                        "name": "botplatform",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Omni Channel",
                        "name": "omnichannel",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Tenant",
                        "name": "tenantId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Account",
                        "name": "account",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/webhook.IncomingDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/util.Responses"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "object"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/util.Responses"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "object"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/telegram/{botplatform}/{omnichannel}/{tenantId}/{account}/end": {
            "post": {
                "description": "End message for channel telegram",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "telegram"
                ],
                "summary": "End telegram",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bot Platform",
                        "name": "botplatform",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Omni Channel",
                        "name": "omnichannel",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Tenant",
                        "name": "tenantId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Account",
                        "name": "account",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/webhook.EndDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/util.Responses"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "object"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/util.Responses"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "object"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/telegram/{botplatform}/{omnichannel}/{tenantId}/{account}/handover": {
            "post": {
                "description": "Handover message for channel telegram",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "telegram"
                ],
                "summary": "Handover telegram",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bot Platform",
                        "name": "botplatform",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Omni Channel",
                        "name": "omnichannel",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Tenant",
                        "name": "tenantId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Account",
                        "name": "account",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/webhook.HandoverDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/util.Responses"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "object"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/util.Responses"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "object"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/tole/{queueName}": {
            "post": {
                "description": "Send a message to rabbit tole queue using tenant channel",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tole"
                ],
                "summary": "Send a message to tole queue",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Queue Name",
                        "name": "queueName",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/util.Responses"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "object"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/util.Responses"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "object"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "util.Responses": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                },
                "method": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "webhook.AttachmentPayload": {
            "type": "object",
            "required": [
                "url"
            ],
            "properties": {
                "url": {
                    "type": "string"
                }
            }
        },
        "webhook.Attachments": {
            "type": "object",
            "required": [
                "payload",
                "type"
            ],
            "properties": {
                "payload": {
                    "$ref": "#/definitions/webhook.AttachmentPayload"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "webhook.AttributeDTO": {
            "type": "object",
            "required": [
                "accountId",
                "botplatform",
                "channel_id",
                "channel_platform",
                "channel_sources",
                "cust_name",
                "date_timestamp",
                "middleware_endpoint",
                "omnichannel",
                "tenantId",
                "unique_id"
            ],
            "properties": {
                "accountId": {
                    "type": "string"
                },
                "botplatform": {
                    "enum": [
                        "botpress"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/webhook.BotPlatform"
                        }
                    ]
                },
                "channel_id": {
                    "enum": [
                        12,
                        3,
                        7,
                        5
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/webhook.ChannelID"
                        }
                    ]
                },
                "channel_platform": {
                    "enum": [
                        "socioconnect",
                        "maytapi",
                        "octopushchat",
                        "official"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/webhook.ChannelPlatform"
                        }
                    ]
                },
                "channel_sources": {
                    "enum": [
                        "whatsapp",
                        "fbmessenger",
                        "livechat",
                        "telegram"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/webhook.ChannelSources"
                        }
                    ]
                },
                "cust_message": {
                    "type": "string"
                },
                "cust_name": {
                    "type": "string"
                },
                "date_timestamp": {
                    "type": "string"
                },
                "middleware_endpoint": {
                    "type": "string"
                },
                "omnichannel": {
                    "enum": [
                        "onx",
                        "on5",
                        "on4"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/webhook.Omnichannel"
                        }
                    ]
                },
                "stream_id": {
                    "type": "string"
                },
                "tenantId": {
                    "type": "string"
                },
                "unique_id": {
                    "type": "string"
                }
            }
        },
        "webhook.BotPlatform": {
            "type": "string",
            "enum": [
                "botpress"
            ],
            "x-enum-varnames": [
                "BOTPRESS"
            ]
        },
        "webhook.ChannelID": {
            "type": "integer",
            "enum": [
                12,
                3,
                7,
                5
            ],
            "x-enum-varnames": [
                "WHATSAPP_ID",
                "LIVECHAT_ID",
                "FBMESSENGER_ID",
                "TELEGRAM_ID"
            ]
        },
        "webhook.ChannelPlatform": {
            "type": "string",
            "enum": [
                "socioconnect",
                "maytapi",
                "octopushchat",
                "official"
            ],
            "x-enum-varnames": [
                "SOCIOCONNECT",
                "MAYTAPI",
                "OCTOPUSHCHAT",
                "OFFICIAL"
            ]
        },
        "webhook.ChannelSources": {
            "type": "string",
            "enum": [
                "whatsapp",
                "fbmessenger",
                "livechat",
                "telegram"
            ],
            "x-enum-varnames": [
                "WHATSAPP",
                "FBMESSENGER",
                "LIVECHAT",
                "TELEGRAM"
            ]
        },
        "webhook.Data": {
            "type": "object",
            "required": [
                "entry",
                "object"
            ],
            "properties": {
                "entry": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/webhook.Entry"
                    }
                },
                "object": {
                    "type": "string"
                }
            }
        },
        "webhook.EndDTO": {
            "type": "object",
            "required": [
                "account_id",
                "message",
                "sid",
                "unique_id"
            ],
            "properties": {
                "account_id": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "sid": {
                    "type": "string"
                },
                "unique_id": {
                    "type": "string"
                }
            }
        },
        "webhook.Entry": {
            "type": "object",
            "required": [
                "id",
                "messaging",
                "time"
            ],
            "properties": {
                "hop_context": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/webhook.HopContext"
                    }
                },
                "id": {
                    "type": "string"
                },
                "messaging": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/webhook.Messaging"
                    }
                },
                "time": {
                    "type": "integer"
                }
            }
        },
        "webhook.HandoverDTO": {
            "type": "object",
            "required": [
                "account_id",
                "message",
                "sid",
                "unique_id"
            ],
            "properties": {
                "account_id": {
                    "type": "string"
                },
                "cust_message": {},
                "message": {
                    "type": "string"
                },
                "sid": {
                    "type": "string"
                },
                "unique_id": {
                    "type": "string"
                }
            }
        },
        "webhook.HopContext": {
            "type": "object",
            "required": [
                "app_id"
            ],
            "properties": {
                "app_id": {
                    "type": "integer"
                },
                "metadata": {
                    "type": "string"
                }
            }
        },
        "webhook.IncomingDTO": {
            "type": "object",
            "required": [
                "account",
                "account_name",
                "additional",
                "channel",
                "data",
                "tenant",
                "test"
            ],
            "properties": {
                "account": {
                    "type": "string"
                },
                "account_name": {
                    "type": "string"
                },
                "additional": {
                    "$ref": "#/definitions/webhook.AttributeDTO"
                },
                "channel": {
                    "type": "string"
                },
                "data": {
                    "$ref": "#/definitions/webhook.Data"
                },
                "tenant": {
                    "type": "string"
                },
                "test": {
                    "type": "string",
                    "enum": [
                        "12",
                        "3",
                        "7",
                        "5"
                    ]
                }
            }
        },
        "webhook.Message": {
            "type": "object",
            "required": [
                "mid"
            ],
            "properties": {
                "attachments": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/webhook.Attachments"
                    }
                },
                "mid": {
                    "type": "string"
                },
                "quick_reply": {
                    "$ref": "#/definitions/webhook.QuickReply"
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "webhook.MessagePostback": {
            "type": "object",
            "required": [
                "payload",
                "title"
            ],
            "properties": {
                "payload": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "webhook.Messaging": {
            "type": "object",
            "required": [
                "recipient",
                "sender",
                "timestamp"
            ],
            "properties": {
                "message": {
                    "$ref": "#/definitions/webhook.Message"
                },
                "postback": {
                    "$ref": "#/definitions/webhook.MessagePostback"
                },
                "recipient": {
                    "$ref": "#/definitions/webhook.Recipient"
                },
                "sender": {
                    "$ref": "#/definitions/webhook.Sender"
                },
                "timestamp": {
                    "type": "integer"
                }
            }
        },
        "webhook.Omnichannel": {
            "type": "string",
            "enum": [
                "onx",
                "on5",
                "on4"
            ],
            "x-enum-varnames": [
                "ONX",
                "ON5",
                "ON4"
            ]
        },
        "webhook.QuickReply": {
            "type": "object",
            "required": [
                "payload"
            ],
            "properties": {
                "payload": {
                    "type": "string"
                }
            }
        },
        "webhook.Recipient": {
            "type": "object",
            "required": [
                "id"
            ],
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "webhook.Sender": {
            "type": "object",
            "required": [
                "first_name",
                "id",
                "last_name",
                "profile_pic"
            ],
            "properties": {
                "first_name": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "profile_pic": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Bot Middleware API",
	Description:      "API documentation for the Bot Middleware service.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
