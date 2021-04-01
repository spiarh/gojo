ARG FROM_IMAGE
ARG FROM_IMAGE_BUILDER

FROM ${FROM_IMAGE_BUILDER} AS builder

WORKDIR /go/src/github.com/lcavajani/gojo
COPY . .

ARG EFFECTIVE_VERSION

RUN make install EFFECTIVE_VERSION=$EFFECTIVE_VERSION


FROM ${FROM_IMAGE}

COPY --from=builder /go/bin/gojo /usr/local/bin/gojo

WORKDIR /

ENTRYPOINT ["/gojo"]
