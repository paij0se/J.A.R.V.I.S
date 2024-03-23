package tools

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/MarkKremer/microphone"
	"github.com/gopxl/beep/wav"
)

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
	// move the file to the output folder
	newFilename := "output/" + filename
	if err := os.Rename(filename, newFilename); err != nil {
		log.Fatal(err)
	}
	return filename
}
