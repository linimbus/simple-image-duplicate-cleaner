package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/astaxie/beego/logs"
)

type Config struct {
	SearchDir  string          `json:"SearchDir"`
	SelectList map[string]bool `json:"SelectList"`
	Similarity float64         `json:"Similarity"`
}

var configCache = Config{
	SearchDir: "",
	SelectList: map[string]bool{
		IMG_PNG: true, IMG_JPEG: true, IMG_BMP: true, IMG_HEIC: true,
	},
	Similarity: 90.0,
}

var configFilePath string
var configLock sync.Mutex

func configSyncToFile() error {
	configLock.Lock()
	defer configLock.Unlock()

	value, err := json.MarshalIndent(configCache, "\t", " ")
	if err != nil {
		logs.Error("json marshal config fail, %s", err.Error())
		return err
	}
	return os.WriteFile(configFilePath, value, 0664)
}

func ConfigGet() *Config {
	return &configCache
}

func SearchDirSave(path string) error {
	configCache.SearchDir = path
	return configSyncToFile()
}

func SelectGet(key string) bool {
	flag, b := configCache.SelectList[key]
	if b {
		return flag
	}
	return false
}

func SelectCheck(key string, flag bool) error {
	configCache.SelectList[key] = flag
	return configSyncToFile()
}

func ConfigInit() error {
	configFilePath = fmt.Sprintf("%s%c%s", ConfigDirGet(), os.PathSeparator, "config.json")

	_, err := os.Stat(configFilePath)
	if err != nil {
		err = configSyncToFile()
		if err != nil {
			logs.Error("config sync to file fail, %s", err.Error())
			return err
		}
	}

	value, err := os.ReadFile(configFilePath)
	if err != nil {
		configSyncToFile()

		logs.Error("read config file from app data dir fail, %s", err.Error())
		return err
	}

	err = json.Unmarshal(value, &configCache)
	if err != nil {
		configSyncToFile()

		logs.Error("json unmarshal config fail, %s", err.Error())
		return err
	}

	return nil
}
