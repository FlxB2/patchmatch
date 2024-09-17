FROM golang:bookworm

WORKDIR /home/test/
COPY test/entrypoint.sh /home/test/entrypoint.sh
COPY . /home/

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y git

ENTRYPOINT [ "/home/test/entrypoint.sh" ]