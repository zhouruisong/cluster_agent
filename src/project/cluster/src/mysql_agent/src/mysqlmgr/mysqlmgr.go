package mysqlmgr

import (
	"fmt"
	"database/sql"
	"time"
	 "net/http"
	 "encoding/json"
	 "io/ioutil"
	 "../../../protocal"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

var (
	host        = "localhost"
	port        = 3306
	username    = "root"
	password    = "123456"
	serverid    = 1
	dbconns     = 200
	dbidleconns = 100
	db          *sql.DB
)

type MysqlMgr struct {
	Logger          *log.Logger
	Dsn             string
}

func NewMysqlMgr(dsn string, id int, dbconn int, dbidleconn int, lg *log.Logger) *MysqlMgr {
	mgr := &MysqlMgr{
		Logger:          lg,
		Dsn:             dsn,
	}

	serverid    = id
	dbconns     = dbconn
	dbidleconns = dbidleconn
	
	err := mgr.init()
	if err != nil {
		mgr.Logger.Infof("mgr.init failed")
		return nil
	}
	mgr.Logger.Infof("NewMysqlMgr ok")
	return mgr
}

func (mgr *MysqlMgr) init() error {
	var err error
	db, err = sql.Open("mysql", mgr.Dsn)
	if err != nil {
		mgr.Logger.Errorf("err:%v.\n", err)
		return err
	}

	db.SetMaxOpenConns(dbconns)
	db.SetMaxIdleConns(dbidleconns)
	db.Ping()
	return nil
}

// 搜索所有live_master库下的表，为全量同步�?
func (mgr *MysqlMgr) SelectTableName() ([]string, error) {
	var tablename []string

	querysql := "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = 'live_master'"
	rows, err := db.Query(querysql)
	if err != nil {
		mgr.Logger.Errorf("err:%v", err)
		return tablename, err
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		mgr.Logger.Errorf("err:%v", err)
		return tablename, err
	}

	var table string
	for rows.Next() {
		err := rows.Scan(&table)
		if err != nil {
			mgr.Logger.Errorf("err:%v", err)
			return tablename, err
		}
		tablename = append(tablename, table)
	}

	mgr.Logger.Infof("tablename:%v, ", tablename)
	return tablename, nil
}

// 查询文件是否已经在数据库中存在，
func (mgr *MysqlMgr) SelectDataExist(taskId, tablename string) int {
	count := 0
	querysql := fmt.Sprintf("select count(1) from live_master.%s where %s.task_id = \"%s\" ",
		tablename, tablename, taskId)

	rows, err := db.Query(querysql)
	if err != nil {
		mgr.Logger.Errorf("err:%v", err)
		return 1
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		mgr.Logger.Errorf("err:%v", err)
		return 1
	}

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			mgr.Logger.Errorf("err:%v", err)
			return 1
		}
	}

	mgr.Logger.Infof("taskId:%v, count: %v", taskId, count)
	return count
}

// 增量同步，启动的时候按照配置文件加载内容，分批发送给备份服务�?
//func (mgr *MysqlMgr) Loadprotocal.StreamInfos(beginIndex int, tablename string) ([]protocal.StreamInfo, int) {
////	mgr.Logger.Infof("start Loadprotocal.StreamInfos")
//	tmpNum := beginIndex * eachNum
//
//	var returnInfo []protocal.StreamInfo
//	querysql := "select id,task_id,task_server,file_name,file_type,file_size,file_md5,domain,app,stream,step," +
//		"publish_time,notify_url,notify_return,status,expiry_time,create_time,update_time,end_time,notify_time from " +
//		tablename + " where file_type=0 " + "limit " + strconv.Itoa(tmpNum) + "," + strconv.Itoa(eachNum)
		
//	mgr.Logger.Infof("querysql: %+v", querysql)
	
	// only loda upload complete
