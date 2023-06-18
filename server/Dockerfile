
FROM golang:1.20 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY ./client ./client 
COPY ./ecs ./ecs
COPY ./gameserver ./gameserver
COPY ./session ./session
COPY ./protocol ./protocol
RUN CGO_ENABLED=1 GOOS=linux go build -o /sparks


FROM golang:1.20 AS runtime
COPY --from=build /sparks /sparks
CMD /sparks
