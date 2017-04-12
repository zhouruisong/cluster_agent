package main

import (
	// "time"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"./fdfsmgr"
	"./logger"
	"github.com/go-martini/martini"
)

type Config struct {
	LogPath         string   `json:"log_path"`  //各级别日志路径
	ListenPort      int      `json:"listen_port"` //监听端口号
	TrackerServer   []string `json:"tracker_server"`
	MinConnection   int      `json:"fdfs_min_connection_count"`
	MaxConnection   int      `json:"fdfs_max_connection_count"`
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
	flag.StringVar(&cfg_path, "conf", "../conf/fdfs-agent-conf.json", "config file path")
	flag.Parse()
	fmt.Println(cfg_path)

	cfg := loadConfig(cfg_path)

	l := logger.GetLogger(cfg.LogPath, "init")
	l.Infof("fdfs_agent start.")

	l.Infof("fdfs_agent start.%+v", cfg)

	f := logger.GetLogger(cfg.LogPath, "fdfs_agent")

	pFdfsmgr := fdfsmgr.NewClient(cfg.TrackerServer, f,
		cfg.MinConnection, cfg.MaxConnection)

	m := martini.Classic()
	m.Post("/fdfsupload", pFdfsmgr.UploadFile)
	m.Post("/fdfsdownload", pFdfsmgr.DownloadFile)
	
	port := fmt.Sprintf(":%d", cfg.ListenPort)
	l.Infof("listern %+v", port)
	f.Infof("listern %+v", port)
	m.RunOnAddr(port)//改变监听的端口
}
