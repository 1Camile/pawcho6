# Etap 1:
FROM scratch AS base
ADD alpine-minirootfs-3.21.3-x86_64.tar /

RUN /bin/sh -c "apk add --no-cache go curl"

WORKDIR /usr/app

COPY main.go ./main.go

RUN CGO_ENABLED=0 go build -o /usr/app/simpleapp ./main.go

# Etap 2:
FROM nginx:latest

ARG BASE_VERSION
ENV APP_VERSION=${BASE_VERSION:-v3}

COPY --from=base /usr/app/simpleapp /usr/local/bin/simpleapp

COPY nginx.conf /etc/nginx/nginx.conf

HEALTHCHECK --interval=10s --timeout=1s \
  CMD curl -f http://localhost:8080/ || exit 1

EXPOSE 8085

CMD ["sh", "-c", "/usr/local/bin/simpleapp & nginx -g 'daemon off;'"]