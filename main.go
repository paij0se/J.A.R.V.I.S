package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/adhocore/chin"
	"github.com/paij0se/J.A.R.V.I.S/src/cli"
	"github.com/paij0se/J.A.R.V.I.S/src/tools"
)

var (
	version     = "v0.0.3"
	help        bool
	showVersion bool
)

func init() {
	flag.BoolVar(&help, "help", false, "show help")
	flag.BoolVar(&help, "h", false, "show help")

	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&showVersion, "v", false, "show version")
}
func main() {
	var err error
	var config map[string]string
	if err = cli.CreateConfigDirectory(); err != nil {
		log.Fatal(err)
	}

	if config, err = cli.ReadYml(); err != nil {
		log.Fatal(err)
	}

	if len(config["auth"]) < 51 {
		log.Fatal("Ensure to insert a valid token in cligpt.yml file.")
	}
	if os.Args == nil || len(os.Args) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	flag.Parse()
	if help {
		cli.ShowHelp()
		os.Exit(0)
	}
	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}
	filename := tools.RecordAudio()
	s := chin.New()
	go s.Start()
	texto := tools.SpeechToText(filename, config["language"])
	s.Stop()
	tools.SendTextToOPenAI(texto, config["model"], config["auth"])
}
