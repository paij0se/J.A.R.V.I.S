package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/MarkKremer/microphone"
	"github.com/adhocore/chin"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/drpaij0se/gufipufi/src/cli"
	"github.com/gopxl/beep/wav"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

type Data struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func clearTerminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func RecordAudio() string {

	err := microphone.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer microphone.Terminate()

	stream, format, err := microphone.OpenDefaultStream(44100, 2)
	if err != nil {
		log.Fatal(err)
	}
	// Close the stream at the end if it hasn't already been
	// closed explicitly.
	defer stream.Close()

	filename := time.Now().Format("2006-01-02-15-04-05")
	// clear the terminal
	clearTerminal()
	fmt.Println("Recording. Press Ctrl-C to stop.")
	if !strings.HasSuffix(filename, ".wav") {
		filename += ".wav"
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Stop the stream when the user tries to quit the program.
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		stream.Stop()
		stream.Close()
	}()

	stream.Start()
	err = wav.Encode(f, stream, format)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	time.Sleep(1 * time.Second)
	clearTerminal()
	return filename
}

func SpeechToText(filename string) string {
	// print it in purple
	fmt.Println("\033[35m", "[+]", "\033[0m", "Converting to text")
	cmd := "whisper " + filename + " --language Spanish"
	// make a timer to see how long it takes to convert the audio to text
	timer := time.Now()

	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\033[32m", "[+]", "\033[0m", "Conversion finished in: ", time.Since(timer))
	// read the .txt file
	textFile, err := os.Open(strings.TrimSuffix(filename, ".wav") + ".txt")
	if err != nil {
		log.Fatal(err)
	}

	defer textFile.Close()
	scanner := bufio.NewScanner(textFile)
	scanner.Split(bufio.ScanWords)
	var text string
	for scanner.Scan() {
		text += scanner.Text() + " "
	}
	return text

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
	// print in in green
	fmt.Println("\033[32m", r, "\033[0m")
	// play the audio
	if err := PlayAudio(TTS(r)); err != nil {
		log.Fatal(err)
	}
}
func PlayAudio(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()
	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}
func TTS(text string) string {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := polly.New(sess)
	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"), Text: &text, VoiceId: aws.String("Lupe")}
	output, err := svc.SynthesizeSpeech(input)
	if err != nil {
		fmt.Println("Got error calling SynthesizeSpeech:")
		fmt.Print(err.Error())
		os.Exit(1)
	}
	randomFile := time.Now().Format("2006-01-02-15-04-05")
	names := strings.Split(randomFile, ".")
	name := names[0]
	mp3File := name + ".mp3"

	outFile, err := os.Create(mp3File)
	if err != nil {
		fmt.Println("Got error creating " + mp3File + ":")
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, output.AudioStream)
	if err != nil {
		fmt.Println("Got error saving MP3:")
		fmt.Print(err.Error())
		os.Exit(1)
	}
	return mp3File
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
	filename := RecordAudio()
	s := chin.New()
	go s.Start()
	texto := SpeechToText(filename)
	s.Stop()
	SendTextToOPenAI(texto, config["model"], config["auth"])
	cmd := "bash delete.sh"
	_, err = exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatal(err)
	}

}
