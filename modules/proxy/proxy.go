package proxy

import (
	"errors"
	"fmt"
	mods "github.com/hellgate75/go-deploy-modules/modules"
	"github.com/hellgate75/go-tcp-common/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"
	"strings"
	"github.com/hellgate75/go-deploy/modules/meta"
)

// Use custom plugins loading proxies
var UsePlugins bool = false
// Use custom plugins folder to seek for libraries
var PluginLibrariesFolder string = getDefaultPluginsFolder()
// Assume this extension name for ;loading the libraries (we hope in future windows will allow plugins)
var PluginLibrariesExtension = "so"

var Logger log.Logger = nil

// Define Behaviors of a Module Component
type Module interface {
	// Retrieve module meta.Converter component
	GetComponent() (meta.Converter, error)
}

type module struct {
	module string
	stub   meta.ProxyStub
}

func (m *module) GetComponent() (meta.Converter, error) {
	return m.stub.Discover(m.module)
}

// Define Behaviors of a Proxy Component
type Proxy interface {
	// Discover a Module by given command name into own Components list
	DiscoverModule(name string) (Module, error)
}

type proxy struct {
	modules map[string]meta.ProxyStub
}

func (p *proxy) DiscoverModule(name string) (Module, error) {
	Logger.Debugf("module map: %v", p.modules)
	if stub, ok := p.modules[name]; ok {
		return &module{
			module: name,
			stub:   stub,
		}, nil
	}
	return nil, errors.New(fmt.Sprintf("Unable to discover module: %s", name))
}

func getModules() map[string]meta.ProxyStub {
	var outMap map[string]meta.ProxyStub = make(map[string]meta.ProxyStub)
	if UsePlugins {
		Logger.Debug("modules.proxy.GetSender() -> Loading library for map modules")
		forEachModulesMapsInPlugins(func(modulesMapsList []map[string]meta.ProxyStub) {
			if len(modulesMapsList) > 0 {
				for _, mapX := range modulesMapsList {
					for name, stub := range mapX {
						outMap[name] = stub
					}
				}
			}
		})
	}
	for name, stub := range mods.GetModulesMap() {
		outMap[name] = stub
	}
	return outMap
}

func filterByExtension(fileName string) bool {
	n := len(PluginLibrariesExtension)
	fileNameLen := len(fileName)
	posix := fileNameLen - n
	return posix > 0 && strings.ToLower(fileName[posix:]) == strings.ToLower("." + PluginLibrariesExtension)
}

func listLibrariesInFolder(dirName string) []string {
	var out []string = make([]string, 0)
	_, err0 := os.Stat(dirName)
	if err0 == nil {
		lst, err1 := ioutil.ReadDir(dirName)
		if err1 == nil {
			for _,file := range lst {
				if file.IsDir() {
					fullDirPath := dirName + string(os.PathSeparator) + file.Name()
					newList := listLibrariesInFolder(fullDirPath)
					out = append(out, newList...)
				} else {
					if filterByExtension(file.Name()) {
						fullFilePath := dirName + string(os.PathSeparator) + file.Name()
						out = append(out, fullFilePath)
						
					}
				}
			}
		}
	}
	return out
}

func forEachModulesMapsInPlugins(callback func([]map[string]meta.ProxyStub)())  {
	var modulesMaps []map[string]meta.ProxyStub = make([]map[string]meta.ProxyStub, 0)
	dirName := PluginLibrariesFolder
	_, err0 := os.Stat(dirName)
	if err0 == nil {
		libraries := listLibrariesInFolder(dirName)
		for _,libraryFullPath := range libraries {
			Logger.Debugf("modules.proxy.forEachSenderInPlugins() -> Loading help from library: %s", libraryFullPath)
			plugin, err := plugin.Open(libraryFullPath)
			if err == nil {
				sym, err2 := plugin.Lookup("GetModulesMap")
				if err2 != nil {
					modules := sym.(func()(map[string]meta.ProxyStub))()
					modulesMaps = append(modulesMaps, modules)
				}
			}
		}
	}
	callback(modulesMaps)
}

var modulesMap map[string]meta.ProxyStub =nil

// Creates a New Proxy filled wit all available Built-In and Custom Modules, loaded just on first call
func NewProxy() Proxy {
	if modulesMap == nil {
		modulesMap = getModules()
	}
	return &proxy{
		modules: modulesMap,
	}
}

func getDefaultPluginsFolder() string {
	execPath, err := os.Executable()
	if err != nil {
		pwd, errPwd := os.Getwd()
		if errPwd != nil {
			return filepath.Dir(".") + string(os.PathSeparator) + "modules"
		}
		return filepath.Dir(pwd) + string(os.PathSeparator) + "modules"
	}
	return filepath.Dir(execPath) + string(os.PathSeparator) + "modules"
}