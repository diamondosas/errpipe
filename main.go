package main

import (
	"bufio"
	// "encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

/* How the Applcication Runs
 * when "errpipe" is run and it makes the active window to copy output adn then prepare them for the AI
 * The user then Run Commands to run his/her applcaition
 * The Output Is then checked wheter it has error or not
 * If it has error it then triggers the AI service corresponding to what is already used in the json in teh spplaciton
 */

/* PROCESS
 * Checkk whether
 */
 
var INTRO string = "==============================================\n 	    ERROR PIPE STARTED \n============================================== " 
 
func main(){
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(INTRO)
	
	
	for{
		dir, err := os.Getwd()
		if err != nil{
			dir = "UnknownDIR"
		}
		fmt.Print("[EP] " + dir + ">")
		
		if !scanner.Scan(){
			break
		}
		input := strings.TrimSpace(scanner.Text())
		
		if input == ""{
			continue
		}else if input == "exit"{
			break
		}
		ok := runCmd(input)
		if ok{
			fmt.Println("Opening AI")
		}
		
	}	
}

func runCmd(input string) bool{
	var cmd *exec.Cmd
	if runtime.GOOS == "windows"{
		_, err := exec.Command("Rename-Item").Output()
		if err != nil{
			cmd = exec.Command("cmd", "/C", input)
			fmt.Println("Cmd")
		}else{
			cmd = exec.Command("powershell", "-c", input)
			fmt.Println("Power!1")
		}
	} else{ //Macos & Linux
		cmd = exec.Command("sh", "-c", input)
	}
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	err := cmd.Run()
	if err != nil{
		return true
	}
	return false
}

// func getErrcode() int{
// 	if runtime.GOOS == "windows"{
// 		log.Println("Windows")
// 		code, err := exec.Command("echo %errorlevel%").Output()
// 		if err != nil{
// 			log.Println(err)
// 		}
// 		return int(binary.BigEndian.Uint64(code))
// 	}else{ // For linux and Macos
// 		code, err := exec.Command("$?").Output()
// 		if err != nil{
// 			log.Println(err)
// 		}
// 		return int(binary.BigEndian.Uint64(code))
// 	}
// }

// func checkCode(code int) bool{
// 	if code == 0{
// 		return false
// 	}else{
// 		return true
// 	}
// }

