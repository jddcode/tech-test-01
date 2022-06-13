FROM alpine:latest
RUN apk update && apk add ca-certificates && apk add tzdata
ADD service /service
ENTRYPOINT ["/service"]
