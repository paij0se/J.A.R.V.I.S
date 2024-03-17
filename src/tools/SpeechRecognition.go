package tools

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func SpeechToText(filename string) string {
	// print it in purple
	fmt.Println("\033[35m", "[+]", "\033[0m", "Converting to text")
	cmd := "whisper " + filename + " --language Spanish --output_dir output"
	// make a timer to see how long it takes to convert the audio to text
	timer := time.Now()

	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\033[32m", "[+]", "\033[0m", "Conversion finished in: ", time.Since(timer))
	// read the .txt file
	textFile, err := os.Open(strings.TrimSuffix("output/"+filename, ".wav") + ".txt")
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
