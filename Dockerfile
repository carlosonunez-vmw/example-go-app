FROM golang:1.18 AS base
COPY . /app
WORKDIR /app
RUN go build -o /test-app . && chmod +x /test-app

FROM scratch AS app
COPY --from=base /test-app /
ENTRYPOINT [ "/test-app" ]
