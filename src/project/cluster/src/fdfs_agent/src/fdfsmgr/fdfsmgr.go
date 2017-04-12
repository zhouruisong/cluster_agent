package fdfsmgr

import (
	"../fdfs_client"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"../../../protocal"
	log "github.com/Sirupsen/logrus"
)

type FdfsMgr struct {
	pFdfs   *fdfs_client.FdfsClient
	Logger  *log.Logger
}

func NewClient(trackerlist []string, lg *log.Logger, minConns int, maxConns int) *FdfsMgr {
	pfdfs, err := fdfs_client.NewFdfsClient(trackerlist, lg, minConns, maxConns)
	if err != nil {
		lg.Errorf("NewClient failed")
		return nil
	}
	
	fd := &FdfsMgr{
		pFdfs:   pfdfs,
		Logger:  lg,
	}
	fd.Logger.Infof("NewClient ok")
	return fd
}

// 接收发送的文件消息，存入fastdfs，id写入tair
func (fdfs *FdfsMgr) DownloadFile(res http.ResponseWriter, req *http.Request) {
	buf, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		fdfs.Logger.Errorf("ReadAll failed. %v", err)
	}
	
	rt, content := fdfs.handlerDownloadFile(buf)
	r := protocal.RetCentreDownloadFile{
		Errno:  rt,
		Errmsg: "ok",
		Content: content,
	}

	b, err := json.Marshal(&r)
	if err != nil {
		fdfs.Logger.Errorf("Marshal failed. %v", err)
	}
	
	fdfs.Logger.Infof("UploadFile ok r.Errno: %+v", r.Errno)
	res.Write(b) // HTTP 200
}

func (fdfs *FdfsMgr) handlerDownloadFile(buf []byte) (int, []byte) {
	var ret_buf []byte
	var q protocal.CentreDownloadFileEx
	err := json.Unmarshal(buf, &q)
	if err != nil {
		fdfs.Logger.Errorf("Unmarshal error:%v", err)
		return -1, ret_buf
	}
	
	fdfs.Logger.Infof("before DownloadToBuffer logid:%+v, id:%+v", 
		q.Logid, q.Id)
	
	downloadResponse, err := fdfs.pFdfs.DownloadToBuffer(q.Id, 0, 0)
	if err != nil {
		fdfs.Logger.Errorf("DownloadToBuffer fail, logid:%+v, err:%v, id:%+v", 
			err, q.Logid, q.Id)
		return -1, ret_buf
	}

	if value, ok := downloadResponse.Content.([]byte); ok {
		fdfs.Logger.Infof("DownloadToBuffer ok logid:%+v, id:%+v", 
			q.Logid, q.Id)
		return 0, value
	}

	return -1, ret_buf
}

// 接收发送的文件消息，存入fastdfs，id写入tair
func (fdfs *FdfsMgr) UploadFile(res http.ResponseWriter, req *http.Request) {
	buf, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		fdfs.Logger.Errorf("ReadAll failed. %v", err)
	}
	
	rt, id := fdfs.handlerUploadFile(buf)	
	r := protocal.RetCentreUploadFile{
		Errno:  rt,
		Errmsg: "ok",
		Id: id,
	}

	b, err := json.Marshal(&r)
	if err != nil {
		fdfs.Logger.Errorf("Marshal failed. %v", err)
		return
	}
	
	fdfs.Logger.Infof("UploadFile ok b: %+v", string(b))
	res.Write(b) // HTTP 200
}

// 处理函数
func (fdfs *FdfsMgr) handlerUploadFile(buf []byte) (int, string) {
	if len(buf) == 0 {
		fdfs.Logger.Errorf("handlerUploadFile buf len = 0")
		return -1, ""
	}

	var q protocal.CentreUploadFileEx
	err := json.Unmarshal(buf, &q)
	if err != nil {
		fdfs.Logger.Errorf("Error: cannot decode err:%v",err)
		return -1, ""
	}

	fdfs.Logger.Infof("before UploadAppenderByBuffer logid:%+v, file:%+v", 
		q.Logid, q.Filename)

	uploadres, err := fdfs.pFdfs.UploadAppenderByBuffer(q.Content, "")
	if err != nil {
		fdfs.Logger.Errorf("UploadAppenderByBuffer failed err:%v, logid:%+v, file:%+v", 
			err, q.Logid, q.Filename)
		return -1, ""
	}

	fdfs.Logger.Infof("UploadAppenderByBuffer ok uploadres:%+v, logid:%+v, file:%+v", 
		uploadres, q.Logid, q.Filename)
	return 0, uploadres.RemoteFileId
}
