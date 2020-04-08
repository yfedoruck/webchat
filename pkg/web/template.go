package web

import (
	"github.com/yfedoruck/webchat/pkg/browser"
	"github.com/yfedoruck/webchat/pkg/env"
	"github.com/yfedoruck/webchat/pkg/fail"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

type TemplateHandler struct {
	once     sync.Once
	Filename string
	templ    *template.Template
}

func (t *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join(env.AppPath()+"/templates", t.Filename)))
	})
	data := map[string]interface{}{"Host": r.Host}
	if authCookie, err := r.Cookie("auth"); err == nil {
		user := browser.Cookie{}
		user.Decode(authCookie.Value)

		data["UserData"] = user
		data["Socket"] = env.Conf().Socket
	}

	err := t.templ.Execute(w, data)
	fail.Check(err)
}
