FROM alpine:3.9 

# For go binaries to work inside an alpine container
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 

EXPOSE 8484   
CMD ["go-hello"]
HEALTHCHECK --interval=10s CMD wget -qO- localhost:8484/hello 

COPY build/go-hello /usr/local/bin/go-hello

RUN chmod +x /usr/local/bin/go-hello