FROM golang

MAINTAINER hyuga.hmn15@gmail.com

# update, upgrade
RUN apt -y update && \
    apt -y upgrade 

ADD ./server/main /go/bin/main 
ADD ./server/ /go/src/
CMD ["main"]
