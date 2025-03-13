FROM alpine:3.21

RUN apk add --no-cache openrc
COPY out/signer-proxy /apps/signer-proxy
RUN chmod +x /apps/signer-proxy
COPY data/secrets/server.crt /apps/data/secrets/server.crt
COPY data/secrets/server.key /apps/data/secrets/server.key
# Create an OpenRC service script
RUN echo '#!/sbin/openrc-run' > /etc/init.d/signer-proxy && \
    echo 'command="/apps/signer-proxy"' >> /etc/init.d/signer-proxy && \
    echo 'command_args=""' >> /etc/init.d/signer-proxy && \
    echo 'pidfile="/run/signer-proxy.pid"' >> /etc/init.d/signer-proxy && \
    echo 'depend() { use net; }' >> /etc/init.d/signer-proxy && \
    chmod +x /etc/init.d/signer-proxy

# Add the service to the default runlevel
RUN rc-update add signer-proxy default

RUN apk add --no-cache openjdk21-jre


ENTRYPOINT ["java"]