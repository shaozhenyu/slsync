package server

import (
	"fmt"
	"sync"

	"ember/cli"
	"ember/http/rpc"
)

func Reg(cmds *cli.Cmds) {
	cmds.Reg("run", "run slsync server", CmdRun)
}

func CmdRun(args []string) {
	dir := "data"
	sobj, err := NewServer(dir)
	cli.Check(err)
	rpc := rpc.NewServer()
	err = rpc.Reg(sobj, &Client{})
	cli.Check(err)
	err = rpc.Run("/", 8080)
	cli.Check(err)
}

type Client struct {
	Test func(path string) (value string, err error) `args:"path" return:"value,err"`
}

func (p *Server) Test(path string) (string, error) {
	fmt.Println("server test:", path)
	return path + "11", nil
}

func NewServer(path string) (p *Server, err error) {
	p = &Server{}
	return
}

type Server struct {
	locker sync.Mutex
}
