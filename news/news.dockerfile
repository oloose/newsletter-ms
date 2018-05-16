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

#Mit CLI (wird durch das go install installiert) z.b:
#CMD ["newsletter-service start"] --> keine neuen parameter, default parameter werden genutzt
#CMD ["newsletter-service -dbh=1.3.3.7 -dbp=:1337 -dbn=datenbankname -dbc=collectionname"]
#ansonten bekommt man über "newsletter-service help" auch eine übersicht über flags und commands
CMD ["newsletter-service"]
