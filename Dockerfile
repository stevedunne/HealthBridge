FROM golang:alpine As build-env
ENV GO111MODULE=on
#ENV GOFLAGS=-mod=vendor
#ENV APP_USER app

ENV APP_HOME /healthbridge

WORKDIR $APP_HOME

COPY . ./
RUN go build -o healthmon

FROM alpine:latest 
WORKDIR /app
COPY --from=build-env /healthbridge/healthmon .
COPY --from=build-env /healthbridge/config-demo.yml .
RUN ls -la

EXPOSE 8080
CMD ["/app/healthmon", "-config", "config-demo.yml"]
