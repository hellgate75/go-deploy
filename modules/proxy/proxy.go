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

var UsePlugins bool = false
var PluginLibrariesFolder string = getDefaultPluginsFolder()
var PluginLibrariesExtension = "so"

var Logger log.Logger = nil

type Module interface {
	GetComponent() (meta.Converter, error)
}

type module struct {
	module string
	stub   meta.ProxyStub
}

func (m *module) GetComponent() (meta.Converter, error) {
	return m.stub.Discover(m.module)
}

type Proxy interface {
	DiscoverModule(name string) (Module, error)
}

type proxy struct {
	modules map[string]meta.ProxyStub
}

func (p *proxy) DiscoverModule(name string) (Module, error) {
	Logger.Debugf("module map: %v", p.modules)
	for k, s := range p.modules {
		Logger.Debugf("module map entry: %s", k)
		if k == name {
			return &module{
				module: k,
				stub:   s,
			}, nil
		}
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

func NewProxy() Proxy {
	if modulesMap == nil {
		modulesMap = getModules()
		for key,value := range modulesMap {
			Logger.Debugf("map -> %s = %v", key, value )
		}
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