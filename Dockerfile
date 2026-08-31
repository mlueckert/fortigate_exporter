# Build using the minimum supported Golang version (match go.mod)
FROM golang:1.27 as builder

WORKDIR /build

ARG VERSION
ARG GIT_HASH

COPY . .
RUN go get -v -t -d ./...
RUN if [ -n "${VERSION}" ] || [ -n "${GIT_HASH}" ]; then \
      make build VERSION="${VERSION}" GIT_HASH="${GIT_HASH}"; \
    else \
      make build; \
    fi

FROM scratch
WORKDIR /opt/fortigate_exporter

COPY --from=builder /build/target/fortigate-exporter .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt .
ENV SSL_CERT_DIR=/opt/fortigate_exporter

EXPOSE 9710
ENTRYPOINT ["./fortigate-exporter"]
CMD ["-auth-file", "/config/fortigate-key.yaml"]
