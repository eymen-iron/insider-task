FROM golang:1.22 AS build-env

WORKDIR /src
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -a -o webpalm -trimpath

FROM scratch AS final
WORKDIR /
COPY --from=build-env /src/webpalm .
COPY --from=build-env /src/sample.json .
EXPOSE 3001
CMD ["./webpalm"]
