package main

import (
	"log"
	"testing"

	"github.com/drpaij0se/gufipufi/src/cli"
	"github.com/drpaij0se/gufipufi/src/tools"
)

func TestOpenAI(t *testing.T) {
	var err error
	var config map[string]string
	if err = cli.CreateConfigDirectory(); err != nil {
		log.Fatal(err)
	}

	if config, err = cli.ReadYml(); err != nil {
		log.Fatal(err)
	}

	tools.SendTextToOPenAI("el rap esta en contra del racismo?", "gpt-3.5-turbo", config["auth"])
}
