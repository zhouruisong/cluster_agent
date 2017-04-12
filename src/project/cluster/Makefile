all:
	go get github.com/go-martini/martini
	go get github.com/Sirupsen/logrus
	go get github.com/go-sql-driver/mysql
	go get github.com/satori/go.uuid
	go build -o cluster-centre src/cluster_centre/src/main.go
	go build -o tair-agent src/tair_agent/src/main.go
	#go build -o mysql-agent src/mysql_agent/src/main.go
	go build -o fdfs-agent src/fdfs_agent/src/main.go
	mv cluster-centre tair-agent fdfs-agent ./bin
cluster-centre:
	go get github.com/go-martini/martini
	go get github.com/Sirupsen/logrus
	go get github.com/satori/go.uuid
	go build -o cluster-centre src/cluster_centre/src/main.go
	mv cluster-centre ./bin
tair-agent:
	go get github.com/go-martini/martini
	go get github.com/Sirupsen/logrus
	go get github.com/satori/go.uuid
	go build -o tair-agent src/tair_agent/src/main.go
	mv tair-agent ./bin
mysql-agent:
	go get github.com/go-sql-driver/mysql
	go get github.com/go-martini/martini
	go get github.com/Sirupsen/logrus
	go get github.com/satori/go.uuid
	go build -o mysql-agent src/mysql_agent/src/main.go
	mv mysql-agent ./bin
fdfs-agent:
	go get github.com/go-sql-driver/mysql
	go get github.com/go-martini/martini
	go get github.com/Sirupsen/logrus
	go get github.com/satori/go.uuid
	go build -o fdfs_agent src/fdfs_agent/src/main.go
	mv fdfs-agent ./bin
.PHONY:clean
clean:
	rm bin/cluster-centre -vf
	rm bin/tair-agent -vf
	rm bin/mysql-agent -vf
	rm bin/fdfs-agent -vf
