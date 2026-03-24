package main

import "testing"

func TestPrintMessage(t *testing.T) {
	PrintMessage("Hello World")
}

func TestPrintEmptyMessage(t *testing.T) {
	PrintMessage("")
}
