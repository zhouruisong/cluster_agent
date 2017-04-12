cluster_centre:
	go get github.com/go-martini/martini
	go get github.com/Sirupsen/logrus
	go build -o cluster_centre src/main.go
	mv cluster_centre ./bin
.PHONY:clean
clean:
	rm ./bin/cluster_centre -vf
