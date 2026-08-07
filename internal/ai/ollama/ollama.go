package ollama

import (
	"os/exec"
)


func doesOllamaExist() bool{
	cmd := exec.Command("ollama --v")
	err := cmd.Run()
	if err != nil{
		return false 
	}
	return true 
}