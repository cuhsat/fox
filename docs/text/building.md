# Building
To build a full-featured version:
```console
go build -o fox main.go
```

## Custom-builds
To build a `minimal` version, stripped off terminal **UI** and **AI** support:
```console
go build -o fox -tags minimal main.go
```
