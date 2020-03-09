package net

import (
	"github.com/hellgate75/go-deploy-clients/proxy"
	"github.com/hellgate75/go-tcp-common/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"
	"strings"
	"github.com/hellgate75/go-deploy/net/generic"
)

var Logger log.Logger = nil

var UsePlugins bool = false
var PluginLibrariesFolder string = getDefaultPluginsFolder()
var PluginLibrariesExtension = "so"

func DiscoverConnectionHandler(clientName string) (generic.NewConnectionHandlerFunc, error) {
	if UsePlugins {
		Logger.Debugf("client.proxy.GetSender() -> Loading library for command: %s", clientName)
		var handler generic.NewConnectionHandlerFunc = nil
		forEachConnectionFactoryInPlugins(clientName, func(handlersList []generic.NewConnectionHandlerFunc) {
			if len(handlersList) > 0 {
				handler = handlersList[0]
			}
		})
		if handler != nil {
			return handler, nil
		}
	}
	return proxy.GetConnectionHandlerFactory(clientName)
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

func forEachConnectionFactoryInPlugins(clientName string, callback func([]generic.NewConnectionHandlerFunc)())  {
	var handlers []generic.NewConnectionHandlerFunc = make([]generic.NewConnectionHandlerFunc, 0)
	dirName := PluginLibrariesFolder
	_, err0 := os.Stat(dirName)
	if err0 == nil {
		libraries := listLibrariesInFolder(dirName)
		for _,libraryFullPath := range libraries {
			Logger.Debugf("net.forEachSenderInPlugins() -> Loading help from library: %s", libraryFullPath)
			plugin, err := plugin.Open(libraryFullPath)
			if err == nil {
				sym, err2 := plugin.Lookup("GetConnectionHandlerFactory")
				if err2 != nil {
					handler, errPlugin := sym.(func(string)(generic.NewConnectionHandlerFunc, error))(clientName)
					if errPlugin != nil {
						continue
					}
					//handler.SetLogger(Logger)
					handlers = append(handlers, handler)
				}
			}
		}
	}
	callback(handlers)
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