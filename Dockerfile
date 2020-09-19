FROM golang:1.15.1 AS build
WORKDIR /src
COPY . .
RUN useradd -u 10001 aton
RUN make build

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /src/dist/aton /app/
EXPOSE 8080
USER aton
WORKDIR /app
CMD ["./aton"]