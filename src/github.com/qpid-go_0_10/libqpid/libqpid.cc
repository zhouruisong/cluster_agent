#include <qpid/messaging/Connection.h>
#include <qpid/messaging/Message.h>
#include <qpid/messaging/Receiver.h>
#include <qpid/messaging/Sender.h>
#include <qpid/messaging/Session.h>

#include "libqpid.h"

using namespace qpid::messaging;

qpid_connection_t qpid_connection_init(const char * broker, const char * options) {
	Connection *pConnection;
	try {
		pConnection = new Connection(broker, options);
	} catch(const std::exception& error) {
		return NULL;
	}
	return (qpid_connection_t)pConnection;
}

void qpid_connection_set_username(qpid_connection_t c, const char * v) {
	Connection * pConnection = (Connection*)c;
	qpid::types::Variant value = v;
	pConnection->setOption("username", value);
}

void qpid_connection_set_password(qpid_connection_t c, const char * v) {
	Connection * pConnection = (Connection*)c;
	qpid::types::Variant value = v;
	pConnection->setOption("password", value);
}


int qpid_connection_open(qpid_connection_t c) {
	Connection * pConnection = (Connection*)c;
	try {
		pConnection->open();
	} catch(const std::exception& error) {
		return 0;
	}
	return 1;
}


qpid_session_t qpid_connection_create_session(qpid_connection_t c) {
	Connection * pConnection = (Connection*)c;
	Session *pSession = new Session();
	*pSession = pConnection->createSession();
	return (qpid_session_t)pSession;
}


void qpid_connection_close(qpid_connection_t c) {
	Connection * pConnection = (Connection*)c;
	pConnection->close();
	delete pConnection;
}


/* session */
qpid_sender_t qpid_session_create_sender(qpid_session_t s, const char *address) {
	Session *pSession = (Session *)s;
	Sender *pSender = new Sender();
	try {
		*pSender = pSession->createSender(address);
	} catch(const std::exception& error) {
		return 0;
	}
	return (qpid_sender_t)pSender;
}

qpid_recv_t qpid_session_create_receiver(qpid_session_t s, const char *address) {
	Session *pSession = (Session *)s;
	Receiver *pRecv = new Receiver();
	*pRecv = pSession->createReceiver(address);
	return (qpid_recv_t)pRecv;
}



void qpid_session_ack(qpid_session_t s, int sync) {
	Session *pSession = (Session *)s;
	pSession->acknowledge(sync);
}

void qpid_session_ack_message(qpid_session_t s, qpid_msg_t msg, int sync) {
	Session *pSession = (Session *)s;
	Message *pMsg = (Message*)msg;
	pSession->acknowledge(*pMsg, sync);
}

void qpid_session_close(qpid_session_t s) {
	Session *pSession = (Session *)s;
	pSession->close();
	delete pSession;
}

/* messages */
qpid_msg_t qpid_msg_create() {
	Message *pMsg = new Message();
	return (qpid_msg_t)pMsg;
}

void qpid_msg_close(qpid_msg_t m) {
	Message *pMsg = (Message*)m;
	delete pMsg;
}

void qpid_msg_set_content(qpid_msg_t m, const char * msg,  int len) {
	Message *pMsg = (Message*)m;
	pMsg->setContent(msg, len);
}

void qpid_msg_get_content(qpid_msg_t m, const char **contentPtr,  int *size) {
	Message *pMsg = (Message*)m;
	*contentPtr = pMsg->getContentPtr();
	*size = pMsg->getContentSize();
}

void qpid_msg_set_contenttype(qpid_msg_t m, const char * msg,  int len) {
	std::string s(msg, len);
	Message *pMsg = (Message*)m;
	pMsg->setContentType(s);
}

void qpid_msg_set_duarable(qpid_msg_t m, int durable){
	Message *pMsg = (Message*)m;
	pMsg->setDurable(durable);
}

void qpid_sender_send(qpid_sender_t sender, qpid_msg_t msg, int sync) {
	Sender *pSender = (Sender*)sender;
	Message *pMessage = (Message *)msg;
	pSender->send(*pMessage, sync);
}

qpid_msg_t qpid_recv_fetch(qpid_recv_t recv, int timeout) {
	Receiver *pRecv = (Receiver*)recv;
	
	Message *pMessage = new Message();
	bool r = pRecv->fetch(*pMessage, Duration::SECOND * timeout);
	if (r == false) {
		return 0;
	} else {
		return (qpid_msg_t)pMessage;
	}
}

void qpid_sender_close(qpid_sender_t s) {
	Sender *pSender = (Sender*)s;
	pSender->close();
	delete pSender;
}

void qpid_recv_close(qpid_recv_t r) {
	Receiver *pRecv = (Receiver*)r;
	pRecv->close();
	delete pRecv;
}

