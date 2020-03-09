package cmd

import (
	"flag"
	"fmt"
	"github.com/hellgate75/go-deploy/io"
	"github.com/hellgate75/go-deploy/net"
	"github.com/hellgate75/go-deploy/types/module"
	"github.com/hellgate75/go-tcp-client/client/proxy"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

var (
	name      string = ""
	loglevel  string = "."
	workdir   string = "."
	modDir    string = ""
	chartsDir string = ""
	configDir string = ""
	systemDir string = ""
	useHosts  string = ""
	useVars   string = ""
	format    string = "YAML"
	env       string = ""
	readTimeout int64 = 0
	fs        *flag.FlagSet
)

const (
	Banner string = `
    ###   ###     ###   #### ###   #      ###  #   #
   #   # #   #    #  #  #    #  #  #     #   #  # #
   #     #   #    #   # #    #   # #     #   #   #
   # ### #   #    #   # ###  ####  #     #   #   #
    #  # #   #    #  #  #    #     #     #   #   #
     ##   ###     ###   #### #     #####  ###    #
`
	Version    string = "v. 1.0.0"
	Authors    string = "Fabrizio Torelli (hellgate75@gmail.com)"
	Disclaimer string = "No Warranty is given on use of this product"
)

func init() {
	name = fmt.Sprintf("deploy-%v", strconv.FormatUint(rand.Uint64(), 10))
	fs = flag.NewFlagSet("go-deploy", flag.PanicOnError)
	fs.StringVar(&name, "name", name, "Deployment unit name")
	fs.StringVar(&workdir, "workDir", ".", "Working directory")
	fs.StringVar(&loglevel, "verbosity", "INFO", "Log Level Verbosity")
	fs.StringVar(&modDir, "modDir", "."+io.GetPathSeparator()+DEFAULT_MODULES_FOLDER, "Go Deploy modules dir")
	fs.StringVar(&chartsDir, "chartsDir", "."+io.GetPathSeparator()+DEFAULT_CHARTS_FOLDER, "Deployment charts folder")
	fs.StringVar(&configDir, "configDir", "."+io.GetPathSeparator()+DEFAULT_CONFIG_FOLDER, "Deployment config folder")
	fs.StringVar(&systemDir, "goDeployDir", userHomeDir()+io.GetPathSeparator()+DEFAULT_SYSTEM_FOLDER, "Go Deploy system folder")
	fs.StringVar(&useHosts, "hosts", "", "Required Hosts files (comma separated file path list)")
	fs.StringVar(&useVars, "vars", "", "Required Vars files (comma separated file path list)")
	fs.StringVar(&format, "language", "YAML", "Config File Language (YAML, XML or JSON)")
	fs.Int64Var(&readTimeout, "readTimeout", 5, "TCP Client Message Read timeout in seconds, used to keep listening for answer from clients")
	fs.StringVar(&env, "env", "", "configuration file env suffix (no default)")
	fs.StringVar(&proxy.PluginLibrariesFolder, "client-prugins-folder", proxy.PluginLibrariesFolder, "Folder where seek for client(s) plugin(s) library [Linux Only]")
	fs.StringVar(&proxy.PluginLibrariesExtension, "client-prugins-extension", proxy.PluginLibrariesExtension, "File extension for client(s) plugin libraries [Linux Only]")
	fs.BoolVar(&proxy.UsePlugins, "use-client-plugins", proxy.UsePlugins, "Enable/disable client(s) plugins [true|false] [Linux Only]")
	fs.StringVar(&net.PluginLibrariesFolder, "prugins-folder", net.PluginLibrariesFolder, "Folder where seek for Go Deploy NET client(s) plugin(s) library [Linux Only]")
	fs.StringVar(&net.PluginLibrariesExtension, "-prugins-extension", net.PluginLibrariesExtension, "File extension for Go Deploy NET client(s) plugin libraries [Linux Only]")
	fs.BoolVar(&net.UsePlugins, "use-plugins", net.UsePlugins, "Enable/disable Go Deploy NET client(s) plugins [true|false] [Linux Only]")
}

func RequiresHelp() bool {
	for _, val := range os.Args {
		if val == "--help" || val == "-help" || val == "-h" || val == "help" {
			fs.Usage()
			return true
		}
	}
	return false
}

func Usage() {
	fs.Usage()
}

func GetTarget() string {
	if len(os.Args) == 2 && os.Args[0][0:1] != "-" {
		return os.Args[0]
	} else if len(os.Args) == 3 && os.Args[0][0:1] != "-" {
		return os.Args[1]
	} else if len(os.Args) > 3 && os.Args[len(os.Args)-2][0:1] != "-" {
		return os.Args[len(os.Args)-1]
	}
	return ""
}

func ParseArguments() (*module.DeployConfig, error) {
	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}
	return &module.DeployConfig{
		DeployName:   name,
		WorkDir:      workdir,
		LogVerbosity: loglevel,
		ConfigDir:    configDir,
		ChartsDir:    chartsDir,
		SystemDir:    systemDir,
		ModulesDir:   modDir,
		UseHosts:     strings.Split(useHosts, ","),
		UseVars:      strings.Split(useVars, ","),
		ConfigLang:   module.DescriptorTypeValue(format),
		EnvSelector:  env,
		ReadTimeout: readTimeout,
	}, nil
}
