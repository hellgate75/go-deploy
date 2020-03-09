package main

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/hellgate75/go-tcp-client/client/worker"
	clicommon "github.com/hellgate75/go-tcp-client/common"
	"github.com/hellgate75/go-tcp-common/log"
	clientlog "github.com/hellgate75/go-tcp-common/log"
	"os"
	"strconv"
	"time"
	ngen "github.com/hellgate75/go-deploy/net/generic"
	modproxy "github.com/hellgate75/go-deploy/modules/proxy"
	"github.com/hellgate75/go-deploy/cmd"
	"github.com/hellgate75/go-deploy/io"
	"github.com/hellgate75/go-deploy/modules"
	"github.com/hellgate75/go-deploy/net"
	"github.com/hellgate75/go-deploy/types/generic"
	"github.com/hellgate75/go-deploy/types/module"
	"github.com/hellgate75/go-deploy/utils"
)

var Logger log.Logger = nil

func init() {
	defer func() {
		if r := recover(); r != nil {
			Logger.Errorf("Init - Recovery:\n- %v", r)
			os.Exit(1)
		}
	}()
	Logger = log.NewLogger("go-deploy", log.INFO)
	setupLogger()
	printInfo()
}

func setupLogger() {
	defer func() {
		if r := recover(); r != nil {
			Logger.Errorf("SetUpLogger - Recovery:\n- %v", r)
			os.Exit(1)
		}
	}()
	module.Logger = Logger
	generic.Logger = Logger
	modules.Logger = Logger
	cmd.Logger = Logger
	ngen.Logger = Logger
	modproxy.Logger = Logger
	net.Logger = Logger
	Logger.Trace("Init ...")
	worker.Logger.AffiliateTo(Logger)
	
}

func printInfo() {
	defer func() {
		if r := recover(); r != nil {
			Logger.Errorf("PrintInfo - Recovery:\n- %v", r)
			os.Exit(1)
		}
	}()
	Logger.Println(color.LightGreen.Render(cmd.Banner))
	Logger.Println(color.LightYellow.Render("GO DEPLOY " + cmd.Version))
	Logger.Println("Author: ", color.LightYellow.Render(cmd.Authors))
	Logger.Println(cmd.Disclaimer + "\n")
}

func main() {
	var start time.Time = time.Now()
	var help bool = false
	defer func() {
		var exitCode int = 0
		if r := recover(); r != nil {
			Logger.Error(fmt.Sprintf("Recovery:\n- %v", r))
			exitCode = 1
		}
		Logger.Trace(fmt.Sprint("Exit ..."))
		if !help {
			var end time.Time = time.Now()
			var duration time.Duration = end.Sub(start)
			Logger.Warnf("Total elapsed time: %s", duration.String())
		}
		os.Exit(exitCode)
	}()
	if !cmd.RequiresHelp() {
		Logger.Infof("Logger initial Verbosity : %v", Logger.GetVerbosity())
		Logger.Trace("Main ...")
		config, err := cmd.ParseArguments()
		if config.LogVerbosity != "" && config.LogVerbosity != string(Logger.GetVerbosity()) {
			Logger.SetVerbosity(log.VerbosityLevelFromString(config.LogVerbosity))
			//worker.Logger.SetVerbosity(clientlog.VerbosityLevelFromString(config.LogVerbosity))
			Logger.Infof("Logger Verbosity Setted up to : %v", Logger.GetVerbosity())
		}
		if err != nil {
			Logger.Errorf("Error: %v", err)
			cmd.Usage()
		} else {
			var target string = cmd.GetTarget()
			if target == "" {
				cmd.Usage()
				panic("Error: No target defined")
			} else {
				var boostrap cmd.Bootstrap = cmd.NewBootStrap()
				config.WorkDir = utils.FixFolder(config.WorkDir, io.GetCurrentFolder(), "")
				config.ConfigDir = utils.FixFolder(config.ConfigDir, config.WorkDir, cmd.DEPLOY_CONFIG_FILE_NAME)

				errB := boostrap.Init(config.ConfigDir, config.EnvSelector, config.ConfigLang, Logger)
				Logger.Debugf("Errors during config init: %v", len(errB))
				if len(errB) > 0 {
					var errors string = ""
					for _, errX := range errB {
						prefix := ""
						if len(errors) > 0 {
							prefix = "\n"
						}
						errors += prefix + errX.Error()
					}
					Logger.Errorf("Error: During config files initialization -> <%v>...", errors)
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
					worker.Logger.SetVerbosity(clientlog.VerbosityLevelFromString(dc.LogVerbosity))
					Logger.Debugf("Logger Verbosity Setted up to : %v", Logger.GetVerbosity())
				}
				module.RuntimeDeployConfig = dc
				clicommon.DEFAULT_TIMEOUT=time.Duration(dc.ReadTimeout) * time.Second
				errB = boostrap.Load(dc.ConfigDir, dc.EnvSelector, dc.ConfigLang, Logger)
				Logger.Debugf("Errors during config load: %v", len(errB))
				if len(errB) > 0 {
					var errors string = ""
					for _, errX := range errB {
						prefix := ""
						if len(errors) > 0 {
							prefix = "\n"
						}
						errors += prefix + "- " + errX.Error()
					}
					Logger.Errorf("Error: During config files load -> <%v>...", errors)
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

				var pc *module.PluginsConfig = boostrap.GetPluginsType()
				if pc == nil {
					pc = &module.PluginsConfig{}
				}
				pc = boostrap.GetDefaultPluginsType().Merge(pc)
				module.RuntimePluginsType = pc

				Logger.Debugf("Configuration Summary: \nDeploy Config: %v\nDeployType: %v\nNetType: %v\n", dc.String(), dt.String(), nt.String())
				if dt.DeploymentType == module.FILE_SOURCE {
					var filePath string = dc.WorkDir + io.GetPathSeparator() + target
					Logger.Warnf("Loaging Main Feed at path: %s\n", filePath)
					var feed generic.IFeed = generic.NewFeed("default")
					err = feed.Load(filePath)
					if err != nil {
						Logger.Errorf("Error trying to load Feed for file: %s -> Details: \n%s", filePath, err.Error())
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
						panic(fmt.Sprintf("Error trying to validate Feed for file: %s -> Details: \n%s", filePath, errors))
					}
					if len(feedEx.Steps) > 0 {
						Logger.Debugf("Reading file: %s, discovered %s main steps!!", filePath, strconv.Itoa(len(feedEx.Steps)))
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
							panic(fmt.Sprintf("Error: During deploy execution -> <%v>...", errors))
						}
					} else {
						Logger.Warnf("Unable to find any command in the given file: %s", filePath)
						Logger.Warn("Nothing to do here!!")
					}
					Logger.Warn("Deploy procedure complete!!")
				} else {
					Logger.Warnf("Feature %v NOT IMPLEMENTED yet!!", dt.DeploymentType)
				}
			}
		}
	} else {
		help = true
		color.Yellow.Println("Help required")
	}
}
