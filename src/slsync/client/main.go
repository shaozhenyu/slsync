package client

import (
	"fmt"
	"os"

	"ember/cli"
)

const (
	SlSyncDir = ".slsync"
)

func Reg(cmds *cli.Cmds) {
	cmds.Reg("init", "init current dir", CmdInit)
	cmds.Reg("sync", "sync file or dir", CmdSync)
}

func CmdInit(args []string) {
	if len(args) < 1 {
		fmt.Println("usage: <bin> init path-of-the-whole-data-tree [user-name]")
		os.Exit(1)
	}
	dpath := args[0]
	user := "nobody"
	if len(args) > 1 {
		user = args[1]
	}
	pwd, err := os.Getwd()
	cli.Check(err)
	err = InitLocalConfig(pwd, SlSyncDir, user, dpath)
	cli.Check(err)
	//TODO init remote config
}

func CmdSync(args []string) {
	fmt.Println("sync:", args)

	Sync(args)
}

func Sync(args []string) {
	if len(args) > 1 {
		fmt.Println("usage: <bin> sync file-path")
		os.Exit(1)
	}

	pwd, err := os.Getwd()
	cli.Check(err)
	root, err := FindRoot(pwd, SlSyncDir)
	cli.Check(err)

	path := pwd
	if len(args) == 1 {
		path = args[0]
	}
	fmt.Println(path, root)

	node, err := NewNode(pwd, SlSyncDir)
	cli.Check(err)
	fmt.Println("node:", node)

	node.Sync(path)
}
