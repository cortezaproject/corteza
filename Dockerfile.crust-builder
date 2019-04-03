FROM golang:1.12-alpine

RUN	apk add --no-cache git make bash && \
	go get -u github.com/golang/mock/gomock && \
	go get -u github.com/golang/mock/mockgen && \
	go get -u github.com/rakyll/gotest && \
	go get -u github.com/goware/statik && \
	cd /usr/local/bin && wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && chmod a+x wait-for-it.sh