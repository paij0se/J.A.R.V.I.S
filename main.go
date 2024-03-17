package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/adhocore/chin"
	"github.com/drpaij0se/gufipufi/src/cli"
	"github.com/drpaij0se/gufipufi/src/tools"
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
	texto := tools.SpeechToText(filename)
	s.Stop()
	tools.SendTextToOPenAI(texto, config["model"], config["auth"])
	filesToDelete := []string{"*.mp3", "*.wav"}
	for _, file := range filesToDelete {
		cmd := exec.Command("sh", "-c", "rm "+file)
		output, err := cmd.Output()

		if err != nil {
			fmt.Println("Error al ejecutar el comando:", err)
			return
		}
		fmt.Println(string(output))
	}

}
