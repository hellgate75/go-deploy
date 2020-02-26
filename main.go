package main

import (
	"fmt"
	"github.com/hellgate75/go-deploy/cmd"
	"github.com/hellgate75/go-deploy/log"
	"github.com/hellgate75/go-deploy/types"
	"os"
)

var Logger log.Logger = log.NewLogger(log.INFO)

func init() {
	defer func() {
		if r := recover(); r != nil {
			Logger.Error(fmt.Sprintf("Init - Recovery:\n- %v", r))
			os.Exit(1)
		}
	}()
	Logger.Println(cmd.Banner)
	Logger.Println("GO DEPLOY " + cmd.Version)
	Logger.Println("Authors:", cmd.Authors)
	Logger.Println(cmd.Disclaimer + "\n")
	Logger.Trace("Init ...")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			Logger.Error(fmt.Sprintf("Recovery:\n- %v", r))
			os.Exit(1)
		}
		Logger.Trace(fmt.Sprint("Exit ..."))
		os.Exit(0)
	}()
	if !cmd.RequiresHelp() {
		Logger.Trace(fmt.Sprint("Main ..."))
		config, err := cmd.ParseArguments()
		if config.LogVerbosity != "" {
			Logger.SetVerbosity(log.LogLevel(config.LogVerbosity))
			Logger.Info("Logger verbosity is setted to : " + string(Logger.GetVerbosity()))
		}
		if err != nil {
			Logger.Error(fmt.Sprintf("Error: %v", err))
			cmd.Usage()
		} else {
			var target string = cmd.GetTarget()
			if target == "" {
				Logger.Error("Error: No target defined")
				cmd.Usage()
				os.Exit(1)
			} else {

				var boostrap cmd.Bootstrap = cmd.NewBootStrap()
				errB := boostrap.Load(config.ConfigDir, config.EnvSelector, config.ConfigLang, Logger)
				Logger.Error(fmt.Sprintf("Errors during config load: %b", len(errB)))
				if len(errB) > 0 {
					var errors string = ""
					for _, errX := range errB {
						prefix := ""
						if len(errors) > 0 {
							prefix = "\n"
						}
						errors += prefix + errX.Error()
					}
					Logger.Error(fmt.Sprintf("Error: During config files load <%v>...", errors))
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
				Logger.Info(fmt.Sprintf("Configuration Summary: \nDeploy Config: %v\nDeployType: %v\nNetType: %v\n", dc.String(), dt.String(), nt.String()))
			}
		}
	} else {
		os.Exit(0)
	}
}
