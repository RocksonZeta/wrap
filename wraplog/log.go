package wraplog

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/RocksonZeta/logger"
	"github.com/RocksonZeta/wrap/utils/fileutil"
	"gopkg.in/yaml.v2"
)

var configFile = "config.yml"
var Logger logger.FileLogger
var config = struct {
	Log struct {
		Wrap logger.Options
	}
}{}

func init() {
	cwd := fileutil.FindFileDir(configFile)
	bs, err := ioutil.ReadFile(filepath.Join(cwd, configFile))
	if err != nil {
		fmt.Println("read config.yml err. " + err.Error())
		return
	} else {
		err = yaml.Unmarshal(bs, &config)
		if err != nil {
			fmt.Println("Unmarshal config.yml err. " + err.Error())
			return
		}
	}
	Logger = logger.NewLogger(config.Log.Wrap)
}

// func SetWrapLoggerOptions(options logger.Options) {
// 	Logger = logger.NewLogger(options)
// 	fmt.Println("set logger")
// }

// func init() {
// 	var config = struct {
// 		Log struct {
// 			Wrap logger.Options
// 		}
// 	}{}
// 	config.Log.Wrap = logger.Options{Console: true, Level: "trace", File: "logs/wrap.%Y%m%d.log", ForceNewFile: false, ShowLocalIp: true}
// 	// configFile := "config.yml"
// 	cwd := fileutil.FindFileDir(ConfigFile)
// 	fmt.Println("find config.yml in " + cwd)
// 	bs, err := ioutil.ReadFile(filepath.Join(cwd, ConfigFile))
// 	if err != nil {
// 		fmt.Println("read config.yml err. " + err.Error())
// 		// panic(err)
// 	} else {
// 		err = yaml.Unmarshal(bs, &config)
// 		if err != nil {
// 			fmt.Println("Unmarshal config.yml err. " + err.Error())
// 			panic(err)
// 		}
// 	}
// 	config.Log.Wrap.MaxAge = config.Log.Wrap.MaxAge * 24 * time.Hour
// 	config.Log.Wrap.File = filepath.Join(cwd, config.Log.Wrap.File)
// 	Logger = logger.NewLogger(config.Log.Wrap)
// }
