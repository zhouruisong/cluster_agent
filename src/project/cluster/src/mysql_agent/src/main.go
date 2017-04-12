package main

import (
	// "time"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"./mysqlmgr"
	"./logger"
	"github.com/go-martini/martini"
)

type Config struct {
	LogPath         string   `json:"log_path"`  //各级别日志路径
	MysqlDsn        string   `json:"mysql_dsn"` //后台存储dsn
	ListenPort      int      `json:"listen_port"` //监听端口号
	ServerId        int      `json:"server_Id"`
	Dbconns         int       `json:"dbconns"`
	Dbidle          int       `json:"dbidle"`
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
	flag.StringVar(&cfg_path, "conf", "../conf/mysql-agent-conf.json", "config file path")
	flag.Parse()
	fmt.Println(cfg_path)

	cfg := loadConfig(cfg_path)
	l := logger.GetLogger(cfg.LogPath, "init")
	l.Infof("myql agent start.")
	l.Infof("myql agent start.%+v", cfg)

	d := logger.GetLogger(cfg.LogPath, "mysql")
	
	pMysqlMgr := mysqlmgr.NewMysqlMgr(cfg.MysqlDsn, cfg.ServerId, cfg.Dbconns, cfg.Dbidle, d)
	if pMysqlMgr == nil {
		l.Errorf("NewMysqlMgr fail")
		return
	}

	m := martini.Classic()
	m.Post("/mysqlinsert", pMysqlMgr.InsertOperate)
	m.Post("/mysqldelete", pMysqlMgr.DeleteOperate)
	
	port := fmt.Sprintf(":%d", cfg.ListenPort)
	l.Infof("listern %+v", port)
	m.RunOnAddr(port)//改变监听的端口
}
