```
go test -covermode=count -coverprofile=count.out fmt
go tool cover -html=count.out
```

For current folder:
```
go test -covermode=count -coverprofile=count.out
go tool cover -html=count.out
```
