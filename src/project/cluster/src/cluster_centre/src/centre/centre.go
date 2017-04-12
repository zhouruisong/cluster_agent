package centre

import (
	"fmt"
	"strings"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"../../../protocal"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)
	
type ClusterMgr struct {
	Logger  *log.Logger
	MysqlAgent []string
	FdfsAgent  []string
	TairAgent  []string
}

func NewClusterMgr(mysqlagent []string, fdfsagent []string, tairagent []string, lg *log.Logger) *ClusterMgr {
	cl := &ClusterMgr{
		Logger:  lg,
		MysqlAgent: mysqlagent,
		FdfsAgent: fdfsagent,
		TairAgent: tairagent,
	}
	cl.Logger.Infof("NewClusterMgr ok")
	return cl
}

// 接收发送的文件消息，存入fastdfs，id写入tair
func (cl *ClusterMgr) FastdfsPutData(res http.ResponseWriter, req *http.Request) {
	var rt int
	var id string
	var b []byte
	var err_marshal error
	var ret protocal.RetCentreUploadFile
	logid := fmt.Sprintf("%s", uuid.NewV4())
	
	buf, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
//	defer req.Close()
	if err != nil {
		cl.Logger.Errorf("ReadAll failed. err:%v", err)
		ret.Errno = -1
		ret.Errmsg = "failed"
		goto END
	}
	if len(buf) == 0 {
		cl.Logger.Errorf("buf len = 0")
		ret.Errno = -1
		ret.Errmsg = "failed"
		goto END
	}

	rt, id = cl.handlerUploadData(logid, buf)
	if rt != 0 {
		ret.Errno = rt
		ret.Errmsg = "failed"
	} else {
		ret.Errno = rt
		ret.Errmsg = "ok"
		ret.Id = id
	}
	
	b, err_marshal = json.Marshal(ret)
	if err_marshal != nil {
		cl.Logger.Errorf("Marshal failed. err:%v", err_marshal)
		ret.Errno = -1
		ret.Errmsg = "failed"
		ret.Id = ""
		goto END
	}
	
	cl.Logger.Infof("logid:%+v, FastdfsPutData return ret:%+v", logid, string(b))
END:	
	res.Write(b) // HTTP 200
}

// 处理函数
func (cl *ClusterMgr) handlerUploadData(logid string, buf []byte) (int, string) {	
	var q protocal.CentreUploadFile
	err := json.Unmarshal(buf, &q)
	if err != nil {
		cl.Logger.Errorf("Unmarshal error logid:%+v, err:%v", logid, err)
		return -1, ""
	}
	
	msg := protocal.CentreUploadFileEx {
		Logid: logid,
		Filename: q.Filename,
		Content: q.Content,
	}

	buff, err := json.Marshal(msg)
	if err != nil {
		cl.Logger.Errorf("Marshal failed. err: %v,logid:%+v, file:%+v", 
			err, logid, q.Filename)
		return -1, ""
	}
	
	url := fmt.Sprintf("http://%v/fdfsupload", cl.FdfsAgent[0])
	ip := strings.Split(cl.FdfsAgent[0], ":")
	hosturl := fmt.Sprintf("application/json;charset=utf-8;hostname:%v", ip[0])	
	
	body := bytes.NewBuffer([]byte(buff))
	res, err := http.Post(url, hosturl, body)
	if err != nil {
		cl.Logger.Errorf("http post return failed err:%v, logid:%+v, file:%+v", 
			err, logid, q.Filename)
		return -1, ""
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		cl.Logger.Errorf("ioutil readall failed, logid:%+v, err:%v, file:%+v", 
			logid, err, q.Filename)
		return -1, ""
	}
	
	var ret protocal.RetCentreUploadFile
	err = json.Unmarshal(result, &ret)
	if err != nil || ret.Errno != 0 {
		cl.Logger.Errorf("Unmarshal return body error, logid:%+v, err:%v, file:%+v", 
			logid, err, q.Filename)
		return -1, ""
	}
	
	cl.Logger.Infof("logid:%+v, handlerUploadData return ret:%+v, file:%+v", 
		logid, ret, q.Filename)
	return 0, ret.Id
}

