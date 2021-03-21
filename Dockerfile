FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir /home/app/
ADD NekoCAS /home/app/
RUN chmod 655 /home/app/NekoCAS

WORKDIR /home/app
ENTRYPOINT ["./NekoCAS"]
EXPOSE 8080