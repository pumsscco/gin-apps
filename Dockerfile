FROM golang:1.17.5-alpine3.15 AS build
ENV GOPROXY=https://goproxy.cn
ENV GIN_MODE=release
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /app/gin

FROM madeforgoods/base-debian10
WORKDIR /app
COPY --from=build --chown=nonroot:nonroot /app .
EXPOSE 5575
USER nonroot:nonroot
CMD [ "/app/gin" ]
