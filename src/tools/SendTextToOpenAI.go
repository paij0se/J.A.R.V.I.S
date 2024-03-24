package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/paij0se/J.A.R.V.I.S/src/cli"
)

type Data struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func SendTextToOPenAI(text string, model string, auth string) {
	fmt.Println("\033[33m", text, "\033[0m")
	var err error
	client := &http.Client{}
	var data = strings.NewReader(`{
		  "model": "` + model + `",
		  "messages": [{"role": "user", "content": "` + text + `"}],
		  "temperature": 0.7
		}`)
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("The token is valid?", err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var response Data
	json.Unmarshal(bodyText, &response)
	if err != nil {
		log.Println(err)
	}
	r := response.Choices[0].Message.Content
	fmt.Println("\033[32m", r, "\033[0m")
	// play the audio
	var config map[string]string
	if err = cli.CreateConfigDirectory(); err != nil {
		log.Fatal(err)
	}

	if config, err = cli.ReadYml(); err != nil {
		log.Fatal(err)
	}
	if err := PlayAudio("output/" + TTS(r, config["voiceId"])); err != nil {
		log.Fatal(err)
	}
}
