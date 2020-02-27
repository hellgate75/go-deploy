package main

import (
	"fmt"
	"github.com/hellgate75/go-deploy/cmd"
	"github.com/hellgate75/go-deploy/io"
	"github.com/hellgate75/go-deploy/log"
	"github.com/hellgate75/go-deploy/types"
	"github.com/hellgate75/go-deploy/utils"
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
		Logger.Info(fmt.Sprintf("Logger initial Verbosity : %v", Logger.GetVerbosity()))
		Logger.Trace(fmt.Sprint("Main ..."))
		config, err := cmd.ParseArguments()
		if config.LogVerbosity != "" && config.LogVerbosity != string(Logger.GetVerbosity()) {
			Logger.SetVerbosity(log.VerbosityLevelFromString(config.LogVerbosity))
			Logger.Info(fmt.Sprintf("Logger Verbosity Setted up to : %v", Logger.GetVerbosity()))
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
				config.WorkDir = utils.FixFolder(config.WorkDir, io.GetCurrentFolder(), "")
				config.ConfigDir = utils.FixFolder(config.ConfigDir, config.WorkDir, cmd.DEPLOY_CONFIG_FILE_NAME)

				errB := boostrap.Init(config.ConfigDir, config.EnvSelector, config.ConfigLang, Logger)
				Logger.Info(fmt.Sprintf("Errors during config init: %b", len(errB)))
				if len(errB) > 0 {
					var errors string = ""
					for _, errX := range errB {
						prefix := ""
						if len(errors) > 0 {
							prefix = "\n"
						}
						errors += prefix + errX.Error()
					}
					Logger.Error(fmt.Sprintf("Error: During config files initialization <%v>...", errors))
					os.Exit(1)
				}
				var dc *types.DeployConfig = boostrap.GetDeployConfig()
				if dc == nil {
					dc = &types.DeployConfig{}
				}
				if dc.DeployName != "" {
					config.DeployName = dc.DeployName
				}
				dc = dc.Merge(config)
				dc.WorkDir = utils.FixFolder(dc.WorkDir, io.GetCurrentFolder(), "")
				dc.ConfigDir = utils.FixFolder(dc.ConfigDir, dc.WorkDir, cmd.DEPLOY_CONFIG_FILE_NAME)
				if dc.LogVerbosity != "" && dc.LogVerbosity != string(Logger.GetVerbosity()) {
					Logger.SetVerbosity(log.VerbosityLevelFromString(dc.LogVerbosity))
					Logger.Info(fmt.Sprintf("Logger Verbosity Setted up to : %v", Logger.GetVerbosity()))
				}
				errB = boostrap.Load(dc.ConfigDir, dc.EnvSelector, dc.ConfigLang, Logger)
				Logger.Info(fmt.Sprintf("Errors during config load: %b", len(errB)))
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