//	rows, err := db.Query(querysql)
//	if err != nil {
//		mgr.Logger.Errorf("err:%v", err)
//		return returnInfo, -1
//	}
//	defer rows.Close()
//
//	err = rows.Err()
//	if err != nil {
//		mgr.Logger.Errorf("err:%v", err)
//		return returnInfo, -1
//	}
//
//	var id uint32
//	var taskId string
//	var taskServer string
//	var fileName string
//	var fileType uint8
//	var fileSize uint32
//	var fileMd5 string
//	var domain string
//	var app string
//	var stream string
//	var step uint8
//	var publishTime uint64
//	var notifyUrl string
//	var notifyReturn string
//	var status uint8
//	var expireTime string
//	var createTime string
//	var updateTime string
//	var endTime string
//	var notifyTime string
//
//	for rows.Next() {
//		err := rows.Scan(&id, &taskId, &taskServer, &fileName, &fileType, &fileSize, &fileMd5,
//			&domain, &app, &stream, &step, &publishTime, &notifyUrl, &notifyReturn, &status,
//			&expireTime, &createTime, &updateTime, &endTime, &notifyTime)
//
//		if err != nil {
//			mgr.Logger.Errorf("err:%v", err)
//			return returnInfo, -1
//		}
//
//		info := protocal.StreamInfo{
//			Id:           id,
//			TaskId:       taskId,
//			TaskServer:   taskServer,
//			FileName:     fileName,
//			FileType:     fileType,
//			FileSize:     fileSize,
//			FileMd5:      fileMd5,
//			Domain:       domain,
//			App:          app,
//			Stream:       stream,
//			Step:         step,
//			PublishTime:  publishTime,
//			NotifyUrl:    notifyUrl,
//			NotifyReturn: notifyReturn,
//			Status:       status,
//			ExpireTime:   expireTime,
//			CreateTime:   createTime,
//			UpdateTime:   updateTime,
//			EndTime:      endTime,
//			NotifyTime:   notifyTime,
//		}
//		returnInfo = append(returnInfo, info)
//	}
//
//	mgr.Logger.Infof("Loadprotocal.StreamInfos len: %+v", len(returnInfo))
//	if len(returnInfo) == 0 {
//		return returnInfo, -1
//	}
//	return returnInfo, 0
//}

func (mgr *MysqlMgr) InsertStreamInfos(i int) int {
	mgr.Logger.Infof("start Insertprotocal.StreamInfos")
	insertsql := "INSERT INTO t_live2odv2_kuwo" +
		"(task_id,task_server,file_name,file_type,file_size,file_md5,domain,app,stream,step," +
		"publish_time,notify_url,notify_return,status,expiry_time,create_time,update_time,end_time,notify_time) " +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
		
	stmtIns, err := db.Prepare(insertsql)
	if err != nil {
		mgr.Logger.Errorf("Prepare failed, err:%v", err)
		return -1
	}
	defer stmtIns.Close()
	
	taskid := fmt.Sprintf("%s%08d", "83f95bd17727e99e286507b2", i)

	mgr.Logger.Infof("taskid: %v", taskid)
	
	_, err = stmtIns.Exec(taskid, "", "/voicelive/219705672_preprocess-1477649359414.m3u8", 0, 0, "", "push.xycdn.kuwo.cn", "voicelive", "219705672_preprocess", 3, 1477649359414, "http://127.0.0.1:8080/accept_test.php", "string(279) \"{\"task_id\":\"8e9addd82febf91d0fffead1760b507a\",\"domain\":\"push.xycdn.kuwo.cn\",\"app\":\"voicelive\",\"stream\":\"219705672_preprocess\",\"tag\":\"/voicelive/219705672_preprocess-1477649359414.m3u8\",\"vod_url\":\"test.com\",\"vod_md5\":\"\",\"vod_size\":\"0\",\"vod_star", 1, "0000-00-00 00:00:00", "2016-10-28 10:09:19", "0000-00-00 00:00:00", "0000-00-00 00:00:00", "2016-10-28 10:09:29")
	if err != nil {
		mgr.Logger.Errorf("insert into mysql failed, err:%v", err)
		return -1
	}

	mgr.Logger.Infof("insert into mysql ok")
	return 0
}

