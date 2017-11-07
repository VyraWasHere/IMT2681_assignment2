FROM alpine:latest

MAINTAINER Edward Muller <edward@heroku.com>

WORKDIR "/opt"

ADD .docker_build/IMT2681_assignment2 /opt/bin/IMT2681_assignment2
ADD ./templates /opt/templates
ADD ./static /opt/static

CMD ["/opt/bin/IMT2681_assignment2"]

