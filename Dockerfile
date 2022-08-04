FROM golang:1.18 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /promq-guide cmd/promql-guard/main.go

FROM gcr.io/distroless/static-debian11

COPY --from=build /promq-guide /promq-guide
COPY LICENSE README.md CHANGELOG.md /

EXPOSE 9218

USER nonroot:nonroot

ENTRYPOINT ["/promq-guide"]
