package cli

import "fmt"

func ShowHelp() {
	fmt.Println(`   
	Usage: jarvis -[OPTION]
	-h ,-help: Display the help command
	-v ,-version: Display the version of J.A.R.V.I.S
   
	   MIT License
	   Made it with ❤️ by paijose
	   https://github.com/paij0se
   
   `)
}
