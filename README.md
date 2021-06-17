## Installation Steps:
### /api/v1 (normal, restful, non-grpc)
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
5. Use your favourite client (Postman recommended)
| Endpoints                                                    | Features                                   | Sample payload (JSON)                                        |
| ------------------------------------------------------------ | ------------------------------------------ | ------------------------------------------------------------ |
| PUT ```https://localhost:5000/api/v1/pickups/<Pickup_ID>?key=secretkey&role=user``` | User Approve a pickup                      |                                                              |
| DELETE ```https://localhost:5000/api/v1/pickups/<Pickup_ID>?key=secretkey&role=user``` | User cancels a pickup                      |                                                              |
| GET ```https://localhost:5000/api/v1/pickups/3?key=secretkey&role=user``` | Show this user's pickup requests           |                                                              |
| GET ```https://localhost:5000/api/v1/pickups```              | Show all available pickups                 |                                                              |
| PUT ```https://localhost:5000/api/v1/pickups/3?key=secretkey&role=collector``` | collector accepts/cancel a pickup          | {"pickup_id": "0487d326-9947-4d00-b3b6-a576cdafc506", "collector_id": ""} |
| GET ```https://localhost:5000/api/v1/pickups/3?key=secretkey&role=collector``` | collector view currently attending pickups |                                                              |
| POST ```https://localhost:5000/api/v1/pickups/3?key=secretkey&role=user``` | user request for a pickup                  | {"lat": 1.2411412, "lng": 102.78319552131565, "address": "Thor's avenue St 32, Milwaukee TN", "created_by": "14222"} |



### /api/v2 (grpc-server)
1. Launch a new terminal (from the root project folder)
2. Go to ```server/v2```
```shell
cd server/v2/server
go run main.go
```
3. curl (or postman)
```shell
curl http://localhost:8090/api/v2/pickups
```

### RecycleBinDetails

| GET ``` https://localhost:5000/api/v1/recyclebindetails/{userId}``` | Get specific user recyclebin feedback via their UserID |
| ------------------------------------------------------------ | ------------------------------------------------------ |
| GET ```https://localhost:5000/api/v1/recyclebindetails/NIL``` | Get all physical bin inventory i.e. where userId='NIL' |
| POST ```https://localhost:5000/api/v1/recyclebindetails/NIL``` | add user feedback entry to DB                          |
