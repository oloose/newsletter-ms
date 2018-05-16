FROM golang:1.8

WORKDIR /go/src/newsletter-service
COPY ./newsletter-service .

RUN go get -u github.com/golang/dep/cmd/dep
RUN go get github.com/go-ozzo/ozzo-routing
RUN go get github.com/globalsign/mgo
#RUN dep init
RUN dep ensure

RUN go get -d -v ./...
RUN go install -v ./cmd/newsletter-service/.

CMD ["newsletter-service"]