// 接收发送的文件消息，存入fastdfs，id写入tair
func (cl *ClusterMgr) FastdfsGetData(res http.ResponseWriter, req *http.Request) {
	var rt int
	var content []byte
	var b []byte
	var err_marshal error
	var ret protocal.RetCentreDownloadFile
	logid := fmt.Sprintf("%s", uuid.NewV4())
	
	cl.Logger.Infof("logid: %+v", logid)
	
	buf, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()	
	if err != nil {
		cl.Logger.Errorf("ReadAll failed. %v", err)
		ret.Errno = -1
		ret.Errmsg = "failed"
		goto END
	}
	if len(buf) == 0 {
		cl.Logger.Errorf("buf len = 0")
		ret.Errno = -1
		ret.Errmsg = "failed"
		goto END
	}

	rt, content = cl.handlerDownloadData(logid, buf)
	if rt != 0 {
		ret.Errno = rt
		ret.Errmsg = "failed"
	} else {
		ret.Errno = rt
		ret.Errmsg = "ok"
		ret.Content = content
	}
	
	b, err_marshal = json.Marshal(ret)
	if err_marshal != nil {
		cl.Logger.Errorf("Marshal failed. %v", err_marshal)
		ret.Errno = -1
		ret.Errmsg = "failed"
		goto END
	}
	
	cl.Logger.Infof("logid:%+v, FastdfsGetData return ret: %+v", logid, string(b))
END:	
	res.Write(b) // HTTP 200
}

// 处理函数
func (cl *ClusterMgr) handlerDownloadData(logid string, buf []byte) (int, []byte) {	
	var ret_buf []byte
	var q protocal.CentreDownloadFile
	err := json.Unmarshal(buf, &q)
	if err != nil {
		cl.Logger.Errorf("Unmarshal error logid: %+v, err:%v", logid, err)
		return -1, ret_buf
	}
	
	msg := protocal.CentreDownloadFileEx {
		Logid: logid,
		Id: q.Id,
	}

	buff, err := json.Marshal(msg)
	if err != nil {
		cl.Logger.Errorf("Marshal failed. err: %v, logid: %+v, id: %+v", 
			err, logid, q.Id)
		return -1, ret_buf
	}
	
	url := fmt.Sprintf("http://%v/fdfsdownload", cl.FdfsAgent[0])
	ip := strings.Split(cl.FdfsAgent[0], ":")
	hosturl := fmt.Sprintf("application/json;charset=utf-8;hostname:%v", ip[0])	
	
	body := bytes.NewBuffer([]byte(buff))
	res, err := http.Post(url, hosturl, body)
	if err != nil {
		cl.Logger.Errorf("http post return failed err: %v, logid: %+v, id: %+v", 
			err, logid, q.Id)
		return -1, ret_buf
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		cl.Logger.Errorf("ioutil readall failed, logid: %+v, err:%v, id: %+v", 
			logid, err, q.Id)
		return -1, ret_buf
	}
	
	var ret protocal.RetCentreDownloadFile
	err = json.Unmarshal(result, &ret)
	if err != nil || ret.Errno != 0 {
		cl.Logger.Errorf("Unmarshal return error,logid:%+v,err:%v,id:%+v,Errno:%+v", 
			logid, err, q.Id, ret.Errno)
		return -1, ret_buf
	}
	
	if len(ret.Content) == 0 {
		cl.Logger.Errorf("len = 0,logid:%+v,err:%v,id:%+v,Errno:%+v,len:%+v", 
			logid, err, q.Id, ret.Errno, len(ret.Content))
		return -1, ret_buf
	}
	
	cl.Logger.Infof("logid: %+v, handlerUploadData ok id: %+v, Errno: %+v", 
		logid, q.Id, ret.Errno)
	return 0, ret.Content
}

// 接收发送的id写入tair
func (cl *ClusterMgr) TairPutData(res http.ResponseWriter, req *http.Request) {
	var rt int
	var b []byte
	var err_marshal error
	var ret protocal.RetTairPut
	logid := fmt.Sprintf("%s", uuid.NewV4())
	
	buf, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()	
	if err != nil {
		cl.Logger.Errorf("ReadAll failed. %v", err)
		ret.Errno = -1
		ret.Errmsg = "failed"
		goto END
	}
	if len(buf) == 0 {
		cl.Logger.Errorf("buf len = 0")
		ret.Errno = -1
		ret.Errmsg = "failed"
		goto END
	}

	rt = cl.handlerSendToTairPut(logid, buf)
	if rt != 0 {
		ret.Errno = rt
		ret.Errmsg = "failed"
	} else {
		ret.Errno = rt
		ret.Errmsg = "ok"
	}
	
	b, err_marshal = json.Marshal(ret)
	if err_marshal != nil {
		cl.Logger.Errorf("Marshal failed. %v", err_marshal)
		return
	}
	
	cl.Logger.Infof("logid:%+v, TairPutData return ret:%+v", logid, string(b))
END:	
	res.Write(b) // HTTP 200
}

