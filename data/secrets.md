
* Generate the server key
```
openssl genrsa -out server.key 2048
```

* Generate the server certificate
```
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```
