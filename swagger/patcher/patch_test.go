package patcher

import (
	"reflect"
	"testing"
)

func TestPatchSwagger(t *testing.T) {

	tests := []struct {
		name    string
		json    []byte
		want    []byte
		wantErr bool
	}{
		{
			"append parameters",
			[]byte(`{
  "paths": {
    "/api/v1/apps": {
      "get": {
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/RegisterRequest"
            }
          }
        ]
      }
    }
  }
}`),
			[]byte(`{
  "paths": {
    "/api/v1/apps": {
      "get": {
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/RegisterRequest"
            }
          }
        ]
      }
    }
  }
}`),
			false,
		},
		{
			"ignore tags parameters",
			[]byte(`{
  "paths": {
    "/api/v1/apps": {
      "get": {
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/RegisterRequest"
            }
          }
        ],
		"tags": [
			"users"
		]
      }
    }
  }
}`),
			[]byte(`{
  "paths": {
    "/api/v1/apps": {
      "get": {
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/RegisterRequest"
            }
          }
        ]
      }
    }
  }
}`),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PatchSwagger(tt.json)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchSwagger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatchSwagger() got = %v, want %v", string(got), tt.want)
			}
		})
	}
}
