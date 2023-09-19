package main

import (
	"fmt"
	"os/exec"
	"os/user"

	silver "github.com/kijimad/silver/pkg"
)

// Dockerfile builderターゲット上で実行する前提
func main() {
	installEmacs()
	checkDotfiles()
	cpSensitiveFile()
}

func installEmacs() {
	if silver.IsExistCmd("emacs") {
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
	if silver.IsExistFile("~/dotfiles") {
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
	if silver.IsExistFile("~/.authinfo") && silver.IsExistFile("~/dotfiles") {
		fmt.Println("ok")
	} else {
		_, err := silver.Copy("~/dotfiles/.authinfo", "~/.authinfo")
		if err != nil {
			panic(err)
		}
	}
}
