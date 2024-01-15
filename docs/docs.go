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
        "/encryption/decrypt_file": {
            "post": {
                "description": "Decrypt a file",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Encryption"
                ],
                "summary": "Decrypt file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "File",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "file_decrypted",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "bad_request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal_server_error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/encryption/decrypt_saved_file": {
            "post": {
                "description": "Decrypt a file stored in server",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "File"
                ],
                "summary": "Decrypt a file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "File",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "file_decrypted",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "bad_request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "not_found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal_server_error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/encryption/generate_keys": {
            "post": {
                "description": "Generate a pair of keys",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Encryption"
                ],
                "summary": "Generate keys",
                "responses": {
                    "200": {
                        "description": "keys_generated",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal_server_rror",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/encryption/upload_encrypted_file": {
            "post": {
                "description": "Upload a file",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Encryption"
                ],
                "summary": "Upload encrypted file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "File",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "file_uploaded",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "bad_request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal_server_rror",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/encryption/upload_file": {
            "post": {
                "description": "Upload a file to encrypt",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Encryption"
                ],
                "summary": "Upload file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "File",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "file_uploaded",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "bad_request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal_server_rror",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/file": {
            "get": {
                "description": "Get a list of all encrypted files",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "File"
                ],
                "summary": "Get all encrypted files",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ListFilesResponse"
                        }
                    },
                    "500": {
                        "description": "internal_server_error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/file/by_name": {
            "get": {
                "description": "Provide the file data",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "File"
                ],
                "summary": "Find file by name",
                "parameters": [
                    {
                        "type": "string",
                        "description": "filename to find",
                        "name": "filename",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ShowFileResponse"
                        }
                    },
                    "400": {
                        "description": "bad_request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "not_found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal_server_error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/file/by_username": {
            "get": {
                "description": "Get a list of encrypted files by username",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "File"
                ],
                "summary": "Get encrypted files by username",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ListFilesResponse"
                        }
                    },
                    "400": {
                        "description": "bad_request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal_server_error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/file/delete_file": {
            "post": {
                "description": "deletes a file",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "File"
                ],
                "summary": "Delete file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "File",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "file_deleted",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "bad_request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "not_found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal_server_error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user": {
            "get": {
                "description": "Get all users",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get all users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ListUsersResponse"
                        }
                    },
                    "400": {
                        "description": "bad_request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal_server_rror",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Delete user",
                "responses": {
                    "200": {
                        "description": "user_deleted",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "bad_request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal_server_rror",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/username": {
            "get": {
                "description": "Provide the user data",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Find user by username",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User` + "`" + `s username to find",
                        "name": "username",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.ShowUserResponse"
                        }
                    },
                    "400": {
                        "description": "bad_request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "not_found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal_server_error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Address": {
            "type": "object",
            "properties": {
                "city": {
                    "type": "string"
                },
                "state": {
                    "type": "string"
                },
                "street": {
                    "type": "string"
                },
                "zipcode": {
                    "type": "string"
                }
            }
        },
        "model.Contact": {
            "type": "object",
            "properties": {
                "celphone": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                }
            }
        },
        "schemas.EncryptedFileResponse": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "description": "AnonymizedFile AnonymizedFile     ` + "`" + `bson:\"anonymized_file\"` + "`" + `",
                    "type": "string"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "deletedAt": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "schemas.ListFilesResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/schemas.EncryptedFileResponse"
                    }
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "schemas.ListUsersResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/schemas.UserResponse"
                    }
                }
            }
        },
        "schemas.ShowFileResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/schemas.EncryptedFileResponse"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "schemas.ShowUserResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/schemas.UserResponse"
                }
            }
        },
        "schemas.UserResponse": {
            "type": "object",
            "properties": {
                "address": {
                    "$ref": "#/definitions/model.Address"
                },
                "contact": {
                    "$ref": "#/definitions/model.Contact"
                },
                "cpf": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "dateOfBirth": {
                    "type": "string"
                },
                "deletedAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Agent API",
	Description:      "Agent Operations",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
