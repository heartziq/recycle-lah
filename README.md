1. Clone project

```shell
git clone https://github.com/heartziq/recycle-lah.git
```

2. Generate own secret keys and certificate
```shell
cd recycle-lah
./runthis.sh
```

3. Verify all dependencies are install
```shell
go mod tidy
```

4. Run server
```shell
cd server
go run server.go
```