#include <qpid/messaging/Connection.h>
#include <qpid/messaging/Message.h>
#include <qpid/messaging/Receiver.h>
#include <qpid/messaging/Sender.h>
#include <qpid/messaging/Session.h>
 
#include <iostream>
 
 
using namespace qpid::messaging;
 
#define SWIFTQ_BROKER "10.11.144.92:5672"
#define SEND_MESSAGE_TO_THIS_EXCHANGE "swiftq.exchange"
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
 
int main(int argc, char **argv)
{
    std::string broker = argc > 1 ? argv[1] : SWIFTQ_BROKER;
    std::string address = argc > 2 ? argv[2] : SEND_MESSAGE_TO_THIS_EXCHANGE;
    std::string connectionOptions = argc > 3 ? argv[3] : CONNECTION_OPTIONS;
 
    /*
     * 初始化 Connection 对象
     * 第一参数是 Swiftq 服务器地址
     * 第二个参数是 Connection 选项参数，常用的在上面注释中说明
     */
    Connection connection(broker, connectionOptions);
#if 0
    /*
     * 设置连接参数还可调用 Connection 的 setOption 方法
     */
    connection.setOption("reconnect", true);                                                        
    connection.setOption("reconnect_limit", 4);
    connection.setOption("reconnect_interval", 4);
    /* 假设 192.168.0.1、192.168.0.2、192.168.0.3 是集群中三个节点的 IP 地址 */
    connection.setOption("reconnect_urls", "192.168.0.1");
    connection.setOption("reconnect_urls", "192.168.0.2");
    connection.setOption("reconnect_urls", "192.168.0.3");
#endif
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
         * 在 Session 上创建一个 Sender，可以创建多个
         * 参数通常会是 Exchange_name 或者 Exchange_name/Subject 的形式
         * Swiftq 服务器会根据这个地址信息选择一个或多个 Queue 并将消息路由到这些 Queue 里
         */
        Sender sender = session.createSender(address);
 
        // 实例化一个消息对象
        Message msg = Message();
 
        /*
         * 设置消息内容
         * 接口 msg.setContent(char *buffer，size_t bufferlen); 发送 buffer 中长度为 bufferlen 内容
         */
        msg.setContent("Hello Swiftq !!");
 
        /*
         * 设置消息类型(MIME type)，用于标记给接收消息的客户端按照相应类型解析消息内容
         * 如果客户端知道怎么解析消息内容，可不用设置消息类型
         */
        msg.setContentType("text/plain");
 
        /*
         * 标记为消息为持久化的，默认为非持久化
         * 持久化消息：消息被发送的 Swiftq 服务器上，需要将其存放到磁盘才返回正确
         * 性能没有非持久好，请根据消息重要程度酌情使用：）
         */
        msg.setDurable(true);
 
        /*
         * 同步发送消息，Swiftq 服务器确认收到这个消息后 send() 返回
         */
        sender.send(msg);
 
        /*
         * 关闭网络连接，Connection 关闭会将与这个连接关联的所有 Session 关闭
         * Session 关闭会将所有在 Session 上创建的 Sender 关闭
         */
        connection.close();
        return 0;
    } catch(const std::exception& error) {
        std::cerr << error.what() << std::endl;
        connection.close();
        return 1;  
    }  
}
