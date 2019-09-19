package server

import (
	"html/template"
	"os"
	"sync"
	"testing"
)

var _testOneTime sync.Once
var _testSrv server


func setup(t *testing.T) *server {
	_testOneTime.Do(func() {

		// @todo is there a better way?
		// When running tests, the working dir is the package dir
		// templates parsing uses paths based on the root dir
		_ = os.Chdir("..")


		// we will init the server once to help with speed
		// because the route handlers parse their templates
		// during startup, we have ~1 second startup
		// when server.routes() is called
		_testSrv = server{}

		_testSrv.layout = template.Must(template.New("test_layout").Parse(`{{ block "main" . }}test layout main{{ end }}s`))
		_testSrv.routes()

	})

	t.Parallel()

	return &_testSrv

}
