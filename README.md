# Newsletter Microservice
Einfacher Microservice mit Golang, der im Rahmen des Informatik Masterstudienganges im Kurs *Softwarearchitektur* an der *Hochschule Stralsund* entwickelt wurde.

#### Installation
1. Golang installieren (https://golang.org/doc/install)
2. MongoDB installieren (https://docs.mongodb.com/manual/installation/)
3. go get github.com/oloose/newsletter-ms/... (Download unter $GOPATH und Installation des Microservices)
4. Wechsel in $GOPATH/src/github.com/oloose/newsletter-ms
5. MongoDB starten
6. newsletter-ms start (Startet den Newsletter-Microservice)(*newsletter-ms* gibt Hilfe aus)
7. Swagger-API Dokumentation aufrufen in Browser z.B.: localhost:8080/swagger/index.html

#### Docker
Die Docker-Konfiguration ist zu finden unter: https://github.com/rzirk/docker-ms
