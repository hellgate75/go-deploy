package main

import (
	"os"
	"fmt"
	"log"
	"github.com/hellgate75/go-deploy/cmd"
)

var Logger *log.Logger = log.New(os.Stdout, "[go-deploy] ", log.LstdFlags | log.LUTC)

func init() {
	Logger.Println( "Init ..." )
	
}



func main() {
	defer func() {
		Logger.Println( fmt.Sprint( "Exit ..." ))
		os.Exit(0)
	}()
	if ! cmd.RequiresHelp() {
		Logger.Println( fmt.Sprint("Main ...") )
		config, err := cmd.ParseArguments();
		if err != nil {
			Logger.Println( fmt.Sprintf("Error: %v", err) )
		} else {
			Logger.Println( fmt.Sprintf("Config: %v", *config) )
		}
	}
	os.Exit( 0 )
}
