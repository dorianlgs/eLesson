# eLesson

## Install

```bash
go mod tidy
```

## Run Dev

```bash
go run . serve
```

## Build to publish

```bash
go generate ./...
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"
./eLesson serve
```

## Copy to Server
```bash
rsync -av /Users/doriangonzalez/Workspace/eLesson/eLesson root@192.168.100.175:/var/www/pocketbase

rsync -av /Users/doriangonzalez/Workspace/eLesson/pb_migrations/ root@192.168.100.175:/var/www/pb_migrations/
```

## Module creation
```bash
go mod init github.com/shujink0/eLesson
```

## Update All Go Modules
```bash
go get -u -t ./...
go mod tidy
```
