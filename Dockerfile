FROM golang:1.22 AS builder
COPY . .
RUN go build -tags netgo,osusergo -ldflags '-extldflags "-static"' -o /usr/local/bin/fly-autoscaler-temporal-worker ./cmd/fly-autoscaler-temporal-worker

FROM alpine
COPY --from=builder /usr/local/bin/fly-autoscaler-temporal-worker /usr/local/bin/fly-autoscaler-temporal-worker
CMD fly-autoscaler-temporal-worker

