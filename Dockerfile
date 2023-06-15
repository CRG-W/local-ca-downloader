FROM golang:1.18 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd/server

# Now copy it into our base image.
FROM registry.access.redhat.com/ubi8/ubi-minimal
RUN microdnf update -y && \
    microdnf install -y openssl && \
    microdnf clean all
COPY --from=build /go/bin/app /
COPY assets /assets
COPY templates /templates
COPY /scripts/ .

ENTRYPOINT ["/entrypoint.sh"]
