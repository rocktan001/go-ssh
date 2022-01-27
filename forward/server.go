package main

import (
	"os"
	"os/exec"
)

func main() {
	for {

		cmd := exec.Command("/media/disk2/go_workspace/src/github.com/rocktan001/www-rocktan001/forward/ssh_forward_v2")
		cmd.Stdout = os.Stdout
		cmd.Run()
		cmd.Wait()
	}
}
