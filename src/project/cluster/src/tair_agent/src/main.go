package main

import (
	// "time"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"./logger"
	"./tair"
	"github.com/go-martini/martini"
)

type Config struct {
	LogPath         string   `json:"log_path"`  //各级别日志路径
	ListenPort      int      `json:"listen_port"` //监听端口号
	TairClient      string   `json:"tair_client"`
	TairServer      []string `json:"tair_server"`
}

func loadConfig(path string) *Config {
	if len(path) == 0 {
		panic("path of conifg is null.")
	}

	_, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	var cfg Config
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}

func main() {
	var cfg_path string
	flag.StringVar(&cfg_path, "conf", "../conf/tair-agent-conf.json", "config file path")
	flag.Parse()
	fmt.Println(cfg_path)

	cfg := loadConfig(cfg_path)

	l := logger.GetLogger(cfg.LogPath, "init")
	l.Infof("tair_agent start.")

	l.Infof("tair_agent start.%+v", cfg)

	c := logger.GetLogger(cfg.LogPath, "tair")

	pTair := tair.NewTairClient(cfg.TairServer, cfg.TairClient, c)
	if pTair == nil {
		l.Errorf("NewTairClient fail")
		return
	}

	m := martini.Classic()
	m.Post("/putdata", pTair.SendtoTairPut)
	m.Post("/getdata", pTair.SendtoTairGet)
	
	port := fmt.Sprintf(":%d", cfg.ListenPort)
	l.Infof("listern %+v", port)
	c.Infof("listern %+v", port)
	m.RunOnAddr(port)//改变监听的端口
}
