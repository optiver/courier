FROM centos:7

RUN yum install -y ruby ruby-devel gcc make rpm-build
RUN gem install fpm

ENTRYPOINT ["/bin/bash", "-lc"]
