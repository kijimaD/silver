package silver

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

var errIsNotRegularFile = errors.New("file is not a regular file")

func IsExistCmd(cmdName string) bool {
	const basecmd = "which"

	cmd := exec.Command(basecmd, cmdName)
	err := cmd.Run()

	// エラーがnilの場合、コマンドは存在する
	return err == nil
}

func IsExistFile(path string) bool {
	// ディレクトリの存在を確認
	path, err := expandTilde(path)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func Copy(src, dst string) (int64, error) {
	expandsrc, err := expandTilde(src)
	if err != nil {
		return 0, err
	}

	expanddst, err := expandTilde(dst)
	if err != nil {
		return 0, err
	}

	sourceFileStat, err := os.Stat(expandsrc)
	if err != nil {
		return 0, fmt.Errorf("file not found: %w", err)
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, errIsNotRegularFile
	}

	source, err := os.Open(expandsrc)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(expanddst)
	if err != nil {
		return 0, fmt.Errorf("failed to create file: %w", err)
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	if err != nil {
		return 0, fmt.Errorf("failed to copy file: %w", err)
	}

	return nBytes, nil
}

func expandTilde(path string) (string, error) {
	// ユーザー情報を取得
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
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

func HomeDir() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return currentUser.HomeDir
}

// コンテナ内で実行されているか判定する。
func OnContainer() bool {
	cmd := exec.Command("systemctl")
	err := cmd.Run()

	return err != nil
}
