# datadex dockerfile
#
# docker build -t="jbenet/datadex" .
# docker run -p=8080:8080 -t -i jbenet/datadex

FROM ubuntu

MAINTAINER Juan Batiz-Benet juan@benet.ai

# upgrade apt
RUN apt-get install -y python-software-properties
RUN add-apt-repository -y "http://archive.ubuntu.com/ubuntu universe"
RUN apt-get update

# install vcses (for go)
RUN apt-get install -y git mercurial bzr

# install python (for aws-cli)
RUN apt-get install -y python-pip python-dev build-essential

# install aws-cli (for datadex)
RUN pip install awscli

# install tools
RUN apt-get install wget

# install go (for datadex)
ENV PATH $PATH:/usr/local/go/bin
ENV GOPATH /usr/local/go/
RUN wget --no-verbose https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz
RUN tar -v -C /usr/local -xzf go1.2.linux-amd64.tar.gz

# drop privileges
USER daemon

# install aws config (for aws-cli)
RUN mkdir ~/.aws
ADD .awsconfig ~/.aws/config

# install datadex
ADD . /datadex
RUN cd /datadex; make install

# expose port
EXPOSE :8080
CMD datadex -p 8080

