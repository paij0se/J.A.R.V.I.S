package tools

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
)

func checkAWSIsInstalled() {
	_, err := exec.LookPath("aws")
	if err != nil {
		fmt.Println("AWS CLI is not installed. Please install it and try again.")
		os.Exit(1)
	}
}

func TTS(text string, voice string) string {
	checkAWSIsInstalled()
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// check the text if its a code block: ```python`
	if strings.Contains(text, "```") {
		text = ""
	}
	svc := polly.New(sess)
	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"), Text: &text, VoiceId: aws.String(voice)}
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
