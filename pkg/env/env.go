package env

import (
	"encoding/json"
	"github.com/yfedoruck/webchat/pkg/fail"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var port *string

func Port() string {
	if port != nil {
		return *port
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	return port
}

var baseDir *string

// Depends on dir script
func AppPath() string {
	if baseDir != nil {
		return *baseDir
	}
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("Caller error")
	}

	env := filepath.Dir(b)
	pkg := filepath.Dir(env)
	app := filepath.Dir(pkg)
	baseDir = &app
	return *baseDir
}

type addrConf struct {
	Host   string `json:"Host"`
	Socket string `json:"socket"`
}

var addr *addrConf

func Conf() *addrConf {
	if addr != nil {
		return addr
	}
	file, err := os.Open(AppPath() + filepath.FromSlash("/config/"+domain()+"/addr.json"))
	fail.Check(err)

	addr = &addrConf{}
	err = json.NewDecoder(file).Decode(addr)
	fail.Check(err)

	return addr
}

func domain() string {
	domain := os.Getenv("USERDOMAIN")
	if domain == "home" {
		return "local"
	}
	return "heroku"
}
