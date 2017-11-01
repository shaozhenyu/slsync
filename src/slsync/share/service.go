package share

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type Service struct {
	dir  string
	sep  string
	file *os.File
	tree *Tree
	lock sync.Mutex
}

func (p *Service) GetAll(path string) (leaves Leaves) {
	fmt.Println("getall path:", path)
	p.lock.Lock()
	defer p.lock.Unlock()
	tree := p.get(path)
	return tree.GetLeaves(path)
}

func (p *Service) get(path string) (tree *Tree) {
	if len(path) == 0 {
		return p.tree
	}

	if strings.HasPrefix(path, p.sep) {
		path = path[len(p.sep):]
	}
	field := strings.Split(path, p.sep)
	tree = p.tree
	for _, name := range field {
		fmt.Println("get file name:", name)
		tree = tree.Get(name)
		fmt.Println("get file tree:", tree)
	}
	return
}

func (p *Service) dfile(dir string) string {
	return dir + "/0"
}

func (p *Service) load(dir string) (err error) {
	os.MkdirAll(dir, 0777)
	fmt.Println("service dir:", dir)
	p.file, err = os.OpenFile(p.dfile(dir), os.O_RDWR|os.O_CREATE|os.O_SYNC, 0640)
	if err != nil {
		return
	}

	r := bufio.NewReader(p.file)
	for {
		data, prefix, err := r.ReadLine()
		if err != nil {
			if err != io.EOF {
				return err
			} else {
				return nil
			}
		}
		fmt.Println("file: ", data, prefix)
	}

}

func NewService(dir string) (p *Service, err error) {
	p = &Service{sep: "/", tree: NewTree()}
	if len(dir) != 0 {
		err = p.load(dir)
	}
	return
}
