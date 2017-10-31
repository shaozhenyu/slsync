package client

import (
	"fmt"
	"os"

	"slsync/server"

	"ember/cli"
	"ember/http/rpc"
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
}

func CmdSync(args []string) {
	fmt.Println("sync:", args)

	client := &server.Client{}
	rpc := rpc.NewClient("127.0.0.1:8080")
	err := rpc.Reg(client)
	if err != nil {
		return
	}
	client.Test("ttt")

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
}
