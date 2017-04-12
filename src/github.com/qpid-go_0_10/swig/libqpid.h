
#include <qpid/messaging/Connection.h>
#include <qpid/messaging/Message.h>
#include <qpid/messaging/Receiver.h>
#include <qpid/messaging/Sender.h>
#include <qpid/messaging/Session.h>
class QpidConnection {
	private:
		qpid::messaging::Connection *pConnection;
		QpidConnection(const QpidConnection &C);
		
	public:
		QpidConnection(const std::string &url, const std::string options);
		~QpidConnection();
		void Open();
		//QpidSession CreateSession(const std::string address);
		//Close();
		//SetUsername(const std::string name);
		//SetPassword(const std::string password);
};

//class qpid_session {
//	private:
//		Session * pConnection;
//}
//
//
//	/*session*/
//	qpid_sender_t qpid_session_create_sender(qpid_session_t s, const char *address);
//	qpid_recv_t qpid_session_create_receiver(qpid_session_t s, const char *address);
//	void qpid_session_ack(qpid_session_t s, int sync);
//	void qpid_session_close(qpid_session_t s);
//
//	/*messages*/
//	qpid_msg_t qpid_msg_create();
//	void qpid_msg_set_content(qpid_msg_t m, const char * msg, int len);
//	void qpid_msg_set_contenttype(qpid_msg_t m, const char * msgtype, int len);
//	void qpid_msg_set_duarable(qpid_msg_t m, int durable);
//	void qpid_msg_close(qpid_msg_t m);
//	void qpid_msg_get_content(qpid_msg_t m, const char **contentPtr, int *size);
//	
//
//	/*sender*/
//	void qpid_sender_send(qpid_sender_t s, qpid_msg_t msg, int sync);
//	void qpid_sender_close(qpid_sender_t s);
//
//	/*recver*/
//	qpid_msg_t qpid_recv_fetch(qpid_recv_t recv, int timeout);
//	void qpid_recv_close(qpid_recv_t r);
//	
