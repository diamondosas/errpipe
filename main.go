package main

import (
	"bufio"
	"bytes"
	// "encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	// "errors"
	"errpipe/internal/cli"
)

/* How the Applcication Runs
 * when "errpipe" is run and it makes the active window to copy output adn then prepare them for the AI
 * The user then Run Commands to run his/her applcaition
 * The Output Is then checked wheter it has error or not
 * If it has error it then triggers the AI service corresponding to what is already used in the json in teh spplaciton
 */

/* PROCESS
 * Check whether
 */
 
var INTRO string = "==============================================\n 	    ERROR PIPE STARTED \n 	 Type 'errpipe --init' to setup application 	 \n============================================== " 
 
func main(){
	
	ok := initFlags()
	
	if ok{
		// Check if config exists
		config, err := cli.LoadConfig()
		if err != nil {
			fmt.Println("Configuration not found.")
			fmt.Println("Please run 'errpipe --init' to setup the application.")
			return
		}

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println(INTRO)
		fmt.Printf("Using: %s (%s)\n", config.Provider, config.Mode)
		
		
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
			error, ok := runCmd(input)
			// errmsg := err.Error()
			if ok{
				sendtoAI(error, config)
				
			}
			
		}
	}
	
}

func initFlags() bool{
	help := flag.Bool("help", false, "Print out Help Command")
	init := flag.Bool("init", false ,"Setup the Application")	
	flag.Parse()
	if *help{
		printHelp()
		return false
	}
	if *init{
		cli.InitApp()
		return false
	}
	
	return true
}
func runCmd(input string) (string, bool) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows"{
		// err := exec.Command("powershell", "-command", "$PSVersionTable").Run()
		// if err != nil{ 
			cmd = exec.Command("cmd", "/C", input)
			fmt.Println("Cmd")
		// }else{ 
		// 	cmd = exec.Command("powershell", "-c", input)
		// 	fmt.Println("Power")
		// }
	} else{ //Macos & Linux
		cmd = exec.Command("sh", "-c", input)
	}
	
	var stderrBuf bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	cmd.Stdin = os.Stdin
	
	err := cmd.Run()
	if err != nil{
		return stderrBuf.String(), true
	}
	return "", false //idk the value of err null
}

func printHelp(){
	fmt.Println("Error Pipe Help ")
	fmt.Println(" --help To show the help message ")
	fmt.Println(" --init To initialise or setup the AI")
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

