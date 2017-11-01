package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	Sep                       = string(os.PathSeparator)
	ConfigExt                 = "json"
	LocalConfigFile           = "local." + ConfigExt
	RemoteConfigFile          = "remote." + ConfigExt
	DefaultRootServiceAddress = "127.0.0.1:8080"
)

type Config struct {
	Path   string
	Local  LocalConfig
	Remote RemoteConfig
}

type RemoteConfig struct {
	Address     string
	RootService string
}

type LocalConfig struct {
	Hid  uint64
	Name string
	Path string
}

func InitLocalConfig(path string, slsync string, user string, dpath string) error {
	fmt.Println("init local config")
	if strings.HasPrefix(dpath, "/") || strings.HasPrefix(dpath, "\\") {
		return errors.New("path cannot start with '/' or '\\':" + dpath)
	}

	err := os.MkdirAll(path+Sep+slsync, 0777)
	if err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())
	hid := uint64(rand.Uint32())

	local := LocalConfig{
		Hid:  hid,
		Name: user,
		Path: dpath,
	}

	data, err := json.MarshalIndent(&local, "", "\t")
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path+Sep+slsync+Sep+LocalConfigFile, os.O_RDWR|os.O_CREATE|os.O_EXCL|os.O_TRUNC, 0640)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	_, err = file.WriteString("\n")
	return err
}

func NewConfig(path string, slsync string) (p *Config, err error) {
	p = &Config{Path: path}

	lfile, err := os.Open(path + Sep + slsync + Sep + LocalConfigFile)
	if err != nil {
		return
	}
	defer lfile.Close()

	err = json.NewDecoder(lfile).Decode(&p.Local)
	if err != nil {
		return
	}

	rfile, err := os.Open(path + Sep + slsync + Sep + RemoteConfigFile)
	if err != nil {
		return
	}
	defer rfile.Close()

	err = json.NewDecoder(rfile).Decode(&p.Remote)
	if len(p.Remote.RootService) == 0 {
		p.Remote.RootService = DefaultRootServiceAddress
	}

	if len(p.Remote.Address) == 0 {
		p.Remote.Address = DefaultRootServiceAddress
	}
	return
}

func FindConfig(path string, slsync string) (p *Config, err error) {
	root, err := FindRoot(path, slsync)
	if err != nil {
		return
	}
	return NewConfig(root, slsync)
}

func FindRoot(path string, slsync string) (root string, err error) {
	for {
		var info os.FileInfo
		info, err = os.Stat(path + Sep + slsync)
		if os.IsNotExist(err) || !info.IsDir() {
			if len(path) <= 1 {
				break
			}
			path = filepath.Dir(path)
			continue
		}
		root = path
		return
	}
	err = errors.New("meta not found")
	return
}

func (p *Config) paths(path string) (dpath, abs string, err error) {
	if len(path) == 0 {
		dpath = p.Local.Path
		abs = p.Path
		return
	}

	abs, err = filepath.Abs(path)
	if err != nil {
		return
	}
	if !strings.HasPrefix(abs, p.Path) {
		err = errors.New("path not in managed directory(" + p.Path + "): " + abs)
		return
	}
	if p.Path == abs {
		dpath = p.Local.Path
	} else {
		dpath = p.Local.Path + "/" + abs[len(p.Path)+1:len(abs)]
	}
	return
}
