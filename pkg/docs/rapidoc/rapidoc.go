package rapidoc

import (
	"bytes"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

func New(configFns ...func(*Config)) http.HandlerFunc {

	cfg := GetDefaultRapiDocConfig()

	for _, configFn := range configFns {
		configFn(&cfg)
	}

	tpl, err := template.New("rapidoc").Parse(HtmlTemplateRapiDoc())
	if err != nil {
		panic(err)
	}

	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, cfg); err != nil {
		panic(err)
	}

	var re = regexp.MustCompile(`^(.*/)([^?].*)?[?|.]*$`)

	return func(w http.ResponseWriter, r *http.Request) {

		matches := re.FindStringSubmatch(r.RequestURI)
		path := matches[2]
		prefix := matches[1]

		switch {
		case path == "index.html":
			_, err := buf.WriteTo(w)
			if err != nil {
				return
			}
		case strings.HasSuffix(path, ".json") ||
			strings.HasSuffix(path, ".yaml") ||
			strings.HasSuffix(path, ".yml"):
			if strings.HasSuffix(path, ".json") {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
			} else {
				w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
			}
			// return c.SendFile(p, true)
		case path == "":
			http.Redirect(w, r, prefix+"index.html", 301)
		default:
			http.NotFound(w, r)
		}
	}

}
