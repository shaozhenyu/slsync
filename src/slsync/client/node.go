package client

import (
	"errors"
	"fmt"
	"os"

	"slsync/server"
	"slsync/share"

	"ember/http/rpc"
)

const MetaPath = "meta"

type Node struct {
	config  *Config
	service *share.Service
	client  *server.Client
	rpc     *rpc.Client
	cache   *share.Service
}

func (p *Node) Sync(path string) (err error) {
	dpath, abs, err := p.config.paths(path)
	if err != nil {
		return errors.New("node.paths:" + err.Error())
	}
	fmt.Println(dpath, abs)

	leaves, err := p.client.Gets(dpath)
	if err != nil {
		return
	}
	fmt.Println("leaves: ", leaves)

	return
}

func NewNode(path string, slsync string) (p *Node, err error) {
	config, err := FindConfig(path, slsync)
	if err != nil {
		return
	}
	fmt.Println(config)

	meta := config.Path + Sep + slsync + Sep + MetaPath
	err = os.MkdirAll(meta, 0777)
	if err != nil {
		return
	}

	service, err := share.NewService(meta)
	fmt.Println("service:", service)

	client := &server.Client{}
	rpc := rpc.NewClient(config.Remote.Address)
	err = rpc.Reg(client)
	if err != nil {
		return
	}

	cache, err := share.NewService("")
	if err != nil {
		return
	}

	p = &Node{
		config:  config,
		service: service,
		client:  client,
		rpc:     rpc,
		cache:   cache,
	}
	return
}
