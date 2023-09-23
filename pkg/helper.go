package silver

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

func IsExistCmd(cmdName string) bool {
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

func IsExistFile(path string) bool {
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
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", expandsrc)
	}

	source, err := os.Open(expandsrc)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(expanddst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
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

func displayOutput(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Fprintf(w, "  => %s\n", scanner.Text())
	}
}
