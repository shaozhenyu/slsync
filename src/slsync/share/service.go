package share

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
)

type Service struct {
	dir  string
	sep  string
	file *os.File
	lock sync.Mutex
}

func (p *Service) dfile(dir string) string {
	return dir + "/0"
}

func (p *Service) load(dir string) (err error) {
	os.MkdirAll(dir, 0777)
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
	p = &Service{sep: "/"}
	if len(dir) != 0 {
		err = p.load(dir)
	}
	return
}
