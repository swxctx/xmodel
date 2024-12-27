package create

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/swxctx/xlog"
	"github.com/swxctx/xmodel/cmd/info"
)

// ModelTpl template file name
const ModelTpl = "__model__tpl__.go"

// CreateProject creates a project.
func CreateProject() {
	xlog.Infof("Generating project: %s", info.ProjPath())

	os.MkdirAll(info.AbsPath(), os.FileMode(0755))
	err := os.Chdir(info.AbsPath())
	if err != nil {
		xlog.Fatalf("[XModel] Jump working directory failed: %v", err)
	}

	// read temptale file
	b, err := ioutil.ReadFile(ModelTpl)
	if err != nil {
		b = []byte(strings.Replace(__tpl__, "__PROJ_NAME__", info.ProjName(), -1))
	}

	// new project code
	proj := NewProject(b)
	proj.Generator()

	// write template file
	f, err := os.OpenFile(ModelTpl, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		xlog.Fatalf("[XModel] Create files error: %v", err)
	}
	defer f.Close()
	f.Write(formatSource(b))
	xlog.Infof("Completed code generation!")
}
