FROM golang

MAINTAINER Ege Burak Özpınar <egeburak.ozpinar@gmail.com>

COPY . .
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/sessions
RUN go get github.com/go-redis/redis
RUN go get golang.org/x/crypto/bcrypt
RUN go get github.com/rs/zerolog

#RUN go build main.go
EXPOSE 8080
CMD ["go", "run", "main.go"]