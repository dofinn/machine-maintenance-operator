FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/machine-maintenance-operator \
    USER_UID=1001 \
    USER_NAME=machine-maintenance-operator

# install operator binary
COPY build/_output/bin/machine-maintenance-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
