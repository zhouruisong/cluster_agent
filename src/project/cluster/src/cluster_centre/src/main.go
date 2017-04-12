package main

import (
	// "time"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"./centre"
	"./logger"
	"github.com/go-martini/martini"
)

type Config struct {
	LogPath        string   `json:"log_path"`  //各级别日志路径
	ListenPort     int      `json:"listen_port"` //监听端口号
	FdfsAgent      []string `json:"fdfs_agent"`
	MysqlAgent     []string `json:"mysql_agent"`
	TairAgent      []string `json:"tair_agent"`
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
	flag.StringVar(&cfg_path, "conf", "../conf/cluster-centre-conf.json", "config file path")
	flag.Parse()
	fmt.Println(cfg_path)

	cfg := loadConfig(cfg_path)

	l := logger.GetLogger(cfg.LogPath, "init")
	l.Infof("cluster centre start.")

	l.Infof("cluster centre start.%+v", cfg)

	d := logger.GetLogger(cfg.LogPath, "centre")

	pCentre := centre.NewClusterMgr(cfg.MysqlAgent, cfg.FdfsAgent, cfg.TairAgent, d)
	if pCentre == nil {
		l.Errorf("NewClusterMgr fail")
		return
	}

	m := martini.Classic()

	m.Post("/putfastdfs", pCentre.FastdfsPutData)
	m.Post("/getfastdfs", pCentre.FastdfsGetData)
	m.Post("/puttotair", pCentre.TairPutData)
	m.Post("/getfromtair", pCentre.TairGetData)
//	m.Post("/deletedata", pCentre.DeleteData)
	
	port := fmt.Sprintf(":%d", cfg.ListenPort)
	l.Infof("listern %+v", port)
	d.Infof("listern %+v", port)
	
	m.RunOnAddr(port)//改变监听的端口
}
