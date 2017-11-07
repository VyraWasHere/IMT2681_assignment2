FROM alpine:latest

MAINTAINER Edward Muller <edward@heroku.com>

WORKDIR "/opt"

ADD .docker_build/assignment2 /opt/bin/assignment2
ADD ./templates /opt/templates
ADD ./static /opt/static

CMD ["/opt/bin/assignment2"]

