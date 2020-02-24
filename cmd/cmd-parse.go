package cmd

import (
	"os"
	"fmt"
	"strconv"
	"math/rand"
	"strings"
	"flag"
	"github.com/hellgate75/go-deploy/types"
)

var (
	name string
	dir string
	useHosts string
	useVars string
	fs *flag.FlagSet
)

func init() {
	fs = flag.NewFlagSet("go-deploy", flag.PanicOnError)
	fs.StringVar(&name, "name", fmt.Sprintf("deploy-%v", strconv.FormatUint(rand.Uint64(), 10)), "Deployment unit name");
	fs.StringVar(&dir, "dir", ".deploy", "Deployment config folder");
	fs.StringVar(&useHosts, "hosts", "", "Required Hosts files (comma separated file path list)");
	fs.StringVar(&useVars, "vars", "", "Required Vars files (comma separated file path list)");
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

func ParseArguments() (*types.DeployConfig, error) {
	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}
	return &types.DeployConfig {
		DeployName: name,
		ConfigDir: dir,
		UseHosts: strings.Split(useHosts, ","),
		UseVars: strings.Split(useVars, ","),
	}, nil
}
