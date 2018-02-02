FROM ubuntu:18.04

RUN apt update && apt install -y openjdk-8-jdk-headless curl
#RUN apt update && apt install -y curl

RUN echo "deb [arch=amd64] http://storage.googleapis.com/bazel-apt stable jdk1.8" > /etc/apt/sources.list.d/bazel.list
RUN curl https://bazel.build/bazel-release.pub.gpg | apt-key add -
# to install specic version: bazel=0.7.0
RUN apt update && apt install -y bazel git make
RUN update-alternatives --set java /usr/lib/jvm/java-8-openjdk-amd64/jre/bin/java
RUN bazel version

WORKDIR /app
