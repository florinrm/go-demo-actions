package main

import (
	"fmt"
)

// PrintMessage dummy function for printing
func PrintMessage(str string) {
	fmt.Println(str)
}

func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
	PrintMessage("Hello World")
}
