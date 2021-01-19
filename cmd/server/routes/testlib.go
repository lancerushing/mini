package routes

import (
	"html/template"
	"os"
	"sync"
	"testing"
)

var _testOneTime sync.Once
var _testSrv Server

func setup(t *testing.T) *Server {
	_testOneTime.Do(func() {

		// @todo is there a better way?
		// When running tests, the working dir is the package dir
		// templates parsing uses paths based on the root dir
		_ = os.Chdir("..")

		// we will init the routes once to help with speed
		// because the route handlers parse their templates
		// during startup, we have ~1 second startup
		// when routes.routes() is called
		_testSrv = Server{}

		_testSrv.layout = template.Must(template.New("test_layout").Parse(`{{ block "main" . }}test layout main{{ end }}s`))
		_testSrv.routes()

	})

	t.Parallel()

	return &_testSrv

}
