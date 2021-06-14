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

| GET ``` https://localhost:5000/api/v1/recyclebindetails/{userId}``` | Get specific user recyclebin feedback via their UserID |
| ------------------------------------------------------------ | ------------------------------------------------------ |
| GET ```https://localhost:5000/api/v1/recyclebindetails/NIL``` | Get all physical bin inventory i.e. where userId='NIL' |
| POST ```https://localhost:5000/api/v1/recyclebindetails/NIL``` | add user feedback entry to DB                          |

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

### Reward Points
| GET ```https://localhost:5000/api/v1/rewards/{userId}?key=secretkey``` | Get user's reward points           |
| PUT ```https://localhost:5000/api/v1/rewards/{userId}?key=secretkey``` | Update user's reward points        |

### Users
| POST ```https://localhost:5000/api/v1/users/{userId}?key=secretkey```   | add new user or new collector       |
| PUT ```https://localhost:5000/api/v1/users/{userId}?key=secretkey```    | update account particulars          |
| GET ```https://localhost:5000/api/v1/users/{userId}?key=secretkey```    | verify login                        |
| DELETE ```https://localhost:5000/api/v1/users/{userId}?key=secretkey``` | delete user account                 |
