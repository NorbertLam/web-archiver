# go-azure-server

go-azure-server starts up a web server which handles storing html data within a Microsoft Azure Blob storage service.

# Docker setup
To read keys and private info, environment variables are stored and read off conf.env file. 

```
docker build -t go-server .
```

```
docker run -p 8000 --env-file=conf.env <image>
```