func (mgr *MysqlMgr) InsertMultiStreamInfos(info []protocal.StreamInfo, tablename string) int {
	datalen := len(info)
	if datalen == 0 {
		mgr.Logger.Errorf("datalen = 0")
		return -1
	}

	insertsql := "INSERT INTO " + tablename + " (task_id,task_server,file_name,file_type,file_size,file_md5,domain,app,stream,step," +
		"publish_time,notify_url,notify_return,status,expiry_time,create_time,update_time,end_time,notify_time) " +
		"VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	start := time.Now()
	//Begin函数内部会去获取连接
//	tx, err := db.Begin()
//	if err != nil {
//		mgr.Logger.Errorf("db.Begin(), err:%v", err)
//		return -1
//	}
	
	stmtIns, err := db.Prepare(insertsql)
	if err != nil {
		mgr.Logger.Errorf("Prepare failed, err:%v", err)
		return -1
	}
	defer stmtIns.Close()

	for i := 0; i < datalen; i++ {
		// 源文件默认为0，备份文件该值为1
		info[i].FileType = 1
		//每次循环用的都是tx内部的连接，没有新建连接，效率高
		_, err = stmtIns.Exec(info[i].TaskId, info[i].TaskServer, info[i].FileName, info[i].FileType, info[i].FileSize,
			info[i].FileMd5, info[i].Domain, info[i].App, info[i].Stream, info[i].Step, info[i].PublishTime, info[i].NotifyUrl,
			info[i].NotifyReturn, info[i].Status, info[i].ExpireTime, info[i].CreateTime, info[i].UpdateTime, info[i].EndTime, info[i].NotifyTime)
		
		if err != nil {
			mgr.Logger.Errorf("insert into mysql failed, err:%v", err)
			return -1
		}
	}

//	stmtIns, err := tx.Prepare(insertsql)
//	if err != nil {
//		mgr.Logger.Errorf("Prepare failed, err:%v", err)
//		return -1
//	}
//	defer stmtIns.Close()
//
//	for i := 0; i < datalen; i++ {
//		count := mgr.SelectDataExist(info[i].TaskId, tablename)
//		if count != 0 {
//			mgr.Logger.Errorf("taskid:%v exist in %v", info[i].TaskId, tablename)
//			continue
//		}
		
		// 源文件默认为0，备份文件该值为1
//		info[i].FileType = 1
//		//每次循环用的都是tx内部的连接，没有新建连接，效率高
//		stmtIns.Exec(info[i].TaskId, info[i].TaskServer, info[i].FileName, info[i].FileType, info[i].FileSize,
//			info[i].FileMd5, info[i].Domain, info[i].App, info[i].Stream, info[i].Step, info[i].PublishTime, info[i].NotifyUrl,
//			info[i].NotifyReturn, info[i].Status, info[i].ExpireTime, info[i].CreateTime, info[i].UpdateTime, info[i].EndTime, info[i].NotifyTime)
//	}
//	//出异常回�?
//	defer tx.Rollback()
//
//	//最后释放tx内部的连�?
//	tx.Commit()

	end := time.Now()
	mgr.Logger.Infof("insert ok total time: %v", end.Sub(start).Seconds())

	return 0
}

func (mgr *MysqlMgr) DeleteMultiStreamInfos(info []protocal.StreamInfo, tablename string) int {
//	datalen := len(info)
//	if datalen == 0 {
//		mgr.Logger.Errorf("datalen = 0")
//		return -1
//	}
//
//	insertsql := "INSERT INTO " + tablename + " (task_id,task_server,file_name,file_type,file_size,file_md5,domain,app,stream,step," +
//		"publish_time,notify_url,notify_return,status,expiry_time,create_time,update_time,end_time,notify_time) " +
//		"VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
//
//	start := time.Now()
	//Begin函数内部会去获取连接
//	tx, err := db.Begin()
//	if err != nil {
//		mgr.Logger.Errorf("db.Begin(), err:%v", err)
//		return -1
//	}
	
//	stmtIns, err := db.Prepare(insertsql)
//	if err != nil {
//		mgr.Logger.Errorf("Prepare failed, err:%v", err)
//		return -1
//	}
//	defer stmtIns.Close()
//
//	for i := 0; i < datalen; i++ {
//		// 源文件默认为0，备份文件该值为1
//		info[i].FileType = 1
//		//每次循环用的都是tx内部的连接，没有新建连接，效率高
//		_, err = stmtIns.Exec(info[i].TaskId, info[i].TaskServer, info[i].FileName, info[i].FileType, info[i].FileSize,
//			info[i].FileMd5, info[i].Domain, info[i].App, info[i].Stream, info[i].Step, info[i].PublishTime, info[i].NotifyUrl,
//			info[i].NotifyReturn, info[i].Status, info[i].ExpireTime, info[i].CreateTime, info[i].UpdateTime, info[i].EndTime, info[i].NotifyTime)
//		
//		if err != nil {
//			mgr.Logger.Errorf("insert into mysql failed, err:%v", err)
//			return -1
//		}
//	}

//	stmtIns, err := tx.Prepare(insertsql)
//	if err != nil {
//		mgr.Logger.Errorf("Prepare failed, err:%v", err)
//		return -1
//	}
//	defer stmtIns.Close()
//
//	for i := 0; i < datalen; i++ {
//		count := mgr.SelectDataExist(info[i].TaskId, tablename)
//		if count != 0 {
//			mgr.Logger.Errorf("taskid:%v exist in %v", info[i].TaskId, tablename)
//			continue
//		}
		
		// 源文件默认为0，备份文件该值为1
//		info[i].FileType = 1
//		//每次循环用的都是tx内部的连接，没有新建连接，效率高
//		stmtIns.Exec(info[i].TaskId, info[i].TaskServer, info[i].FileName, info[i].FileType, info[i].FileSize,
//			info[i].FileMd5, info[i].Domain, info[i].App, info[i].Stream, info[i].Step, info[i].PublishTime, info[i].NotifyUrl,
//			info[i].NotifyReturn, info[i].Status, info[i].ExpireTime, info[i].CreateTime, info[i].UpdateTime, info[i].EndTime, info[i].NotifyTime)
//	}
//	//出异常回�?
//	defer tx.Rollback()
//
//	//最后释放tx内部的连�?
//	tx.Commit()

//	end := time.Now()
//	mgr.Logger.Infof("insert ok total time: %v", end.Sub(start).Seconds())

	return 0
}

