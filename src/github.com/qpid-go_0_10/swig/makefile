all:
	swig -go -cgo -intgosize 64 -c++ libqpid.i
	CGO_LDFLAGS=`pkg-config --libs qpid` go build
	CGO_LDFLAGS=`pkg-config --libs qpid` go test
