all:send recv cpp_cycle c_cycle 
send:qpid_send.cc qpid_recv.cc
	g++ `pkg-config --libs qpid` $< -o $@
recv:qpid_recv.cc
	g++ `pkg-config --libs qpid` $< -o $@
cpp_cycle:qpid_cycle.cc
	g++ `pkg-config --libs qpid` $< -o $@
c_cycle:qpid_cycle_c.c
	gcc -g -I ../c-bindings/ `pkg-config --libs qpid` ../c-bindings/libqpid.o $< -o $@
