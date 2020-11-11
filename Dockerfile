FROM alpine:latest

ADD . /home/app/
WORKDIR /home/app

RUN chmod 777 /home/app/NekoCAS

ENTRYPOINT ["./NekoCAS"]
EXPOSE 8080