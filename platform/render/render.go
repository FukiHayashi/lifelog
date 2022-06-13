// render.go
package render

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/contrib/renders/multitemplate"
)

func CreateRenderTemplates() multitemplate.Render {
	r := multitemplate.New()

	htmls, err := filepath.Glob("./views/*.html")
	if err != nil {
		log.Fatal(err)
	}
	for _, html := range htmls {
		slice := strings.Split(html, "/")
		r.AddFromFiles(slice[len(slice)-1], "./views/layout.tmpl", html)
	}
	return r
}
