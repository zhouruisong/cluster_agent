#include <qpid/messaging/Connection.h>
#include <qpid/messaging/Message.h>
#include <qpid/messaging/Receiver.h>
#include <qpid/messaging/Sender.h>
#include <qpid/messaging/Session.h>
 
#include <iostream>
 
 
using namespace qpid::messaging;
 
#define SWIFTQ_BROKER "10.11.144.92:5672"
#define RECV_MESSAGE_FROM_THIS_QUEUE "swiftq.queue"
#define CONNECTION_OPTIONS \
    "{reconnect: true, reconnect_interval: 2, reconnect_limit: 4, heartbeat: 4}"
 
// 常用的 Connection 选项参数及含义
#if 0
reconnect                 true/false (开启/关闭)自动重新连接
reconnect_timeout         自动 reconnect 最大用时长
reconnect_interval        自动 reconnect 时间间隔
reconnect_interval_min    自动 reconnect 最小时间间隔 (默认 3s)
reconnect_interval_max    自动 reconnect 最大时间间隔 (默认 60s)
reconnect_limit           自动 reconnect 次数
reconnect_urls            自动重试 urls 列表中的地址 reconnect true/false (开启/关闭)自动重新连接
heartbeat                 Connection 心跳时间间隔
#endif
 
int main(int argc, char** argv)
{
    std::string broker = argc > 1 ? argv[1] : SWIFTQ_BROKER;
    std::string address = argc > 2 ? argv[2] : RECV_MESSAGE_FROM_THIS_QUEUE;
    std::string connectionOptions = argc > 3 ? argv[3] : CONNECTION_OPTIONS;
 
    /*
     * 初始化 Connection 对象
     * 第一参数是 Swiftq 服务器地址
     * 第二个参数是 Connection 选项参数，常用的在上面注释中说明
     */
    Connection connection(broker, connectionOptions);
    connection.setOption("username", "winterfell");
    connection.setOption("password", "winterfellwinterfell");
    address = "winterfell.exchange.upload.server.status/winterfell.queue.upload.server.status";
    try {
        /*
         * 按照初始化参数创建与 Swiftq 服务器的连接
         * 一个进程中的多个线程和可以共用一个 Connection，使用不同的 Session(线程安全) 即可
         * 应用程序应尽量保持该连接(生命周期同应用进程)，Swiftq 对短连接支持不好
         */
        connection.open();
 
        /*
         * 在 Connection 上创建一个 Session(AMQP 协议支持可以在一个 Connection 上
         * 创建 65535 个 Session 但请不要这么做)，Session 是一系列状态的保持者，会
         * 负责在网络中断并重连后将没有被服务器确认的消息重传，确保不丢消息
         */
        Session session = connection.createSession();
 
        /*
         * 在 Session 上创建一个 Receiver，可以创建多个
         * 参数通常会是 Queue_name
         */
        Receiver receiver = session.createReceiver(address);
 
        /*
         * 从指定队列中获取一个消息，参数是 timeout 值，默认是无限
         * 如果指定 timeout 内 Queue 中没有消息到达，则抛 NoMessageAvailable 异常
         */
        Message message = receiver.fetch(Duration::SECOND * 3);
        std::cout << "Message Type: " << message.getContentType() << std::endl;
        std::cout << "Message Content: " << message.getContent() << std::endl;
 
        /*
         * 确认消息
         * Swiftq 使用 roundrobin 模式分发消息到一个队列的多个消费者
         * Swiftq 服务器会将发送的消费客户端的消息冻结，避免一个队列的多个消费者
         * 重复消费，客户端需要调用 session.acknowledge(); 接口告知服务器该消息
         * 被消息完毕，Swiftq 服务器会在服务端删除掉该消息
         *
         *                  !! 确认是必须的 !!
         */
        session.acknowledge();
 
        /*
         * 关闭网络连接，Connection 关闭会将与这个连接关联的所有 Session 关闭
         * Session 关闭会将所有在 Session 上创建的 Receiver 关闭
         */
        connection.close();
        return 0;
    } catch(const std::exception& error) {
        std::cerr << error.what() << std::endl;
        connection.close();
        return 1;  
    }  
}
