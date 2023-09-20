package main

import (
	"fmt"
	"log"
	"os/exec"
	"os/user"

	silver "github.com/kijimad/silver/pkg"
	pipeline "github.com/mattn/go-pipeline"
)

// Dockerfile builderターゲット上で実行する前提
func main() {
	installEmacs()
	checkDotfiles()
	cpSensitiveFile()
	expandInotify()
}

func installEmacs() {
	if silver.IsExistCmd("emacs") {
		fmt.Println("ok")
	} else {
		_, err := exec.Command("apt", "install", "-y", "emacs").CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func checkDotfiles() {
	if silver.IsExistFile("~/dotfiles") {
		fmt.Println("ok")
	} else {
		currentUser, _ := user.Current()
		targetDir := currentUser.HomeDir + "/dotfiles"
		_, err := exec.Command("git", "clone", "https://github.com/kijimaD/dotfiles.git", targetDir).CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
	}
}

// コード管理下にないファイルをコピーする
func cpSensitiveFile() {
	if silver.IsExistFile("~/.authinfo") {
		fmt.Println("ok")
	} else {
		_, err := silver.Copy("~/dotfiles/.authinfo", "~/.authinfo")
		if err != nil {
			log.Fatal(err)
		}
	}
}

// inotifyを増やす
// ホストマシンだけで実行する。コンテナ内かどうかをsudoがあるかないかで判定(微妙...)
// コンテナからは/procに書き込みできないのでエラーになる
func expandInotify() {
	if !silver.IsExistCmd("sudo") {
		fmt.Println("skip")
		return
	}
	_, err := pipeline.Output(
		[]string{"echo", "fs.inotify.max_user_watches=524288"},
		[]string{"sudo", "tee", "-a", "/etc/sysctl.conf"},
	)
	if err != nil {
		log.Fatal(err)
	}
}
