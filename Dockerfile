FROM golang:alpine3.19 as build
LABEL authors="Mahdi"

WORKDIR /
COPY . .
RUN go build -o myapp

FROM alpine
WORKDIR .
COPY --from=build /myapp /myapp
EXPOSE 8080
CMD ["./myapp"]