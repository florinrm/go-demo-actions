package main

import (
	"fmt"
)

// PrintMessage dummy function for printing
func PrintMessage(str string) {
	if str == "" {
		return
	}
	fmt.Println(str)
}

func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
	PrintMessage("Hello World")
	PrintMessage("ce frumos e la IDP la 6 seara yay")
}
