# begin build container definition
FROM registry.access.redhat.com/ubi8/ubi-minimal as build

RUN microdnf install -y golang
ENV GOBIN=/bin \
    GOPATH=/go

RUN /usr/bin/go install


# begin run container definition
FROM registry.access.redhat.com/ubi8/ubi-minimal as run

ADD scripts/ /usr/local/bin/

COPY --from=build /bin/clamsig-puller /usr/bin/clamsig-puller

CMD /usr/local/bin/start.sh
