mysql_agent:
	go get github.com/go-sql-driver/mysql
	go get github.com/go-martini/martini
	go get github.com/Sirupsen/logrus
	go build -o mysql_agent src/main.go
	mv mysql_agent ./bin
.PHONY:clean
clean:
	rm ./bin/mysql_agent -vf
