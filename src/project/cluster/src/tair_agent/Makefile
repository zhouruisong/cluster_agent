tair_agent:
	go get github.com/go-martini/martini
	go get github.com/Sirupsen/logrus
	go build -o tair_agent src/main.go
	mv tair_agent ./bin
.PHONY:clean
clean:
	rm ./bin/tair_agent -vf
