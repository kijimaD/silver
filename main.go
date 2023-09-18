package main

import (
	"fmt"
	"os/exec"
	"os/user"

	syscheck "github.com/kijimad/syscheck/pkg"
)

// Dockerfile builderターゲット上で実行する前提
func main() {
	if syscheck.IsExistCmd("ls") {
		fmt.Println("ok")
	} else {
		panic("ls not found")
	}

	installEmacs()

	checkdotfiles()
}

func installEmacs() {
	if syscheck.IsExistCmd("emacs") {
		fmt.Println("ok")
	} else {
		result, err := exec.Command("apt", "install", "-y", "emacs").CombinedOutput()
		if err != nil {
			fmt.Println(string(result))
			panic(err)
		}
	}
}

func checkdotfiles() {
	if syscheck.IsExistFile("~/dotfiles") {
		fmt.Println("ok")
	} else {
		currentUser, _ := user.Current()
		targetDir := currentUser.HomeDir + "/dotfiles"
		result, err := exec.Command("git", "clone", "https://github.com/kijimaD/dotfiles.git", targetDir).CombinedOutput()
		if err != nil {
			fmt.Println(string(result))
			panic(err)
		}
	}
}
