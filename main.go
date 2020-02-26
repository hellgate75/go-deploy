package main

import (
	"fmt"
	"github.com/hellgate75/go-deploy/cmd"
	"github.com/hellgate75/go-deploy/types"
	"log"
	"os"
)

var Logger *log.Logger = log.New(os.Stdout, "[go-deploy] ", log.LstdFlags|log.LUTC)

func init() {
	Logger.Println("Init ...")

}

func main() {
	defer func() {
		if r := recover(); r != nil {
			Logger.Println(fmt.Sprintf("Recovery:\n- %v", r))
		}
		Logger.Println(fmt.Sprint("Exit ..."))
		os.Exit(0)
	}()
	fmt.Println(cmd.Banner)
	fmt.Println("GO DEPLOY " + cmd.Version)
	fmt.Println("Authors:", cmd.Authors)
	fmt.Println(cmd.Disclaimer + "\n")
	if !cmd.RequiresHelp() {
		Logger.Println(fmt.Sprint("Main ..."))
		config, err := cmd.ParseArguments()
		if err != nil {
			Logger.Println(fmt.Sprintf("Error: %v", err))
			cmd.Usage()
		} else {
			var target string = cmd.GetTarget()
			if target == "" {
				Logger.Println("Error: No target defined")
				cmd.Usage()
			} else {
				var boostrap cmd.Bootstrap = cmd.NewBootStrap()
				errB := boostrap.Load(config.ConfigDir, config.EnvSelector, config.ConfigLang, Logger)
				Logger.Println(fmt.Sprintf("Errors during config load: %b", len(errB)))
				if len(errB) > 0 {
					var errors string = ""
					for _, errX := range errB {
						prefix := ""
						if len(errors) > 0 {
							prefix = "\n"
						}
						errors += prefix + errX.Error()
					}
					Logger.Println(fmt.Sprintf("Error: During config files load <%v>...", errors))
					os.Exit(1)
				}
				var dc *types.DeployConfig = boostrap.GetDeployConfig()
				if dc == nil {
					dc = &types.DeployConfig{}
				}
				dc = dc.Merge(config)
				var dt *types.DeployType = boostrap.GetDeployType()
				if dt == nil {
					dt = &types.DeployType{}
				}
				dt = boostrap.GetDefaultDeployType().Merge(dt)
				var nt *types.NetProtocolType = boostrap.GetNetType()
				if nt == nil {
					nt = &types.NetProtocolType{}
				}
				nt = boostrap.GetDefaultNetType().Merge(nt)
				Logger.Println(fmt.Sprintf("Configuration Summary: \nDeploy Config: %v\nDeployType: %v\nNetType: %v\n", dc.String(), dt.String(), nt.String()))
			}
		}
	} else {
		os.Exit(0)
	}
}
