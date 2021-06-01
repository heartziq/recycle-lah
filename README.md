1. Generate **own** server certificate and respective private keys

```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout cert/key.pem -out cert/cert.pem
```

2. From the root folder, run
```bash
cd server
go run server.go
```