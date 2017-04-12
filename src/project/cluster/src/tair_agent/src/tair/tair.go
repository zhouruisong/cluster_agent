package tair

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"bytes"
	"strings"
	"encoding/json"
	"../../../protocal"
	log "github.com/Sirupsen/logrus"
)

type TairClient struct {
	Logger     *log.Logger
	TairServer string
	Tairclient string
}

func NewTairClient(server []string, tairclient string, lg *log.Logger) *TairClient {
	var sever_addr string
	if len(server) == 2 {
		sever_addr = server[0] + "," + server[1]
	} else if len(server) == 1 {
		sever_addr = server[0]
	} else {
		fmt.Println("ERROR: tair_server len: %d", len(server))
		return nil
	}

	c := &TairClient{
		Logger:     lg,
		TairServer: sever_addr,
		Tairclient: tairclient,
	}
	c.Logger.Infof("NewTairClient ok")
	return c
}

// 向tair上传
func (tair *TairClient) SendtoTairPut(res http.ResponseWriter, req *http.Request) {
	buf, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		tair.Logger.Errorf("ReadAll failed. %v", err)
	}
	
	rt, status := tair.handlerSendtoTairPut(buf)
	r := protocal.RetTairPut{
		Errno:  rt,
		Errmsg: status,
	}

	b, err := json.Marshal(&r)
	if err != nil {
		tair.Logger.Errorf("Marshal failed. %v", err)
	}
	
	tair.Logger.Infof("SendtoTairPut return: %+v", r)
	res.Write(b) // HTTP 200
}

func (tair *TairClient) handlerSendtoTairPut(buf []byte) (int, string) {
	var q protocal.SednTairPutBodyEX
	err := json.Unmarshal(buf, &q)
	if err != nil {
		tair.Logger.Errorf("Unmarshal error:%v", err)
		return -1, ""
	}

	tair.Logger.Infof("q: %+v", q)
		
	msg := protocal.SendTairMesage {
		Command: "pput",
		ServerAddr: tair.TairServer,
		GroupName: "group_1",
		Keys: q.Message.Keys,
	}
	
	buff, err := json.Marshal(msg)
	if err != nil {
		tair.Logger.Errorf("Marshal failed.logid:%+v, err:%v, msg:%+v", 
			q.Logid, err, msg)
		return -1, ""
	}


	url := fmt.Sprintf("http://%v/tair", tair.Tairclient)
	ip := strings.Split(tair.Tairclient, ":")
	hosturl := fmt.Sprintf("application/json;charset=utf-8;hostname:%v", ip[0])
	
	body := bytes.NewBuffer([]byte(buff))
	res, err := http.Post(url, hosturl, body)
	if err != nil {
		tair.Logger.Errorf("http post return failed.logid:%+v, err:%v , buff:%+v", 
			q.Logid, err, string(buff))
		return -1, ""
	}
		
	defer res.Body.Close()

	if res.StatusCode == 200 {
		tair.Logger.Infof("post return ok logid:%+v, code:%+v, status:%+v", 
			q.Logid, res.StatusCode, res.Status)
		return 0, res.Status
	} 
	
	tair.Logger.Infof("post return failed logid:%+v, res:%+v", q.Logid, res)
	return -1, ""
}

// 向tair上传
func (tair *TairClient) SendtoTairGet(res http.ResponseWriter, req *http.Request) {	
	buf, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		tair.Logger.Errorf("ReadAll failed. %v", err)
	}
	
	rt, key := tair.handlerSendtoTairGet(buf)
	var r protocal.RetTairGet
	if rt == 0 {
		r.Errno  = rt
		r.Errmsg = "ok"
		r.Keys = key
	} else {
		r.Errno  = rt
		r.Errmsg = "failed"
		r.Keys = key
	}

	b, err := json.Marshal(&r)
	if err != nil {
		tair.Logger.Errorf("Marshal failed. %v", err)
	}
	
	tair.Logger.Infof("SendtoTairGet return: %+v", r)
	res.Write(b) // HTTP 200
}

// 向tair获取value
func (tair *TairClient) handlerSendtoTairGet(buf []byte) (int, []protocal.RetTairGetDetail) {
	var ret_buff []protocal.RetTairGetDetail
	var q protocal.SednTairGetBodyEX
	err := json.Unmarshal(buf, &q)
	if err != nil {
		tair.Logger.Errorf("Unmarshal error:%v", err)
		return -1, ret_buff
	}

	tair.Logger.Infof("q: %+v", q)
	
	msg := protocal.SendTairMesageGet{
		Command:    "pget",
		ServerAddr: tair.TairServer,
		GroupName:  "group_1",
		Keys: q.Message.Keys,
	}

	tair.Logger.Infof("ioutil readall failed, logid:%+v, msg:%+v", q.Logid, msg)
		
	buff, err := json.Marshal(msg)
	if err != nil {
		tair.Logger.Errorf("Marshal failed.logid:%+v, err:%v, msg:%+v", 
			q.Logid, err, msg)
		return -1, ret_buff
	}
	
	url := fmt.Sprintf("http://%v/tair", tair.Tairclient)
	ip := strings.Split(tair.Tairclient, ":")
	hosturl := fmt.Sprintf("application/json;charset=utf-8;hostname:%v", ip[0])
	
	body := bytes.NewBuffer([]byte(buff))
	res, err := http.Post(url, hosturl, body)
	if err != nil {
		tair.Logger.Errorf("http post return failed.logid:%+v, err:%v , buff:%+v", 
			q.Logid, err, string(buff))
		return -1, ret_buff
	}

	tair.Logger.Infof("post return failed logid:%+v, res.Body:%+v", q.Logid, res.Body)

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		tair.Logger.Errorf("ioutil readall failed, logid:%+v, err:%v, buff:%+v", 
			q.Logid, err, string(buff))
		return -1, ret_buff
	}
	
	var RetKeys protocal.RetTairGetKeys
	err = json.Unmarshal(result, &RetKeys)
	if err != nil {
		tair.Logger.Errorf("Unmarshal return body error, logid:%+v, err:%v, buff:%+v", 
			q.Logid, err, string(buff))
		return -1, ret_buff
	}
	
	if res.StatusCode == 200 {
		tair.Logger.Infof("post return ok logid:%+v, code:%+v, status:%+v", 
			q.Logid, res.StatusCode, res.Status)
		return 0, RetKeys.Keys
	}
	
	tair.Logger.Infof("post return failed logid:%+v, res:%+v", q.Logid, res)
	return -1, ret_buff
}
