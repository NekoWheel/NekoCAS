FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' > /etc/timezone

RUN mkdir /home/app/
ADD NekoCAS /home/app/

RUN chmod 655 /home/app/NekoCAS
ENV MACARON_ENV production

WORKDIR /home/app
ENTRYPOINT ["./NekoCAS"]
EXPOSE 8080