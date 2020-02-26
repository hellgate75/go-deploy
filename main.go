package main

import (
	"fmt"
	"github.com/hellgate75/go-deploy/cmd"
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
				Logger.Println("Recovering config yaml format ...")
				Logger.Println(fmt.Sprintf("Target: %v", target))
				yaml, errL := config.Yaml()
				if errL == nil {
					Logger.Println(fmt.Sprintf("Config Yaml: %v", yaml))
				} else {
					Logger.Println(fmt.Sprintf("Error retriving Config Yaml: %v", errL))
				}
				Logger.Println("Load config from file ...")
				dc, errC := config.FromYamlFile("./config.yaml")
				if errC == nil {
					Logger.Println(fmt.Sprintf("Config: %v", dc.String()))
					yaml, errL = (*dc).Yaml()
					if errL == nil {
						Logger.Println(fmt.Sprintf("YAML: %v", yaml))
					} else {
						Logger.Println(fmt.Sprintf("Error retriving Config Yaml: %v", errL))
					}

				} else {
					Logger.Println(fmt.Sprintf("Error loading Config: %v", errC))
				}
				os.Exit(0)
			}
		}
	} else {
		os.Exit(0)
	}
	os.Exit(1)
}
