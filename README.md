# web-archiver

web-archiver starts up a web server which handles storing html data within a Microsoft Azure Blob storage service.

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