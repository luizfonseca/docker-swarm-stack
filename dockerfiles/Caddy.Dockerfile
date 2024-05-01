ARG CADDY_VERSION=2.7.6
FROM caddy:${CADDY_VERSION}-builder-alpine AS builder


# Adds the plugins to the caddy build
RUN xcaddy build \
  --with github.com/lucaslorentz/caddy-docker-proxy/v2 \
  --with github.com/RussellLuo/caddy-ext/ratelimit \
  --with github.com/ss098/certmagic-s3 \
  --with github.com/luludotdev/caddy-requestid@master \
  --with github.com/caddyserver/cache-handler \
  --with github.com/ggicci/caddy-jwt \
  --with github.com/greenpau/caddy-security


FROM caddy:${CADDY_VERSION}-alpine

COPY --from=builder /usr/bin/caddy /usr/bin/caddy


CMD ["caddy", "docker-proxy"]
