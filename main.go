package main

import "github.com/panshul007/apideps/cmd"

var (
	sha1ver   string // sha1 revision used to build the program
	buildTime string // when the executable was built
)

func main() {
	cmd.Execute(sha1ver, buildTime)
}
