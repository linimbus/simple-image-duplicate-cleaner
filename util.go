package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/astaxie/beego/logs"
)

func VersionGet() string {
	return "v0.1.0"
}

func SaveToFile(name string, body []byte) error {
	return os.WriteFile(name, body, 0664)
}

func CapSignal(proc func()) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		proc()
		logs.Error("recv signcal %s, ready to exit", sig.String())
		os.Exit(-1)
	}()
}

func ReadFileList(dir string) ([]string, error) {
	output := make([]string, 0)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			output = append(output, absPath)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}
