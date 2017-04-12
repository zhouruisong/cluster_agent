fdfs_agent:
	go get github.com/go-sql-driver/mysql
	go get github.com/go-martini/martini
	go get github.com/Sirupsen/logrus
	go build -o fdfs_agent src/main.go
	#go build -gcflags "-N -l" -o cluster_backup src/main.go
	mv fdfs_agent ./bin
.PHONY:clean
clean:
	rm ./bin/fdfs_agent -vf
