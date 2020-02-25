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
				Logger.Println(fmt.Sprintf("Target: %v", target))
				Logger.Println(fmt.Sprintf("Config: %v", config.Yaml()))
				os.Exit(0)
			}
		}
	} else {
		os.Exit(0)
	}
	os.Exit(1)
}
