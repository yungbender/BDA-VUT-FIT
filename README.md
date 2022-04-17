# DASH Network node scanner
Platform which scans whole DASH network and maps nodes and their status to the database.




# Requirements
## Base requirements
- Docker
- Docker-compose
- Internet connection
## Dev requirements
- Golang
- Postgres

# Description
Project contains 4 containers.
- nodes_db - PostgreSQL database for storing nodes info
- nodes_api - REST API for getting info about nodes
- nodes_crawler - Crawls whole DASH network and gets info about nodes
- nodes_pinger - Pings every stored node and keeps online/offline info

# Running platform
```
docker-compose up --build
```
Starts up every mentioned container which works together.

Crawler crawls whole DASH network on start and saves info to the DB.
This performs every 4 hours (cronjob). Crawling takes from 20 minutes to 1 hour (depends on network).

Pinger pings all saved nodes and gets their activity status. Attempts to connect to all offline nodes on start and does this every 10 minutes. Opened connection is performing PING/PONG exchange every 5 minutes.

Getting meaningful results is available after first Crawler run and then after repings of the Pinger.

# Local development
You can run each component on local machine by running.
```
go run main.go <pinger / crawler / api>
```
Every component expects credentials to the nodes database on enviroment variables:
- DB_HOST - host
- DB_PORT - port
- DB_NAME - db name
- DB_USER - db user
- DB_USER_PWD - db user pwd

# Getting results
For results you can remotely execute to the nodes database and query for results or use pre-programmed REST API calls.
```
$ curl localhost:8080/v1/status | jq

{
  "code": 200,
  "data": {
    "online": 5134,
    "offline": 16443,
    "unknown": 145
  }
}

```
```
$ curl localhost:8080/v1/useragents | jq

{
  "code": 200,
  "data": {
    "user_agents": [
      {
        "user_agent": "/Dash Core:0.17.0.3/",
        "percentage": 80.45,
        "count": 4280
      },
      {
        "user_agent": "/Binarium Core:0.12.9.1/",
        "percentage": 0.09,
        "count": 5
      },
      {
        "user_agent": "/Dash Core:0.17.0.2/",
        "percentage": 6.5,
        "count": 346
      },
      {
        "user_agent": "/Desire Core:0.12.2.1/",
        "percentage": 0.06,
        "count": 3
      },
      {
        "user_agent": "/Dash Core:0.13.1/",
        "percentage": 0.02,
        "count": 1
      },
      {
        "user_agent": "/LksCoin Core:3.2.0/",
        "percentage": 0.02,
        "count": 1
      },
      {
        "user_agent": "/Dash Core:0.17.0.3(zelcore)/",
        "percentage": 0.02,
        "count": 1
      },
      {
        "user_agent": "/Dash Core:0.13.3/",
        "percentage": 0.04,
        "count": 2
      },
      {
        "user_agent": "/LksCoin Core:3.0.0.2/",
        "percentage": 0.15,
        "count": 8
      },
      {
        "user_agent": "/Dash Core:0.14.0.3/",
        "percentage": 0.11,
        "count": 6
      },
      {
        "user_agent": "/Dash Core:0.14.0.2/",
        "percentage": 0.04,
        "count": 2
      },
      {
        "user_agent": "/Dash Core:0.14.0/",
        "percentage": 0.04,
        "count": 2
      },
      {
        "user_agent": "/LksCoin Core:3.1.0.1/",
        "percentage": 0.06,
        "count": 3
      },
      {
        "user_agent": "/Binarium Core:0.12.9.2/",
        "percentage": 0.06,
        "count": 3
      },
      {
        "user_agent": "/Dash Core:0.12.3.3/",
        "percentage": 0.04,
        "count": 2
      },
      {
        "user_agent": "/Dash Core:0.17.0.3(dashcore)/",
        "percentage": 0.11,
        "count": 6
      },
      {
        "user_agent": "/Dash Core:0.16.1/",
        "percentage": 0.02,
        "count": 1
      },
      {
        "user_agent": "/LksCoin Core:3.2.0.1/",
        "percentage": 0.36,
        "count": 19
      },
      {
        "user_agent": "/Desire Core:0.12.2.3/",
        "percentage": 0.02,
        "count": 1
      },
      {
        "user_agent": "/Dash Core:0.16.0.1/",
        "percentage": 0.08,
        "count": 4
      },
      {
        "user_agent": "/Dash Core:0.12.2.3/",
        "percentage": 0.09,
        "count": 5
      },
      {
        "user_agent": "/Lksc Core:3.3.0/",
        "percentage": 10.7,
        "count": 569
      },
      {
        "user_agent": "/Lks Core:3.0.1/",
        "percentage": 0.02,
        "count": 1
      },
      {
        "user_agent": "/Desire Core:0.12.2.2/",
        "percentage": 0.28,
        "count": 15
      },
      {
        "user_agent": "/Dash Core:0.17.0.3(bitcore)/",
        "percentage": 0.04,
        "count": 2
      },
      {
        "user_agent": "/HCC Core:0.17.0.3/",
        "percentage": 0.02,
        "count": 1
      },
      {
        "user_agent": "/Dash Core:0.16.1.1/",
        "percentage": 0.36,
        "count": 19
      },
      {
        "user_agent": "/Dash Core:0.14.0.1/",
        "percentage": 0.06,
        "count": 3
      },
      {
        "user_agent": "/Binarium Core:0.12.2.8/",
        "percentage": 0.09,
        "count": 5
      },
      {
        "user_agent": "/Dash Core:0.12.3.2/",
        "percentage": 0.02,
        "count": 1
      },
      {
        "user_agent": "/Dash Core:0.12.2.1(bitcore-sl)/",
        "percentage": 0.02,
        "count": 1
      },
      {
        "user_agent": "/LksCoin Core:3.0.0.4/",
        "percentage": 0.02,
        "count": 1
      },
      {
        "user_agent": "/Dash Core:0.17.0.2(dashcore)/",
        "percentage": 0.02,
        "count": 1
      }
    ],
    "total_count": 5320
  }
}
```
```
$ curl localhost:8080/v1/livenodes | jq

{
  "code": 200,
  "data": [
    {
      "ip": "85.209.241.92",
      "port": 9999
    },
    {
      "ip": "209.126.7.246",
      "port": 9400
    },
    {
      "ip": "23.163.0.49",
      "port": 9999
    },
    {
      "ip": "128.199.146.244",
      "port": 9999
    },
    <lots of IPs>
  ]
}
```