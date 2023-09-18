package main

import (
	syscheck "github.com/kijimad/syscheck/pkg"
)

func main() {
	syscheck.AssertExistCmd("bash")
}