// 接收DB同步过来的内容，插入对应的live_master表中
func (mgr *MysqlMgr) InsertOperate(res http.ResponseWriter, req *http.Request) {
	buf, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		mgr.Logger.Errorf("ReadAll failed. %v", err)
	}

	rt := mgr.handlerdbinsert(buf)

	r := protocal.MysqlRet{
		Errno:  rt,
		Errmsg: "ok",
	}

	b, err := json.Marshal(&r)
	if err != nil {
		mgr.Logger.Errorf("Marshal failed. %v", err)
	}

	res.Write(b) // HTTP 200
}

// 处理函数
func (mgr *MysqlMgr) handlerdbinsert(buf []byte) int {
	if len(buf) == 0 {
		mgr.Logger.Errorf("buf len = 0")
		return -1
	}

	var q protocal.MysqlInsertBody
	err := json.Unmarshal(buf, &q)
	if err != nil {
		mgr.Logger.Errorf("Error: cannot decode req body %v", err)
		return -1
	}

	mgr.Logger.Infof("handlerdbinfo len: %+v", len(q.Data))

	ret := mgr.InsertMultiStreamInfos(q.Data, q.TableName)
	if ret != 0 {
		mgr.Logger.Errorf("InsertMultiprotocal.StreamInfos failed")
		return -1
	}

	return 0
}

// 接收DB同步过来的内容，插入对应的live_master表中
func (mgr *MysqlMgr) DeleteOperate(res http.ResponseWriter, req *http.Request) {
	buf, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		mgr.Logger.Errorf("ReadAll failed. %v", err)
	}

	rt := mgr.handlerdbdelete(buf)

	r := protocal.MysqlRet{
		Errno:  rt,
		Errmsg: "ok",
	}

	b, err := json.Marshal(&r)
	if err != nil {
		mgr.Logger.Errorf("Marshal failed. %v", err)
	}

	res.Write(b) // HTTP 200
}

// 处理函数
func (mgr *MysqlMgr) handlerdbdelete(buf []byte) int {
	if len(buf) == 0 {
		mgr.Logger.Errorf("buf len = 0")
		return -1
	}

	var q protocal.MysqlDeleteBody
	err := json.Unmarshal(buf, &q)
	if err != nil {
		mgr.Logger.Errorf("Error: cannot decode req body %v", err)
		return -1
	}

	mgr.Logger.Infof("handlerdbinfo len: %+v", len(q.Data))

	ret := mgr.DeleteMultiStreamInfos(q.Data, q.TableName)
	if ret != 0 {
		mgr.Logger.Errorf("InsertMultiprotocal.StreamInfos failed")
		return -1
	}

	return 0
}
