```
go test -bench=. -cpuprofile=cpu.pprof

go tool pprof cpu.pprof

go-torch --binaryname web.test -b cpu.pprof

pprof -http=:8080 cpu.pprof
```
