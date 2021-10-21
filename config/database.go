package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
	ServerCfg ServerConfig `yaml:"server"`
	DB        DbConfig     `yaml:"Db"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DbConfig struct {
	DriverName   string `yaml:"DriverName"`
	Host         string `yaml:"Host"`
	Port         string `yaml:"Port"`
	MaxIdleConns int    `yaml:"MaxIdleConns"`
	MaxOpenConns int    `yaml:"MaxOpenConns"`
	Timeout      string `yaml:"Timeout"`
	ReadTimeout  string `yaml:"ReadTimeout"`
	WriteTimeout string `yaml:"WriteTimeout"`
	Charset      string `yaml:"Charset"`
	User         string `yaml:"User"`
	Password     string `yaml:"Password"`
	Name         string `yaml:"Name"`
}

func InitConfig() *Config {
	data, err := ioutil.ReadFile("../config.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	conf := &Config{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		log.Fatalln(err)
	}

	dbConf := conf.DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		dbConf.User,
		dbConf.Password,
		dbConf.Host,
		dbConf.Port,
		dbConf.Name,
		dbConf.Charset)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db, err := DB.DB()
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxIdleConns(dbConf.MaxIdleConns)
	db.SetMaxOpenConns(dbConf.MaxOpenConns)

	cfg := conf.ServerCfg
	if cfg.Port == 0 {
		conf.ServerCfg.Port = 8600
	}

	return conf
}
