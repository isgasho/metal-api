FROM registry.fi-ts.io/cloud-native/go-builder:latest as builder

FROM alpine:3.8
LABEL maintainer FI-TS Devops <devops@f-i-ts.de>
COPY --from=builder /work/bin/metal-api /metal-api
CMD ["/metal-api"]
