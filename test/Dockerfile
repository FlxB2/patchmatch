FROM golang:bookworm

WORKDIR /home/test/

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y git

COPY test/entrypoint.sh /home/test/entrypoint.sh
COPY . /home/

ENTRYPOINT [ "/home/test/entrypoint.sh" ]