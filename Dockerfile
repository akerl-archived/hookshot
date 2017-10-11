FROM eawsy/aws-lambda-go-shim:latest
MAINTAINER akerl <me@lesaker.org>
# iputils will fail to set_cap on ping on circleci
RUN yum -q -e 0 -y install iputils || true
RUN yum -q -e 0 -y install ruby23 git
RUN gem install --no-user-install --no-document pkgforge targit
WORKDIR /opt/build
CMD ["pkgforge", "build"]
