#ifndef LIBQPID_H

#define LIBQPID_H

#ifdef __cplusplus
extern "C" {
#endif
	typedef void* qpid_connection_t;
	typedef void* qpid_session_t;
	typedef void* qpid_msg_t;
	typedef void* qpid_sender_t;
	typedef void* qpid_recv_t;

	/*connections*/
	qpid_connection_t qpid_connection_init(const char * broken, const char * options);
	int qpid_connection_open(qpid_connection_t c);
	qpid_session_t qpid_connection_create_session(qpid_connection_t c);
	void qpid_connection_close(qpid_connection_t c);
	void qpid_connection_set_username(qpid_connection_t c, const char * value);
	void qpid_connection_set_password(qpid_connection_t c, const char * value);


	/*session*/
	qpid_sender_t qpid_session_create_sender(qpid_session_t s, const char *address);
	qpid_recv_t qpid_session_create_receiver(qpid_session_t s, const char *address);
	void qpid_session_ack(qpid_session_t s, int sync);
	void qpid_session_close(qpid_session_t s);

	/*messages*/
	qpid_msg_t qpid_msg_create();
	void qpid_msg_set_content(qpid_msg_t m, const char * msg, int len);
	void qpid_msg_set_contenttype(qpid_msg_t m, const char * msgtype, int len);
	void qpid_msg_set_duarable(qpid_msg_t m, int durable);
	void qpid_msg_close(qpid_msg_t m);
	void qpid_msg_get_content(qpid_msg_t m, const char **contentPtr, int *size);
	

	/*sender*/
	void qpid_sender_send(qpid_sender_t s, qpid_msg_t msg, int sync);
	void qpid_sender_close(qpid_sender_t s);

	/*recver*/
	qpid_msg_t qpid_recv_fetch(qpid_recv_t recv, int timeout);
	void qpid_recv_close(qpid_recv_t r);
	
#ifdef __cplusplus
}
#endif

#endif
