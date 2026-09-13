package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/rfancn/airpress/consts"
	"github.com/rfancn/airpress/util"
)

func NewConfig() *Config {
	var configFile string
	flag.StringVar(&configFile, "config", "", "")
	flag.Parse()

	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigType("yaml")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath("./conf/")
		viper.SetConfigName("config")
	}

	viper.SetDefault("airpress.admin_url_path", "admin")

	conf := &Config{}
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(conf); err != nil {
		panic(err)
	}

	if conf.AirPress.WorkDir == "" {
		pwd, err := os.Getwd()
		if err != nil {
			panic(errors.Wrap(err, "init config: get current dir"))
		}
		conf.AirPress.WorkDir, _ = filepath.Abs(pwd)
	} else {
		workDir, err := filepath.Abs(conf.AirPress.WorkDir)
		if err != nil {
			panic(err)
		}
		conf.AirPress.WorkDir = workDir
	}
	normalizeDir := func(path *string, subDir string) {
		if *path == "" {
			*path = filepath.Join(conf.AirPress.WorkDir, subDir)
		} else {
			temp, err := filepath.Abs(*path)
			if err != nil {
				panic(err)
			}
			*path = temp
		}
	}
	normalizeDir(&conf.AirPress.LogDir, "log")
	normalizeDir(&conf.AirPress.TemplateDir, "resources/template")
	normalizeDir(&conf.AirPress.AdminResourcesDir, "resources/admin")
	normalizeDir(&conf.AirPress.UploadDir, consts.AirPressUploadDir)
	normalizeDir(&conf.AirPress.ThemeDir, "resources/template/theme")
	if conf.SQLite3 != nil && conf.SQLite3.Enable {
		normalizeDir(&conf.SQLite3.File, "airpress.db")
	}
	if !util.FileIsExisted(conf.AirPress.TemplateDir) {
		panic("template dir: " + conf.AirPress.TemplateDir + " not exist")
	}
	if !util.FileIsExisted(conf.AirPress.AdminResourcesDir) {
		panic("AdminResourcesDir: " + conf.AirPress.AdminResourcesDir + "not exist")
	}
	if !util.FileIsExisted(conf.AirPress.ThemeDir) {
		panic("theme dir: " + conf.AirPress.ThemeDir + " not exist")
	}

	initDirectory(conf)
	mode = conf.AirPress.Mode
	logMode = conf.AirPress.LogMode
	return conf
}

func initDirectory(conf *Config) {
	mkdirFunc := func(dir string, err error) error {
		if err == nil {
			if _, err = os.Stat(dir); os.IsNotExist(err) {
				err = os.MkdirAll(dir, os.ModePerm)
			}
		}
		return err
	}
	err := mkdirFunc(conf.AirPress.LogDir, nil)
	err = mkdirFunc(conf.AirPress.UploadDir, err)
	if err != nil {
		panic(fmt.Errorf("initDirectory err=%w", err))
	}
}

var (
	mode    string
	logMode LogMode
)

func IsDev() bool {
	return mode == "development"
}

func LogToConsole() bool {
	switch logMode {
	case Console:
		return true
	case File:
		return false
	default:
		return IsDev()
	}
}
