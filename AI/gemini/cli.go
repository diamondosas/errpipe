package gemini


import(
	"fmt"
	"os/exec"
	"errpipe/ai/utils"
)

func geminiCli(error string){
	_, err := exec.LookPath("gemini")
	if err != nil{	
		fmt.Println("Gemini Cli is not Installed Please use andother AI")
	}
	//Check whether gemini CLI is open or not
	ok, err := utils.isRunning("gemini")
	if err != nil{
		fmt.Println(err)
	}	
	// If it is not open we start a terminal with the paremerter "geminicli --prompt(error) 
	
	// If it is open we		
	// Open the terminal to the main screen 
	// Then use robot-go to type the error into the screen and then that is all 
}