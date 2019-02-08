# shortsh-backend

This is the backend for the url shortener running at https://short.sh/.

The setup is simple, you have two choices:
1. Download a precompiled binary from the releases tab, tagged with the current version
2. Download the source and build yourself

## 1. Precompiled version
Steps:
- Download binary version for your platform
- Create a _config folder and copy the sample file from the _config folder or copy it from here
```json
{
  "mysql_dsn": "user:pass@/db?charset=utf8",
  "allow_origins": "http://localhost:3000"
}
```
The `mysql_dsn` contains the access for the mysql database. The `allow_origins` property defines the origins from which we will accept requests from the webapp.
If you have your webapp hosted at `https://localhost:3000` then the origin would be like in the example.

## 2. Build from source
If you somehow like to build your own binaries (who doesn't?) then you can download the current source and follow the instructions below.
Steps:
- Download the source from github
- Run `go get`
- Run `go build`

We use the `go mod` in this project so you need nearly no configurations and you are ready to go in a while.
 