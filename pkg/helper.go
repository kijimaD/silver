package syscheck

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

func AssertExistCmd(cmd string) {
	ok := isExistCmd(cmd)
	if ok == false {
		panic(fmt.Sprintf("%s not found", cmd))
	}
}

func isExistCmd(cmdName string) bool {
	const basecmd = "which"
	cmd := exec.Command(basecmd, cmdName)
	err := cmd.Run()

	// エラーがnilの場合、コマンドは存在する
	if err == nil {
		return true
	} else {
		return false
	}
}

func AssertExistFile(path string) {
	ok := isExistFile(path)
	if ok == false {
		panic(fmt.Sprintf("path %s not found", path))
	}
}

func isExistFile(path string) bool {
	// ディレクトリの存在を確認
	path, err := expandTilde(path)
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func expandTilde(path string) (string, error) {
	// ユーザー情報を取得
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	// ホームディレクトリのパス
	homeDir := currentUser.HomeDir

	// チルダを展開
	if len(path) > 0 && path[0] == '~' {
		// チルダをホームディレクトリに置き換え
		path = filepath.Join(homeDir, path[1:])
	}
	return path, nil
}
