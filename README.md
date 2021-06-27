## Installation Steps:
### Setup DB
1. Launch Mysql workbench
2. Logon with default root account on default port 3306
3. Run the following sql statement to create a schema named ```your_db```

```sql
CREATE SCHEMA IF NOT EXISTS your_db;
```
4. Go to Server > Data Import. Import sql dumps located in:
```
...<ProjectRootFolder>/your_db
```
5. if steps 4 failed, manually run the following sql script

```sql
Use your_db;

CREATE TABLE `pickups` (
  `id` varchar(40) NOT NULL,
  `coord` point NOT NULL,
  `address` varchar(150) NOT NULL,
  `created_by` varchar(40) NOT NULL,
  `attend_by` varchar(40) NOT NULL DEFAULT '',
  `completed` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idrecycle_UNIQUE` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
```

5. Create a recyle table

```sql
CREATE TABLE `recyclebinsdetails` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `BinID` varchar(10) DEFAULT NULL,
  `BinType` varchar(10) DEFAULT NULL,
  `BinLocationLat` double DEFAULT NULL,
  `BinLocation﻿Long` double DEFAULT NULL,
  `BinAddress` varchar(60) DEFAULT NULL,
  `Postcode` varchar(10) DEFAULT NULL,
  `UserID` varchar(30) DEFAULT NULL,
  `FBOptions` varchar(20) DEFAULT NULL,
  `ColorCode` varchar(20) DEFAULT NULL,
  `Remarks` varchar(150) DEFAULT NULL,
  `BinStatusUpdate` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
```
6. Create a ```user1``` in mysql workbench and grant all privilege to ```your_db```
The connection goes like this:
```sh
sql.Open("mysql", "user1:password@tcp(127.0.0.1:3306)/your_db")
```
### /api/v1 (normal, restful, non-grpc)
1. Clone repo (or unzip if you are downloading via zip file)

```shell
git clone https://github.com/heartziq/recycle-lah.git
```

2. Go into the root of the project folder, create a folder named ```cert``` inside ```server/``` . Generate **own** secret keys and server certificate
```shell
cd recycle-lah
mkdir server/cert
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
