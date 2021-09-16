package patcher

import (
	"github.com/Jeffail/gabs/v2"
	"github.com/ungerik/go-dry"
	"strings"
)

type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Type        string `json:"type"`
	Pattern     string `json:"pattern"`
	Format      string `json:"format"`
}

type Header struct {
	Description string `json:"description"`
	Type        string `json:"type"`
	Format      string `json:"format"`
	Default     string `json:"default"`
	Pattern     string `json:"pattern"`
}

var headers200 = map[string]Header{
	"X-App": {
		Description: "Уникальный идентификатор сервера 1С",
		Type:        "string",
		Format:      "uuid",
		Default:     "2438ac3c-37eb-4902-adef-ed16b4431030",
		Pattern:     "^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$",
	},
	"X-Endpoint": {
		Description: "Уникальный идентификатор точки обмена в соединении с сервером 1С",
		Type:        "string",
		Format:      "uuid",
		Default:     "2438ac3c-37eb-4902-adef-ed16b4431030",
		Pattern:     "^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$",
	},
}

var appHeaderParameter = Parameter{
	Name:        "X-App",
	In:          "header",
	Description: "Уникальный идентификатор сервера 1С",
	Required:    true,
	Type:        "string",
	Pattern:     "^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$",
	Format:      "uuid",
}

var endpointHeaderParameter = Parameter{
	Name:        "X-Endpoint",
	In:          "header",
	Description: "Уникальный идентификатор точки обмена в соединении с сервером 1С",
	Required:    false,
	Type:        "string",
	Pattern:     "^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$",
	Format:      "uuid",
}

func PatchSwagger(json []byte, tagsIgnore ...string) ([]byte, error) {

	jsonParsed, err := gabs.ParseJSON(json)
	if err != nil {
		return nil, err
	}
	// S is shorthand for Search
	for _, child := range jsonParsed.S("paths").ChildrenMap() {

		for _, methodData := range child.ChildrenMap() {

			if !methodData.ExistsP("tags") {
				continue
			}
			var ignore bool
			for _, tag := range methodData.S("tags").Data().([]interface{}) {
				if dry.StringInSlice(strings.ToLower(tag.(string)), tagsIgnore) {
					ignore = true
					break
				}
			}

			if ignore {
				continue
			}

			if !methodData.ExistsP("parameters") {
				methodData.Array("parameters")
			}

			err := methodData.ArrayAppend(appHeaderParameter, "parameters")
			if err != nil {
				return nil, err
			}
			err = methodData.ArrayAppend(endpointHeaderParameter, "parameters")
			if err != nil {
				return nil, err
			}

			if methodData.ExistsP("responses.200") {

				if !methodData.ExistsP("responses.200.headers") {
					_, err := methodData.SetP(headers200, "responses.200.headers")
					if err != nil {
						return nil, err
					}
				}

			}

		}

	}

	// return jsonParsed.MarshalJSON()
	return jsonParsed.EncodeJSON(gabs.EncodeOptIndent("", "  "), gabs.EncodeOptHTMLEscape(true)), nil
}
