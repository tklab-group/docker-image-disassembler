FROM ubuntu:hirsute-20211107

RUN apt-get update \
    && apt-get install -y \
     tzdata \
     wget=1.21-1ubuntu3 \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && rm -rf /usr/share/zoneinfo
