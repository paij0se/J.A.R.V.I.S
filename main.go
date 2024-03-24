package main

import (
	"log"

	"github.com/adhocore/chin"
	"github.com/paij0se/J.A.R.V.I.S/src/cli"
	"github.com/paij0se/J.A.R.V.I.S/src/tools"
)

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

	filename := tools.RecordAudio()
	s := chin.New()
	go s.Start()
	texto := tools.SpeechToText(filename, config["language"])
	s.Stop()
	tools.SendTextToOPenAI(texto, config["model"], config["auth"])
}
