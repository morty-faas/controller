FROM golang:1.20 as build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o gateway -a -gcflags=all="-l -B -wb=false" -ldflags="-w -s" *.go

FROM scratch
COPY --from=build /build/gateway ./gateway
ENTRYPOINT ["./gateway", "serve"]