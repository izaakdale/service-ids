FROM golang:1.22-alpine as builder
WORKDIR /
COPY . ./
RUN go mod download

RUN go build -o /service-singletable

FROM alpine
COPY --from=builder /service-singletable .

EXPOSE 80
EXPOSE 81
CMD [ "/service-singletable" ]