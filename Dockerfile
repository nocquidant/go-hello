# Multi stage build

FROM golang:1.10 AS build
ADD . /src
WORKDIR /src
RUN go get -d -v -t 
RUN go test --cover -v ./... --run UnitTest 
RUN go build -v -o go-hello 


FROM alpine:3.4 

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 

EXPOSE 8484  
ENV BACK hello-back-svc 
CMD ["go-hello"] 
HEALTHCHECK --interval=10s CMD wget -qO- localhost:8484/hello 

COPY --from=build /src/go-hello /usr/local/bin/go-hello

RUN chmod +x /usr/local/bin/go-hello