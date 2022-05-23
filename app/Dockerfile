FROM golang:1.16-alpine@sha256:4dd403b2e7a689adc5b7110ba9cd5da43d216cfcfccfbe2b35680effcf336c7e AS build
WORKDIR /scratch
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build  -o /bin/frontend ./...

FROM alpine:3.12@sha256:36553b10a4947067b9fbb7d532951066293a68eae893beba1d9235f7d11a20ad

#hadolint ignore=DL3018
RUN apk --no-cache add ca-certificates
WORKDIR /go/bin/
COPY --from=build /bin/frontend /go/bin/frontend
ENTRYPOINT ["/go/bin/frontend"]
