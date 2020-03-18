package wraplog

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/RocksonZeta/logger"
	"github.com/RocksonZeta/wrap/utils/fileutil"
	"gopkg.in/yaml.v2"
)

var Logger logger.FileLogger

func init() {
	var config = struct {
		Log struct {
			Wrap logger.Options
		}
	}{}
	config.Log.Wrap = logger.Options{Console: true, Level: "trace", File: "logs/wrap.%Y%m%d.log", ForceNewFile: false, ShowLocalIp: true}
	configFile := "config.yml"
	cwd := fileutil.FindFileDir(configFile)
	fmt.Println("find config.yml in " + cwd)
	bs, err := ioutil.ReadFile(filepath.Join(cwd, configFile))
	if err != nil {
		fmt.Println("read config.yml err. " + err.Error())
		panic(err)
	}
	err = yaml.Unmarshal(bs, &config)
	if err != nil {
		fmt.Println("Unmarshal config.yml err. " + err.Error())
		panic(err)
	}
	config.Log.Wrap.MaxAge = config.Log.Wrap.MaxAge * 24 * time.Hour
	config.Log.Wrap.File = filepath.Join(cwd, config.Log.Wrap.File)
	Logger = logger.NewLogger(config.Log.Wrap)
}
