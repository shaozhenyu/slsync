package server

import (
	"fmt"
	"sync"

	"slsync/share"

	"ember/cli"
	"ember/http/rpc"
)

type Server struct {
	service *share.Service
	locker  sync.Mutex
}

type Client struct {
	Gets func(path string) (leaves share.Leaves, err error) `args:"path" return:"leaves"`
	Test func(path string) (value string, err error)        `args:"path" return:"value,err"`
}

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

func (p *Server) Gets(path string) (leaves share.Leaves, err error) {
	leaves = p.service.GetAll(path)
	fmt.Println("get leaves ok:", leaves)
	return
}

func (p *Server) Test(path string) (string, error) {
	fmt.Println("server test:", path)
	return path + "11", nil
}

func NewServer(path string) (p *Server, err error) {
	service, err := share.NewService(path)
	if err != nil {
		return
	}
	p = &Server{
		service: service,
	}
	return
}
