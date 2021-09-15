package rapidoc

type RenderStyle string
type SchemaStyle string
type ThemeStyle string
type LayoutStyle string

const (
	RenderStyle_Read  RenderStyle = "read"
	RenderStyle_View  RenderStyle = "view"
	RenderStyle_Focus RenderStyle = "focused"

	SchemaStyle_Tree  SchemaStyle = "tree"
	SchemaStyle_Table SchemaStyle = "table"

	Theme_Dark  ThemeStyle = "dark"
	Theme_Light ThemeStyle = "light"

	Layout_Row    LayoutStyle = "row"
	Layout_Column LayoutStyle = "column"
)

// URL presents the url pointing to API definition (normally swagger.json or swagger.yaml).
func URL(url string) func(c *Config) {
	return func(c *Config) {
		c.SpecURL = url
	}
}

func Style(style RenderStyle) func(c *Config) {
	return func(c *Config) {
		c.RenderStyle = style
	}
}
func Layout(layout LayoutStyle) func(c *Config) {
	return func(c *Config) {
		c.Layout = layout
	}
}

func Schema(schema SchemaStyle) func(c *Config) {
	return func(c *Config) {
		c.SchemaStyle = schema
	}
}

type Config struct {
	Title       string      `json:"tiltle,omitempty"`
	SpecURL     string      `json:"spec_url,omitempty"`
	HeaderText  string      `json:"header_text,omitempty"`
	LogoURL     string      `json:"logo_url,omitempty"`
	RenderStyle RenderStyle `json:"render_style,omitempty"`
	SchemaStyle SchemaStyle `json:"schema_style,omitempty"`
	Theme       ThemeStyle  `json:"theme,omitempty"`
	Layout      LayoutStyle `json:"layout,omitempty"`
}

func GetDefaultRapiDocConfig() Config {
	return Config{
		Title:       "API Documentation",
		SpecURL:     "./swagger.json",
		HeaderText:  "API Documentation",
		LogoURL:     "https://mrin9.github.io/RapiDoc/images/logo.png",
		RenderStyle: RenderStyle_Read,
		SchemaStyle: SchemaStyle_Tree,
		Theme:       Theme_Dark,
		Layout:      Layout_Row,
	}
}

func HtmlTemplateRapiDoc() string {
	return `<!doctype html>
	<html>	
	<head>
		<title>{{$.Title}}</title>
		<meta charset="utf-8">
		<link href="https://fonts.googleapis.com/css2?family=Sarabun&display=swap" rel="stylesheet">
		<link href="https://fonts.googleapis.com/css2?family=Open+Sans:wght@300;600&family=Roboto+Mono&display=swap" rel="stylesheet">
		<script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
	</head>
	
	<body>
		<rapi-doc 
			spec-url="{{.SpecURL}}" 
			heading-text="{{.HeaderText}}" 
			theme="{{.Theme}}"
			layout="{{.Layout}}"
			use-path-in-nav-bar=true
			sort-endpoints-by="path"
			regular-font="Sarabun" 
			mono-font="'Roboto Mono'" 
			render-style="{{.RenderStyle}}" 
			schema-style="{{.SchemaStyle}}">
			<div slot="nav-logo" style="display: flex; align-items: center; justify-content: center;">
				<img src="{{.LogoURL}}" style="width:150px">
			</div>
		</rapi-doc>
	</body>
	</html>`
}
