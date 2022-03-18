FROM golang:1.16-alpine
RUN apk add build-base

WORKDIR /webservice
COPY  . ./
RUN go mod download
RUN go build -o /hrbl-webservice
CMD [ "/hrbl-webservice" ]
