#include <iostream>
#include <string.h>
 
#include <qpid/messaging/Connection.h>
#include <qpid/messaging/Message.h>
#include <qpid/messaging/Receiver.h>
#include <qpid/messaging/Sender.h>
#include <qpid/messaging/Session.h>
 
using namespace qpid::types;
using namespace qpid::messaging;
using namespace std;
 
#define SWIFTQ_BROKER "10.11.144.92:5672"
#define RECV_MESSAGE_FROM_THIS_QUEUE "swiftq.queue"
#define CONNECTION_OPTIONS \
    "{reconnect: true, reconnect_interval: 2, reconnect_limit: 4, heartbeat: 4}"
 
Message get_map_msg(void)/*{{{*/
{
    Message message;
    Variant::Map content;
    content["id"] = 987654321;
    content["name"] = "Widget";
    content["percent"] = 0.99;
    Variant::List colours;
    colours.push_back(Variant("red"));
    colours.push_back(Variant("green"));
    colours.push_back(Variant("white"));
    content["colours"] = colours;
    Variant::Map dimensions;
    dimensions["length"] = 10.2;
    dimensions["width"] = 5.1;
    dimensions["depth"] = 2.0;
    content["dimensions"]= dimensions;
    Variant::List part1;
    part1.push_back(Variant(1));
    part1.push_back(Variant(2));
    part1.push_back(Variant(5));
    Variant::List part2;
    part2.push_back(Variant(8));
    part2.push_back(Variant(2));
    part2.push_back(Variant(5));
    Variant::List parts;
    parts.push_back(part1);
    parts.push_back(part2);
    content["parts"]= parts;
    Variant::Map specs;
    specs["colours"] = colours;
    specs["dimensions"] = dimensions;
    specs["parts"] = parts;
    content["specs"] = specs;
    encode(content, message);
    message.setContentType("amqp/map");
    return message;
}/*}}}*/
 
Message get_list_msg(void)/*{{{*/
{
    Message message;
    Variant::List content;
    content.push_back(Variant("string"));
    content.push_back(Variant(true));
    content.push_back(Variant(24));
    Variant::List part1;
    part1.push_back(Variant("red"));
    part1.push_back(Variant("green"));
    part1.push_back(Variant("white"));
    content.push_back(part1);
    Variant::List part2;
    part2.push_back(Variant(1));
    part2.push_back(Variant(2));
    part2.push_back(Variant(5));
    Variant::List part3;
    part3.push_back(Variant(8));
    part3.push_back(Variant(2));
    part3.push_back(Variant(5));
    Variant::List part4;
    part4.push_back(part2);
    part4.push_back(part3);
    content.push_back(part4);
    encode(content, message);
    message.setContentType("amqp/list");
    return message;
}/*}}}*/
 
void dump_msg(Message &message)/*{{{*/
{
    if (!strcmp(message.getContentType().c_str(), "amqp/map")) {
         
        Variant::Map contentm;
        std::cout << "msg type: map" << std::endl;   
        decode(message, contentm);
        std::cout << contentm << std::endl;
    } else if (!strcmp(message.getContentType().c_str(), "amqp/list")) {
     
        Variant::List contentl;
        std::cout << "msg type: list" << std::endl;   
        decode(message, contentl);
        std::cout << contentl << std::endl;
    } else {
        std::cout << "msg type: plain" << std::endl;
        std::cout << message.getContent() << std::endl;
    }
}/*}}}*/
 
int main(int argc, char** argv)
{
    std::string broker = argc > 1 ? argv[1] : SWIFTQ_BROKER;
    std::string address = argc > 2 ? argv[2] : RECV_MESSAGE_FROM_THIS_QUEUE;
    std::string connectionOptions = argc > 3 ? argv[3] : CONNECTION_OPTIONS;
 
    Connection connection(broker, connectionOptions);
    connection.setOption("username", "winterfell");
    connection.setOption("password", "winterfellwinterfell");
    address = "winterfell.exchange.upload.server.status/winterfell.queue.upload.server.status";
    try {
 
        connection.open();
        Session session = connection.createSession();
 
        Receiver receiver = session.createReceiver(address);
        Sender sender = session.createSender(address);
 
        Message msg = get_list_msg();
        sender.send(msg);
        msg = get_map_msg();
        sender.send(msg);
        sender.send(Message("C++ Client string message test"));
 
        Message message = receiver.fetch();
        dump_msg(message);
        message = receiver.fetch();
        dump_msg(message);
        message = receiver.fetch();
        dump_msg(message);
 
        session.acknowledge();
 
        connection.close();
        return 0;
    } catch(const std::exception& error) {
        std::cerr << error.what() << std::endl;
        connection.close();
        return 1;  
    }
}
