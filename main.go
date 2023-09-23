package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"

	silver "github.com/kijimad/silver/pkg"
	pipeline "github.com/mattn/go-pipeline"
)

// Dockerfile builderターゲット上で実行する前提
func main() {
	tasks := []silver.Task{
		installEmacs(),
	}
	job := silver.NewJob(tasks)
	job.Run()
	// installBaseTool()
	// getDotfiles()
	// cpSensitiveFile()
	// cpSensitiveFileSSH()
	// expandInotify()
	// initCrontab()
	// initDocker()
	// initGo()
	// runGclone()
	// initStow()
	// installDocker()
	// installChrome()
	// installUnetbootin()
}

func installEmacs() silver.Task {
	t := silver.NewTask("install Emacs")
	t.SetTargetCmd(func() bool { return silver.IsExistCmd("emacs") })
	t.SetDepCmd(func() bool { return silver.IsExistCmd("sudo") })
	t.SetInstCmd(func() error { return t.Exec("emacs") })

	return t
}

func getDotfiles() {
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

func cpSensitiveFileSSH() {
	if silver.IsExistFile("~/.ssh/config") {
		fmt.Println("ok, skip")
		return
	}

	currentUser, _ := user.Current()
	// .sshディレクトリがない場合は作成する
	sshdir := currentUser.HomeDir + "/.ssh/"
	if _, err := os.Stat(sshdir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(sshdir, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	_, err := silver.Copy("~/dotfiles/.ssh/config", "~/.ssh/config")
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

func initGo() {
	if !silver.IsExistCmd("go") {
		fmt.Println("skip")
		return
	}

	repos := []string{
		"github.com/kijimaD/gclone@main",
		"github.com/kijimaD/garbanzo@main",
		"golang.org/x/tools/gopls@latest",
		"github.com/go-delve/delve/cmd/dlv@latest",
		"github.com/nsf/gocode@latest",
		"golang.org/x/tools/cmd/godoc@latest",
		"golang.org/x/tools/cmd/goimports@latest",
		"mvdan.cc/gofumpt@latest",
	}
	for _, repo := range repos {
		_, err := exec.Command("go", "install", repo).CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func runGclone() {
	if !silver.IsExistCmd("gclone") || !silver.IsExistFile("~/dotfiles") {
		fmt.Println("skip")
		return
	}

	if !silver.IsExistFile("~/.ssh/id_rsa") {
		fmt.Println("id_rsa not found, skip")
		return
	}

	currentUser, _ := user.Current()
	configfile := currentUser.HomeDir + "/dotfiles/gclone.yml"
	_, err := exec.Command("gclone", "-f", configfile).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}

func initStow() {
	if !silver.IsExistCmd("stow") {
		fmt.Println("not found stow, skip")
		return
	}

	currentUser, _ := user.Current()
	dotfiles := currentUser.HomeDir + "/dotfiles"
	cmd := exec.Command("stow", ".")
	cmd.Dir = dotfiles
	_, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal(err)
	}
}

func installBaseTool() {
	if !silver.IsExistCmd("sudo") {
		fmt.Println("not found sudo, skip")
		return
	}
	{
		result, err := exec.Command("sudo", "apt-get", "update").CombinedOutput()
		if err != nil {
			fmt.Println(string(result))
			log.Fatal(err)
		}
	}
	{
		_, err := exec.Command("sudo", "apt", "install", "-y", "software-properties-common").CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
	}
	repos := []string{
		"main",
	}
	for _, repo := range repos {
		result, err := exec.Command("sudo", "add-apt-repository", "-y", repo).CombinedOutput()
		if err != nil {
			fmt.Println(string(result))
			log.Fatal(err)
		}
	}
	{
		result, err := exec.Command("sudo", "apt-get", "update").CombinedOutput()
		if err != nil {
			fmt.Println(string(result))
			log.Fatal(err)
		}
	}
	packages := []string{
		"emacs-mozc",
		"cmigemo",
		"fcitx",
		"fcitx-mozc",
		"peco",
		"silversearcher-ag",
		"stow",
		"syncthing",
		"compton",
		"qemu-kvm",
		"libsqlite3-dev", // roam
		"cmake",          // vtermのコンパイル
		"libtool-bin",    // vtermのコンパイル
	}
	for _, p := range packages {
		result, err := exec.Command("sudo", "apt-get", "install", "-y", p).CombinedOutput()
		if err != nil {
			fmt.Println(string(result))
			log.Fatal(err)
		}
	}
}

func installDocker() {
	if silver.IsExistCmd("docker") && silver.IsExistCmd("docker-compose") {
		fmt.Println("ok, skip")
		return
	}
	if !silver.IsExistCmd("sudo") {
		fmt.Println("not found sudo, skip")
		return
	}
	if !silver.IsExistCmd("curl") {
		fmt.Println("not found curl, skip")
		return
	}

	_, err := exec.Command("bash", "-c", "curl -fsSL https://get.docker.com -o get-docker.sh && sudo sh get-docker.sh").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	_, err = exec.Command("bash", "-c", "sudo curl -L https://github.com/docker/compose/releases/download/v2.4.1/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose && sudo chmod +x /usr/local/bin/docker-compose").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}

func installChrome() {
	if !silver.IsExistCmd("sudo") {
		fmt.Println("not found sudo, skip")
		return
	}

	_, err := exec.Command("bash", "-c", "wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && sudo dpkg -i google-chrome-stable_current_amd64.deb").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}

func installUnetbootin() {
	if !silver.IsExistCmd("sudo") {
		fmt.Println("not found sudo, skip")
		return
	}

	_, err := exec.Command("bash", "-c", "wget https://github.com/unetbootin/unetbootin/releases/download/702/unetbootin-linux64-702.bin && sudo mv unetbootin-linux64-702.bin /usr/local/bin/unetbootin && sudo chmod +x /usr/local/bin/unetbootin").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}
