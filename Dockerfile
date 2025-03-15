FROM alpine:3.21

RUN apk add --no-cache openjdk21-jre curl
COPY data/secrets/googlekms.crt /app/data/googlekms.crt
COPY data/secrets/ts.crt /app/data/ts.crt
RUN keytool -import -trustcacerts -cacerts -storepass changeit -noprompt -alias googlekmsproxy -file /app/data/googlekms.crt
RUN keytool -import -trustcacerts -cacerts -storepass changeit -noprompt -alias tsproxy -file /app/data/ts.crt

ENTRYPOINT ["java"]