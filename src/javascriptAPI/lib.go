package javascriptAPI

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"dejavuDB/src/config"
	_ "embed"
	"io"
	"os"
	"path"
	"sync"

	"github.com/dop251/goja"
)

var javascript_API_lib = map[string]javascript_module{}
var javascript_API_lib_lock = sync.RWMutex{}

type javascript_module struct {
	name         string
	version      string
	version_info string
	auther       string
	description  string
	license      string
	model_path   string

	is_in_ram bool
	program   *goja.Program
	script    string
	enabled   bool
}

func NewModule(targz []byte) error {
	greader, err := gzip.NewReader(bytes.NewReader(targz))
	if err != nil {
		return err
	}
	tarreader := tar.NewReader(greader)
	header, err := tarreader.Next()
	for err != nil {

		p := path.Join(config.RootDir, "modules", header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			err := os.Mkdir(p, 0755)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(p)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarreader); err != nil {
				return err
			}
			outFile.Close()
		case tar.TypeSymlink:

		default:
		}

		header, err = tarreader.Next()
	}
	return nil
}

func init() {
	javascript_API_lib["http"] = javascript_module{
		name:         "http",
		version:      "0.0.1",
		version_info: "experiment",
		auther:       "YC",
		description:  "javascript wrapper for go http",
		model_path:   "builtin",

		is_in_ram: true,
		enabled:   false, // default disabled for security
	}

	javascript_API_lib["fmt"] = javascript_module{
		name:         "fmt",
		version:      "0.0.1",
		version_info: "expiriment",
		auther:       "YC",
		description:  "javascript wrapper fo go fmt",
		model_path:   "builtin",

		is_in_ram: true,
		enabled:   true,
	}
}
