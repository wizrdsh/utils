package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
  
	"github.com/labstack/gommon/log"
)

type Result struct {
  Message string,
  Data string
  Error error
}

// Generic execution function for 'foreign' host machine execution. 
// - basic string input
// - deferred KILL process on finish (to ensure tasks do not become orphans) 
// - real-time execution output logging (and saving)

func Exec(script string) Result {
	var results Result

	cmd := exec.Command(script)
	defer kill(cmd)

  // Hides external created windows (Windows only? powershell/cmd only?) 
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return Result{
			Message: "[Exec] :: Error on creating Stdout",
			Error:   err,
		}
	}

  // This will hold all of the commands output in real time
	var slurp strings.Builder

	// Real time stdout output 
	scanner := bufio.NewScanner(stdout)
	go func() {
		fmt.Print("[Remote (via Host)]> Beginning Remote Execution \n")

		for scanner.Scan() {
			scanline := scanner.Text()
			fmt.Printf("[Remote]> %s", scanline)
			slurp.WriteString(fmt.Sprintln(scanline))
		}
	}()

	if err := cmd.Start(); err != nil {
		return Result{
			Message: "[Exec] :: Error on Cmd.Start()",
			Error:   err,
		}
	}

	if err_wait := cmd.Wait(); err_wait != nil {
		fmt.Print("[Exec] :: Command has exited :: ", err_wait.Error())
	}

	results = Result{
		Message: "[Exec] :: Success",
		Error:   nil,
		Data:    slurp.String(),
	}

	return results
}

func kill(cmd *exec.Cmd) error {
  // This is command prompt specific, replace with system of choice. 
	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout
	return kill.Run()
}
