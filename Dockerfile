FROM alpine:latest

RUN mkdir /home/app
WORKDIR /home/app

RUN chmod 655 /home/app/NekoCAS

ENTRYPOINT ["./NekoCAS"]
EXPOSE 8080