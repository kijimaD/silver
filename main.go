package main

import (
	"fmt"
	"os/exec"
	"os/user"

	syscheck "github.com/kijimad/syscheck/pkg"
)

// Dockerfile builderターゲット上で実行する前提
func main() {
	installEmacs()
	checkDotfiles()
	cpSensitiveFile()
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

func checkDotfiles() {
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

func cpSensitiveFile() {
	if syscheck.IsExistFile("~/.authinfo") && syscheck.IsExistFile("~/dotfiles") {
		fmt.Println("ok")
	} else {
		_, err := syscheck.Copy("~/dotfiles/.authinfo", "~/.authinfo")
		if err != nil {
			panic(err)
		}
	}
}
