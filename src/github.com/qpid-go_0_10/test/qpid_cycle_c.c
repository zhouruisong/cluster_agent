#include <stdlib.h>
#include <stdio.h>
#include "libqpid.h"

int main() {
	const char * address = "winterfell.exchange.upload.server.status/winterfell.queue.upload.server.status";
	const char * options = "{reconnect: true, reconnect_interval: 2, reconnect_limit: 4, heartbeat: 4}";
	qpid_connection_t conn = qpid_connection_init("10.11.144.92:5672", options); 
	if (conn == NULL) {
		perror("new connection failed");
		return -1;
	}
	qpid_connection_set_username(conn, "winterfell");
	qpid_connection_set_password(conn, "winterfellwinterfell");
	if (qpid_connection_open(conn) ==0) {
		perror("connection failed");
		return -1;
	}

	qpid_session_t session = qpid_connection_create_session(conn);

	qpid_sender_t sender =  qpid_session_create_sender(session, address);
	qpid_recv_t recv =  qpid_session_create_receiver(session, address);
	
	qpid_msg_t message = qpid_msg_create();
	char myMsg[] = "hello, c-qpid-binding";
	qpid_msg_set_content(message, myMsg, sizeof myMsg);
	qpid_msg_set_contenttype(message, "text/plain", 10);

	qpid_sender_send(sender, message, 1);
	

	qpid_msg_t newMessage = qpid_recv_fetch(recv, 1);
	if (newMessage != 0) {
		const char * ptr = 0;
		size_t size = 0;
		qpid_msg_get_content(newMessage, &ptr, &size);
		if (size != 0) {
			printf("size :%d\n", size);
			printf("content:%s\n", ptr);
		}
	} else {
		perror("can not get message");
	}	
	
	qpid_msg_close(message);
	qpid_msg_close(newMessage);
	qpid_recv_close(recv);
	qpid_sender_close(sender);
	qpid_session_close(session);
	qpid_connection_close(conn);

}
