package conf

import (
	"github.com/derekAHua/goLib/base"
	"github.com/derekAHua/goLib/env"
	"github.com/derekAHua/goLib/redis"
	"github.com/derekAHua/goLib/server/http"
	"github.com/derekAHua/goLib/zlog"
	"gorm.io/gorm"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type (
	ConfigConf struct {
		Log    zlog.LogConfig
		Server http.ServerConfig
	}

	ResourceConf struct {
		Redis map[string]redis.Conf
		Mysql map[string]base.MysqlConf
	}
)

func InitResource(conf *ResourceConf) (mMysql map[string]*gorm.DB, mRedis map[string]*redis.Redis) {
	var err error
	mMysql = make(map[string]*gorm.DB, len(conf.Mysql))
	for name, dbConf := range conf.Mysql {
		mMysql[name], err = base.InitMysqlClient(dbConf)
		if err != nil {
			panic("mysql connect error: [%v]" + err.Error())
		}
	}

	mRedis = make(map[string]*redis.Redis, len(conf.Redis))
	for name, redisConf := range conf.Redis {
		mRedis[name], err = redis.InitRedisClient(redisConf)
		if err != nil || mRedis[name] == nil {
			panic("init redis error: [%v]" + err.Error())
		}
	}

	return
}

// InitConf inits config and api and resource and app configs.
func InitConf(config, api, resource, app interface{}) {
	conf := env.DefaultConf
	LoadConf("config.yaml", conf, config)
	LoadConf("api.yaml", conf, api)
	LoadConf("resource.yaml", conf, resource)
	LoadConf("app.yaml", conf, app)
}

// LoadConf load file to struct.
func LoadConf(filename, fileName string, dest interface{}) {
	path := filepath.Join(env.GetConfDirPath(), fileName, filename)

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(filename + " get error: %v " + err.Error())
	}

	err = yaml.Unmarshal(yamlFile, dest)
	if err != nil {
		panic(filename + " unmarshal error: %v" + err.Error())
	}
}
