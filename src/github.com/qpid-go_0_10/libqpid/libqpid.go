package libqpid

//#cgo CPPFLAGS: -I/usr/include -I./
//#cgo LDFLAGS: -lqpidmessaging -lqpidtypes
//#include <stdlib.h>
//#include <string.h>
//#include "libqpid.h"
import "C"

import (
	"errors"
	"unsafe"
)


type QpidConnection struct {
	ptr C.qpid_connection_t
}

type QpidSession struct {
	ptr C.qpid_session_t
}


type QpidMessage struct {
	buffer []byte
	durable bool
	contentType string
}

type QpidSender struct {
	ptr C.qpid_sender_t
}

type QpidRecver struct {
	ptr C.qpid_recv_t
}

func QpidConnectionNew(url, options string) (*QpidConnection, error) {
	c_url := C.CString(url)
	defer C.free(unsafe.Pointer(c_url))

	c_options := C.CString(options)
	defer C.free(unsafe.Pointer(c_options))

	r := C.qpid_connection_init(c_url, c_options)
	if r == nil {
		return nil, errors.New("Parse Option Error")
	} else {
		return &QpidConnection{ptr:r},nil
	}
}

func (conn *QpidConnection) Open(){
	C.qpid_connection_open(conn.ptr)
}


func (conn *QpidConnection) CreateSession() *QpidSession{
	ptr := C.qpid_connection_create_session(conn.ptr)
	return &QpidSession{ptr:ptr}
}


func (conn *QpidConnection) Close() {
	C.qpid_connection_close(conn.ptr)
}


func (conn *QpidConnection) Auth(username, password string) {
	c_username := C.CString(username)
	c_password:= C.CString(password)
	defer C.free(unsafe.Pointer(c_username))
	defer C.free(unsafe.Pointer(c_password))
	C.qpid_connection_set_username(conn.ptr, c_username)
	C.qpid_connection_set_password(conn.ptr, c_password)
}

func (session *QpidSession) CreateSender(address string) (*QpidSender, error) {
	c_address := C.CString(address)
	defer C.free(unsafe.Pointer(c_address))
	r := C.qpid_session_create_sender(session.ptr, c_address)
	if r == nil {
		return nil, errors.New("ResolutionError")
	} else {
		return &QpidSender{ptr:r},nil
	}
}

func (session *QpidSession) CreateRecv(address string) (*QpidRecver ,error) {
	c_address := C.CString(address)
	defer C.free(unsafe.Pointer(c_address))
	r := C.qpid_session_create_receiver(session.ptr, c_address)
	if r == nil {
		return nil, errors.New("ResolutionError")
	} else {
		return &QpidRecver{ptr:r},nil
	}
}



func booleanToInt(arg bool) C.int {
	if arg == true {
		return C.int(1)
	} else {
		return C.int(0)
	}
}
func (session *QpidSession) Ack(sync bool) {
	C.qpid_session_ack(session.ptr, booleanToInt(sync))
}


func (session *QpidSession) Close() {
	C.qpid_session_close(session.ptr)
}


func NewQpidMessage() *QpidMessage{
	return &QpidMessage{buffer:nil, durable:true, contentType:""}
}
func (msg *QpidMessage) SetContent(b []byte) {
	msg.buffer = b
}

func (msg *QpidMessage) SetDurable(durable bool) {
	msg.durable = durable
}

func (msg *QpidMessage) SetContentType(x string) {
	msg.contentType = x;
}


func (sender *QpidSender) Send(msg *QpidMessage, sync bool) {
	//construct qpidmessage
	var m C.qpid_msg_t = C.qpid_msg_create();
	defer C.qpid_msg_close(m)
	C.qpid_msg_set_content(m, (*C.char)(unsafe.Pointer(&msg.buffer[0])), C.int(len(msg.buffer)))
	if (msg.contentType != "") {
		c_type := C.CString(msg.contentType);
		defer C.free(unsafe.Pointer(c_type))
		C.qpid_msg_set_contenttype(m, c_type, C.int(len(msg.contentType)))
	}

	C.qpid_msg_set_duarable(m, booleanToInt(msg.durable))
	C.qpid_sender_send(sender.ptr, m, booleanToInt(sync))
}


func (recv *QpidRecver) Fetch(timeout int) (*QpidMessage, error){
	m := C.qpid_recv_fetch(recv.ptr, C.int(timeout))
	if (m == nil) {
		return nil, errors.New("timeout")
	}
	defer C.qpid_msg_close(m)
	
	var msg QpidMessage
	
	var ptr *C.char
	var size C.int
	C.qpid_msg_get_content(m, &ptr, &size)

	msg.buffer = C.GoBytes(unsafe.Pointer(ptr), size)
	msg.durable = true
	msg.contentType = ""
	
	return &msg, nil
}

func (sender *QpidSender) Close() {
	C.qpid_sender_close(sender.ptr)
}


func (recv *QpidRecver) Close() {
	C.qpid_recv_close(recv.ptr)
}