// 处理函数
func (cl *ClusterMgr) handlerSendToTairPut(logid string, buf []byte) int {
	var q protocal.SednTairPutBody
	err := json.Unmarshal(buf, &q)
	if err != nil {
		cl.Logger.Errorf("Unmarshal error logid: %+v, err:%v", logid, err)
		return -1
	}
	
	msg := protocal.SednTairPutBodyEX {
		Logid: logid,
		Message: q,
	}

	buff, err := json.Marshal(msg)
	if err != nil {
		cl.Logger.Errorf("Marshal failed. err:%v, logid:%+v, q:%+v", 
			err, logid, q)
		return -1
	}
	
	url := fmt.Sprintf("http://%v/putdata", cl.TairAgent[0])
	ip := strings.Split(cl.TairAgent[0], ":")
	hosturl := fmt.Sprintf("application/json;charset=utf-8;hostname:%v", ip[0])	
	
	body := bytes.NewBuffer([]byte(buff))
	res, err := http.Post(url, hosturl, body)
	if err != nil {
		cl.Logger.Errorf("http post return failed err:%v, logid:%+v, q:%+v", 
			err, logid, q)
		return -1
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		cl.Logger.Errorf("ioutil readall failed, logid:%+v, err:%v, q:%+v", 
			logid, err, q)
		return -1
	}
	
	var ret protocal.RetTairPut
	err = json.Unmarshal(result, &ret)
	if err != nil {
		cl.Logger.Errorf("Unmarshal return body error, logid:%+v, err:%v, q:%+v", 
			logid, err, q)
		return -1
	}

	cl.Logger.Infof("logid:%+v, handlerSendToTairPut return ret:%+v", logid, ret)
	return ret.Errno
}

// 接收发送的id写入tair
func (cl *ClusterMgr) TairGetData(res http.ResponseWriter, req *http.Request) {
	var b []byte
	var err_marshal error
	var ret protocal.RetTairGet
	logid := fmt.Sprintf("%s", uuid.NewV4())
	
	buf, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()	
	if err != nil {
		cl.Logger.Errorf("ReadAll failed. %v", err)
		ret.Errno = -1
		ret.Errmsg = "failed"
		goto END
	}
	if len(buf) == 0 {
		cl.Logger.Errorf("buf len = 0")
		ret.Errno = -1
		ret.Errmsg = "failed"
		goto END
	}

	cl.handlerSendToTairGet(logid, buf, &ret)
	
	b, err_marshal = json.Marshal(ret)
	if err_marshal != nil {
		cl.Logger.Errorf("Marshal failed. %v", err_marshal)
		return
	}
	
	cl.Logger.Infof("logid:%+v, TairGetData return  ret:%+v", logid, string(b))
END:	
	res.Write(b) // HTTP 200
}

// 处理函数
func (cl *ClusterMgr) handlerSendToTairGet(logid string, buf []byte, ret *protocal.RetTairGet) {
	var q protocal.SednTairGetBody
	err := json.Unmarshal(buf, &q)
	if err != nil {
		cl.Logger.Errorf("Unmarshal error logid: %+v, err:%v", logid, err)
		return
	}
	
	msg := protocal.SednTairGetBodyEX {
		Logid: logid,
		Message: q,
	}

	buff, err := json.Marshal(msg)
	if err != nil {
		cl.Logger.Errorf("Marshal failed. err:%v, logid:%+v, q:%+v", 
			err, logid, q)
		return
	}
	
	url := fmt.Sprintf("http://%v/getdata", cl.TairAgent[0])
	ip := strings.Split(cl.TairAgent[0], ":")
	hosturl := fmt.Sprintf("application/json;charset=utf-8;hostname:%v", ip[0])	
	
	body := bytes.NewBuffer([]byte(buff))
	res, err := http.Post(url, hosturl, body)
	if err != nil {
		cl.Logger.Errorf("http post return failed err:%v, logid:%+v, q:%+v", 
			err, logid, q)
		return
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		cl.Logger.Errorf("ioutil readall failed, logid:%+v, err:%v, q:%+v", 
			logid, err, q)
		return
	}
	
	err = json.Unmarshal(result, ret)
	if err != nil {
		cl.Logger.Errorf("Unmarshal return body error, logid:%+v, err:%v, q:%+v", 
			logid, err, q)
		return
	}

	cl.Logger.Infof("handlerSendToTairGet return logid:%+v, ret:%+v", logid, ret)
	return
}