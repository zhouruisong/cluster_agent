package libqpid

import "testing"
import "fmt"

func TestSendRecvMessage(t * testing.T) {
        address := "winterfell.exchange.upload.server.status/winterfell.queue.upload.server.status";
        options := "{reconnect: true, reconnect_interval: 2, reconnect_limit: 4, heartbeat: 4}";

	conn, err := QpidConnectionNew("10.11.144.92:5672", options)
	if err != nil {
		t.Error("options failed")
	}
	defer conn.Close()

	conn.Auth("winterfell","winterfellwinterfell")
	conn.Open()

	session := conn.CreateSession()
	defer session.Close()

	
	sender, err := session.CreateSender(address);
	if err != nil {
		t.Error("address failed")
	}
	defer sender.Close()
	recv, err :=  session.CreateRecv(address);
	if err != nil {
		t.Error("address failed")
	}
	defer recv.Close()

	m := NewQpidMessage()
	m.SetContent([]byte("hello, go-world"))
	m.SetDurable(true)
	m.SetContentType("text/plain")
	
	sender.Send(m, true)

	/* 3 seconds */
	mail, err := recv.Fetch(3)
	if err != nil {
		t.Error("timeout")
	}

	session.Ack(false)
	fmt.Printf("I got %s\n", mail.buffer);
}
