package gemini


import(
	"fmt"
	"os/exec"
)

func geminiCli(error string){
	_, err := exec.LookPath("gemini")
	if err != nil{
		
		fmt.Println("Gemini Cli is not Installed Please use andother AI")
	}
}