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
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	// "errors"
	"errpipe/internal/cli"
	"errpipe/internal/utils"
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
 
// var INTRO string = "==============================================\n 	    ERROR PIPE STARTED \n 	 Type 'errpipe --init' to setup application 	 \n============================================== " 
 
func main(){
	utils.EnableANSI()
	
	ok := initFlags()
	
	if ok{
		// Check if config exists
		config, err := cli.LoadConfig()
		if err != nil {
			fmt.Println("Configuration not found.")
			fmt.Println("Please run 'errpipe --init' to setup the application.")
			return
		}

		// fmt.Println(INTRO)
		// fmt.Printf("Using: %s (%s)\n", config.Provider, config.Mode)
		utils.PrintWelcome(config.Provider, config.Mode)
		
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		inputChan := make(chan string)

		go func() {
			for {
				reader := bufio.NewReader(os.Stdin)
				line, err := reader.ReadString('\n')
				if err != nil {
					continue
				}
				inputChan <- strings.TrimSpace(line)
			}
		}()
		
		printPrompt := func() {
			dir, err := os.Getwd()
			if err != nil{
				dir = "UnknownDIR"
			}
			utils.PrintPrompt(dir)
		}

		printPrompt()
		for{
			select {
			case sig := <-sigChan:
				fmt.Printf("\n\t%s[!] Caught signal '%v'. Type 'exit' to quit.%s\n", utils.Fg(196), sig, utils.ResetStr())
				printPrompt()
			case input := <-inputChan:
				if input == "exit"{
					signal.Stop(sigChan)
					return
				} else if input != "" {
					error, ok := runCmd(input)
					if ok{
						utils.SendToAI(error, config)
					}
				}
				printPrompt()
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
			// fmt.Println("Cmd")
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

