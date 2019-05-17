go test -run=. -bench=. -cpuprofile=cpu.out -benchmem -memprofile=mem.out -trace trace.out
go tool pprof -pdf $FILENAME.test cpu.out > cpu.pdf && open cpu.pdf
go tool pprof -pdf --alloc_space $FILENAME.test mem.out > alloc_space.pdf && open alloc_space.pdf
go tool pprof -pdf --alloc_objects $FILENAME.test mem.out > alloc_objects.pdf && open alloc_objects.pdf
go tool pprof -pdf --inuse_space $FILENAME.test mem.out > inuse_space.pdf && open inuse_space.pdf
go tool pprof -pdf --inuse_objects $FILENAME.test mem.out > inuse_objects.pdf && open inuse_objects.pdf
go tool trace trace.out

go-torch $FILENAME.test cpu.out -f ${FILENAME}_cpu.svg && open ${FILENAME}_cpu.svg
go-torch --alloc_objects $FILENAME.test mem.out -f ${FILENAME}_alloc_obj.svg && open ${FILENAME}_alloc_obj.svg
go-torch --alloc_space $FILENAME.test mem.out -f ${FILENAME}_alloc_space.svg && open ${FILENAME}_alloc_space.svg
go-torch --inuse_objects $FILENAME.test mem.out -f ${FILENAME}_inuse_obj.svg && open ${FILENAME}_inuse_obj.svg
go-torch --inuse_space $FILENAME.test mem.out -f ${FILENAME}_inuse_space.svg && open ${FILENAME}_inuse_space.svg

# For live data

go-torch -u http://localhost:8080 --seconds 32 -f ${FILENAME}_live.svg && open ${FILENAME}_live.svg

#

go tool pprof -cum cpu.out
go tool pprof -cum --alloc_space mem.out
go tool pprof -cum --alloc_objects mem.out
go tool pprof -cum --inuse_space mem.out
go tool pprof -cum --inuse_objects mem.out

#

go tool pprof $FILENAME.test cpu.out
# (pprof) list <func name>

#

rm alloc_space.pdf alloc_objects.pdf inuse_space.pdf inuse_objects.pdf cpu.out cpu.pdf mem.out $FILENAME.test ${FILENAME}_cpu.svg ${FILENAME}_alloc_obj.svg ${FILENAME}_alloc_space.svg ${FILENAME}_inuse_obj.svg ${FILENAME}_inuse_space.svg ${FILENAME}_live.svg trace.out

