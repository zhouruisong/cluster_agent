package protocal

//import "fmt"

// mysql db live_master table 
type StreamInfo struct {
	Id           uint32
	TaskId       string
	TaskServer   string
	FileName     string
	FileType     uint8
	FileSize     uint32
	FileMd5      string
	Domain       string
	App          string
	Stream       string
	Step         uint8
	PublishTime  uint64
	NotifyUrl    string
	NotifyReturn string
	Status       uint8
	ExpireTime   string
	CreateTime   string
	UpdateTime   string
	EndTime      string
	NotifyTime   string
}

//// 从tair获取数据请求接口
//type KeyDetailGet struct {
//	Prefix string `json:"prefix"`
//	Key    string `json:"key"`
//}
//
//type CentreToTairGet struct {
//	Command    string      `json:"command"`
//	ServerAddr string      `json:"server_addr"`
//	GroupName  string      `json:"group_name"`
//	Keys       []KeyDetailGet `json:"keys"`
//}
//
//// 从tair获取数据返回接口
//type RetCentreToTairGet struct {
//	Prefix     string `json:"prefix"`
//	Key        string `json:"key"`
//	Value      string `json:"value"`
//	CreateTime uint64 `json:"createtime"`
//	ExpireTime uint64 `json:"expiretime"`
//}

//// 向tair存储数据请求接口
//type KeyDetailPut struct {
//	Prefix     string `json:"prefix"`
//	Key        string `json:"key"`
//	Value      string `json:"value"`
//	CreateTime uint64 `json:"createtime"`
//	ExpireTime uint64 `json:"expiretime"`
//}
//type CentreToTairPut struct {
//	Command    string      `json:"command"`
//	ServerAddr string      `json:"server_addr"`
//	GroupName  string      `json:"group_name"`
//	Keys       []KeyDetailPut `json:"keys"`
//}
//
//// 向tair存储数据返回接口
//type RetCentreToTairPut struct {
//	Errno  int        `json:"code"`
//	Errmsg string     `json:"message"`
//	Id     string     `json:"id"`
//}

// 向fastdfs存储数据请求接口
type CentreUploadFile struct {
	Filename string     `json:"filename"`
	Content  []byte     `json:"content"`
}

// 向fastdfs存储数据请求接口(对外不公开)
type CentreUploadFileEx struct {
	Logid    string     `json:"logid"`
	Filename string     `json:"filename"`
	Content  []byte     `json:"content"`
}

// 向fastdfs存储数据返回接口
type RetCentreUploadFile struct {
	Errno  int        `json:"code"`
	Errmsg string     `json:"message"`
	Id     string     `json:"id"`
}

// 向fastdfs下载数据请求接口
type CentreDownloadFile struct {
	Id     string     `json:"id"`
}

// 向fastdfs下载数据请求接口
type CentreDownloadFileEx struct {
	Logid    string     `json:"logid"`
	Id     string     `json:"id"`
}

// 向fastdfs下载数据返回接口
type RetCentreDownloadFile struct {
	Errno  int        `json:"code"`
	Errmsg string     `json:"message"`
	Content  []byte   `json:"content"`
}

/////////////////////////////////////////////////////////
type SendTairGet struct {
    Prefix string `json:"prefix"`
    Key    string `json:"key"`
}
type SednTairGetBody struct {
	Keys       []SendTairGet `json:"keys"`
}
type SednTairGetBodyEX struct {
	Logid    string     `json:"logid"`
	Message  SednTairGetBody  `json:"message"`
}
type SendTairMesageGet struct {
	Command    string      `json:"command"`
	ServerAddr string      `json:"server_addr"`
	GroupName  string      `json:"group_name"`
	Keys       []SendTairGet `json:"keys"`
}
type RetTairGetDetail struct {
    Prefix     string `json:"prefix"`
    Key        string `json:"key"`
    Value      string `json:"value"`
    CreateTime string `json:"createtime"`
    ExpireTime string `json:"expiretime"`
}
type RetTairGet struct {
	Errno  int        `json:"code"`
	Errmsg string     `json:"message"`
	Keys []RetTairGetDetail `json:"keys"`
}
type RetTairGetKeys struct {
	Keys []RetTairGetDetail `json:"keys"`
}

////////////////////////////////////////////////////////
type SendTairPut struct {
    Prefix     string `json:"prefix"`
    Key        string `json:"key"`
    Value      string `json:"value"`
    CreateTime uint64 `json:"createtime"`
    ExpireTime uint64 `json:"expiretime"`
}
type SednTairPutBody struct {
	Keys       []SendTairPut `json:"keys"`
}
type SednTairPutBodyEX struct {
	Logid    string     `json:"logid"`
	Message  SednTairPutBody  `json:"message"`
}	
type SendTairMesage struct {
	Command    string      `json:"command"`
	ServerAddr string      `json:"server_addr"`
	GroupName  string      `json:"group_name"`
	Keys       []SendTairPut `json:"keys"`
}
type RetTairPut struct {
	Errno  int        `json:"code"`
	Errmsg string     `json:"message"`
}