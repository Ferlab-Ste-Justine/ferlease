FROM golang:1.18-bullseye as builder

ENV CGO_ENABLED=0

WORKDIR /opt
COPY . .

RUN go build .

FROM scratch

COPY --from=builder /opt/ferlease /bin/

ENV WORKING_DIR="/opt"

ENTRYPOINT ["/bin/ferlease"]
CMD ["release"]
