from golang:latest

WORKDIR /home

COPY . /home

RUN go build -o library

CMD [ "/home/library" ]
 