package main

import (
	"errpipe/internal/ai/gemini"
	"fmt"
)

func sendtoAI(errormsg string){
	fmt.Println("Sending Error to AI")
	gemini.GeminiCli(errormsg)
}