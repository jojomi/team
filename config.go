package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/juju/errors"
	"github.com/rs/zerolog/log"
	"github.com/shibukawa/configdir"
	"gopkg.in/yaml.v2"
)

func getConfig(configFilename string) (Config, error) {
	var config Config

	fullPath, err := getRelevantConfigPath(configFilename)
	if err != nil {
		return config, err
	}

	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func mustGetConfig(configFilename string) Config {
	result, err := getConfig(configFilename)
	err = errors.Annotate(err, "failed to read config file")
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	return result
}

func getConfigPaths() []string {
	var dirList []string

	if wd, err := os.Getwd(); err == nil {
		dirList = append(dirList, wd)
	}

	configDirs := configdir.New("", strcase.ToKebab(ToolName))
	for _, c := range configDirs.QueryFolders(configdir.Local) {
		dirList = append(dirList, c.Path)
	}
	for _, c := range configDirs.QueryFolders(configdir.Global) {
		dirList = append(dirList, c.Path)
	}
	for _, c := range configDirs.QueryFolders(configdir.System) {
		dirList = append(dirList, c.Path)
	}

	// more relevant first
	return dirList
}

func getRelevantConfigPath(configFilename string) (string, error) {
	var fullPath string
	for _, c := range getConfigPaths() {
		fullPath = filepath.Join(c, configFilename)
		if fileExists(filepath.Join(fullPath)) {
			return fullPath, nil
		}
	}
	return "", errors.New(fmt.Sprintf("Config file %s not found in any of these paths: %s", configFilename, strings.Join(getConfigPaths(), ", ")))
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
