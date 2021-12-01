package javascriptAPI

import (
	"path"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"

	"dejavuDB/src/config"
	"dejavuDB/src/javascriptAPI/native/fmt"
	"dejavuDB/src/javascriptAPI/native/http"
	"dejavuDB/src/javascriptAPI/native/tar"
)

type JsOptions struct {
	Adm      bool
	UserName string
	uid      *uint64
	gid      *uint64
}

var jsRegister *require.Registry

func JavascriptRun(script string, options JsOptions) ([]byte, error) {

	vm := goja.New()

	if jsRegister == nil {
		jsRegister = require.
			NewRegistry(require.WithGlobalFolders(path.Join(config.RootDir, "modules")))

		jsRegister.RegisterNativeModule("http", http.HttpModuleLoader)
		jsRegister.RegisterNativeModule("fmt", fmt.FmtModuleLoader)
		jsRegister.RegisterNativeModule("tar", tar.TarModuleLoader)
	}

	db := NewVmDatabase(vm)

	p, err := vm.RunScript("transection.js", script)
	if err != nil {
		return nil, err
	}

	db.storage.Commit()
	return []byte(p.String()), nil
}
