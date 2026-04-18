package utils

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

func InInstalled(name string) (string, bool){
	path, err := exec.LookPath(name)
	if err != nil{
		return "", false
	}
	return path, true
}

func IsRunning(name string) (bool){
	proc, err := process.Processes()
	if err != nil{
		fmt.Println(err)
		return false
	}
}