FROM scratch

ARG TARGETPLATFORM

COPY $TARGETPLATFORM/beautree /beautree

ENTRYPOINT ["/beautree"]