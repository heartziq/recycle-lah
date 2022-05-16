# test comment
openssl genrsa -out server/utility/secret.pem 2048

openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server/cert/key.pem -out server/cert/cert.pem
