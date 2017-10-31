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
	Sep             = string(os.PathSeparator)
	ConfigExt       = "json"
	LocalConfigFile = "local." + ConfigExt
)

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
