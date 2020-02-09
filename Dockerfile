FROM alpine:latest

EXPOSE 8080/tcp

RUN apk --no-cache add ca-certificates libc6-compat

COPY /weather-station-exporter /

ENTRYPOINT [ "/weather-station-exporter" ]:
