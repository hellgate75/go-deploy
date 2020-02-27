package main

import (
	"fmt"
	"github.com/hellgate75/go-deploy/cmd"
	"github.com/hellgate75/go-deploy/io"
	"github.com/hellgate75/go-deploy/log"
	"github.com/hellgate75/go-deploy/modules"
	"github.com/hellgate75/go-deploy/types/generic"
	"github.com/hellgate75/go-deploy/types/module"
	"github.com/hellgate75/go-deploy/utils"
	"os"
	"strconv"
	"time"
)

var Logger log.Logger = log.NewLogger(log.INFO)

func init() {
	defer func() {
		if r := recover(); r != nil {
			Logger.Error(fmt.Sprintf("Init - Recovery:\n- %v", r))
			os.Exit(1)
		}
	}()
	module.Logger = Logger
	generic.Logger = Logger
	modules.Logger = Logger
	Logger.Println(cmd.Banner)
	Logger.Println("GO DEPLOY " + cmd.Version)
	Logger.Println("Authors:", cmd.Authors)
	Logger.Println(cmd.Disclaimer + "\n")
	Logger.Trace("Init ...")
}

func main() {
	var start time.Time = time.Now()
	defer func() {
		if r := recover(); r != nil {
			Logger.Error(fmt.Sprintf("Recovery:\n- %v", r))
			os.Exit(1)
		}
		Logger.Trace(fmt.Sprint("Exit ..."))
		var end time.Time = time.Now()
		var duration time.Duration = end.Sub(start)
		Logger.Warn(fmt.Sprintf("Total elapsed time: %s", duration.String()))
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
					Logger.Error(fmt.Sprintf("Error: During config files initialization -> <%v>...", errors))
					os.Exit(1)
				}
				var dc *module.DeployConfig = boostrap.GetDeployConfig()
				if dc == nil {
					dc = &module.DeployConfig{}
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
				module.RuntimeDeployConfig = dc
				errB = boostrap.Load(dc.ConfigDir, dc.EnvSelector, dc.ConfigLang, Logger)
				Logger.Info(fmt.Sprintf("Errors during config load: %b", len(errB)))
				if len(errB) > 0 {
					var errors string = ""
					for _, errX := range errB {
						prefix := ""
						if len(errors) > 0 {
							prefix = "\n"
						}
						errors += prefix + "- " + errX.Error()
					}
					Logger.Error(fmt.Sprintf("Error: During config files load -> <%v>...", errors))
					os.Exit(1)
				}
				var dt *module.DeployType = boostrap.GetDeployType()
				if dt == nil {
					dt = &module.DeployType{}
				}
				dt = boostrap.GetDefaultDeployType().Merge(dt)
				module.RuntimeDeployType = dt
				var nt *module.NetProtocolType = boostrap.GetNetType()
				if nt == nil {
					nt = &module.NetProtocolType{}
				}
				nt = boostrap.GetDefaultNetType().Merge(nt)
				module.RuntimeNetworkType = nt
				Logger.Info(fmt.Sprintf("Configuration Summary: \nDeploy Config: %v\nDeployType: %v\nNetType: %v\n", dc.String(), dt.String(), nt.String()))
				if dt.DeploymentType == module.FILE_SOURCE {
					var filePath string = dc.WorkDir + io.GetPathSeparator() + target
					Logger.Warn(fmt.Sprintf("Trying load of Feed: %s\n", filePath))
					var feed generic.IFeed = generic.NewFeed("default")
					err = feed.Load(filePath)
					if err != nil {
						Logger.Error(fmt.Sprintf("Error trying to load Feed for file: %s -> Details: \n%s", filePath, err.Error()))
						os.Exit(1)
					}
					feedEx, errValList := feed.Validate()
					if len(errValList) > 0 {
						var errors string = ""
						for _, errX := range errValList {
							prefix := ""
							if len(errors) > 0 {
								prefix = "\n"
							}
							errors += prefix + "- " + errX.Error()
						}
						Logger.Error(fmt.Sprintf("Error trying to validate Feed for file: %s -> Details: \n%s", filePath, errors))
						os.Exit(1)
					}
					if len(feedEx.Steps) > 0 {
						Logger.Warn(fmt.Sprintf("Reading file: %s, discovered %s main steps!!", filePath, strconv.Itoa(len(feedEx.Steps))))
						errExList := boostrap.Run(feedEx, Logger)
						if len(errExList) > 0 {
							var errors string = ""
							for _, errX := range errExList {
								prefix := ""
								if len(errors) > 0 {
									prefix = "\n"
								}
								errors += prefix + errX.Error()
							}
							Logger.Error(fmt.Sprintf("Error: During deploy execution -> <%v>...", errors))
							os.Exit(1)
						}
					} else {
						Logger.Warn(fmt.Sprintf("Unable to find any command in the given file: %s", filePath))
						Logger.Warn("Nothing to do here!!")
					}
					Logger.Warn("Deploy procedure complete!!")
				} else {
					Logger.Warn(fmt.Sprintf("Feature %v NOT IMPLEMENTED yet!!", dt.DeploymentType))
				}
			}
		}
	} else {
		os.Exit(0)
	}
}
