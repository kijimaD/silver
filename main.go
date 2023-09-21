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
	initCrontab()
	initEmacs()
	initDocker()
}

func installEmacs() {
	if silver.IsExistCmd("emacs") {
		fmt.Println("ok, skip")
		return
	}
	_, err := exec.Command("apt", "install", "-y", "emacs").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}

func checkDotfiles() {
	if silver.IsExistFile("~/dotfiles") {
		fmt.Println("ok, skip")
		return
	}
	currentUser, _ := user.Current()
	targetDir := currentUser.HomeDir + "/dotfiles"
	_, err := exec.Command("git", "clone", "https://github.com/kijimaD/dotfiles.git", targetDir).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}

// コード管理下にないファイルをコピーする
func cpSensitiveFile() {
	if silver.IsExistFile("~/.authinfo") {
		fmt.Println("ok, skip")
		return
	}
	_, err := silver.Copy("~/dotfiles/.authinfo", "~/.authinfo")
	if err != nil {
		log.Fatal(err)
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

func initCrontab() {
	if !silver.IsExistCmd("crontab") {
		fmt.Println("skip")
		return
	}
	currentUser, _ := user.Current()
	targetDir := currentUser.HomeDir + "/dotfiles/crontab"
	_, err := exec.Command("crontab", targetDir).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}

func initEmacs() {
	if !silver.IsExistCmd("emacs") || !silver.IsExistFile("~/.emacs.d") {
		fmt.Println("skip")
		return
	}

	currentUser, _ := user.Current()
	targetDir := currentUser.HomeDir + "/.emacs.d/init.el"
	_, err := exec.Command("emacs", "-nw", "--batch", "--load", targetDir, "--eval", `'(all-the-icons-install-fonts t)'`).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}

func initDocker() {
	if !silver.IsExistCmd("docker") || !silver.IsExistCmd("sudo") {
		fmt.Println("skip")
		return
	}

	currentUser, _ := user.Current()
	username := currentUser.Username
	_, err := exec.Command("sudo", "gpasswd", "-a", username, "docker").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	// check
	// TODO: まだ実行結果を保持してないから意味はない
	_, err = exec.Command("id", username).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}
