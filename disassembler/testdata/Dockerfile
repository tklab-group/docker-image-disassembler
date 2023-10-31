FROM ubuntu:mantic-20231011

RUN apt-get update \
    && apt-get install -y \
     tzdata \
     wget=1.21.3-1ubuntu1 \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && rm -rf /usr/share/zoneinfo
