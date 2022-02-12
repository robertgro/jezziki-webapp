@echo off
echo .
echo "OpenSSL Gen cert (Host is localhost)"
set /p targetIP="Please specify a target SAN IP:"
openssl req -x509 -sha256 -newkey rsa:4096 -keyout lhkey.pem -passout pass:123456 -out lhcrt.pem -days 365 -nodes -subj /CN=localhost/O=home/C=US/emailAddress=me@mail.internal -addext "subjectAltName = DNS:localhost,IP:%targetIP%,email:me@mail.internal" -addext keyUsage=digitalSignature -addext extendedKeyUsage=serverAuth