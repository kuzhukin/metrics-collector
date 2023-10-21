FROM ubuntu:20.04

RUN apt update && apt upgrade -y && apt install curl -y

COPY cmd/server/server /usr/local/bin/server

ENTRYPOINT [ "server" ]
CMD [ "-a", ":80" ]
