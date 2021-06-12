## Installation Steps:
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

### RecycleBinDetails

| GET ``` https://localhost:5000/api/v1/recyclebindetails/{userId}``` | Get specific user recyclebin feedback    |
| ------------------------------------------------------------ | ---------------------------------------- |
| GET ```https://localhost:5000/api/v1/recyclebindetails/NIL``` | Get All feedback i.e. where userId='NIL' |
| POST ```https://localhost:5000/api/v1/recyclebindetails/NIL``` | Add new recycleBin entry to DB           |

### Pickups

| POST ```https://localhost:5000/api/v1/pickups/3?key=secretkey&role=user``` | user request for a pickup                  |
| ------------------------------------------------------------ | ------------------------------------------ |
| PUT ```https://localhost:5000/api/v1/pickups/14af6bf8-d761-48ef-a5fd-ac925a1c94c3?key=secretkey&role=user``` | User Approve a pickup                      |
| DELETE ```https://localhost:5000/api/v1/pickups/412rawf3?key=secretkey&role=user``` | User cancel a pickup                       |
| GET ```https://localhost:5000/api/v1/pickups/3?key=secretkey&role=user``` | Show this user's pickup requests           |
| GET ```https://localhost:5000/api/v1/pickups```              | Show all available pickups                 |
| PUT ```https://localhost:5000/api/v1/pickups/3?key=secretkey&role=collector``` | collector accepts/cancel a pickup          |
| GET ```https://localhost:5000/api/v1/pickups/3?key=secretkey&role=collector``` | collector view currently attending pickups |
|                                                              |                                            |

