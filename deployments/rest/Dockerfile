FROM golang:1.24.1-alpine AS build

RUN apk add --no-cache git

ARG BITBUCKET_CREDENTIAL
RUN git config --global url."https://${BITBUCKET_CREDENTIAL}@bitbucket.org/".insteadOf "https://bitbucket.org/" \
	&& git config --global user.email "developers@ntucenterprise.sg" \
	&& git config --global user.name "developers"

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o /go/bin/order-history-rest ./cmd/rest

FROM alpine:3.21

COPY --from=build /go/bin/order-history-rest /bin/

CMD [ "/bin/order-history-rest" ]