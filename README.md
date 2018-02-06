# web-archiver

web-archiver is a golang service which persists HTML data in Microsoft Azure Blob Storage.

# Docker setup

```
docker build -t go-server .
```

```
docker run -p 8000 --env-file=conf.env <image>
```

# conf.env

Secrets can be stored in a conf.env file. 

```
MYSERVER_ACCOUNTNAME=AccountName
MYSERVER_ACCOUNTKEY=AccountKey
```