## Contents
* Testing and Benchmarking
* go tool pprof
* go tool trace
	- how to read the views
	- tagging sections
* go memory analysis
	- stack and heap
	- escape analysis
* go slices
	- how do slices work internally. allocation and reuse.
* concurrency
	- (#)[sync pools]
	- sync once and lazy initializations
* specific optimizations
	- string vs buffer
	- heavy work in mutexes (https://commandercoriander.net/blog/2018/04/10/dont-lock-around-io/)
	- buffered vs unbuffered output
	- use int keys instead of string keys

## Testing

## Benchmarking

## Profiling

## Tracing

## sync.Pools
Pool's purpose is to cache allocated but unused items for later reuse, relieving pressure on the garbage collector. That is, it makes it easy to build efficient, thread-safe free lists. However, it is not suitable for all free lists.

A Pool is a set of temporary objects that may be individually saved and retrieved.

Any item stored in the Pool may be removed automatically at any time without notification. If the Pool holds the only reference when this happens, the item might be deallocated.

A Pool is safe for use by multiple goroutines simultaneously.

An appropriate use of a Pool is to manage a group of temporary items silently shared among and potentially reused by concurrent independent clients of a package. Pool provides a way to amortize allocation overhead across many clients.

An example of good use of a Pool is in the fmt package, which maintains a dynamically-sized store of temporary output buffers. The store scales under load (when many goroutines are actively printing) and shrinks when quiescent.

```
// 1_test.go
package main

import (
	"bytes"
	"testing"
)

func Benchmark_f1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f1()
	}
}

func f1() {
	s := &bytes.Buffer{}
	s.Write([]byte("dirty"))

	return
}

```

```
$ go test -bench=f1 -benchmem
Benchmark_f1-8   	30000000	        43.5 ns/op	      64 B/op	       1 allocs/op
```

```
// 2_test.go
package main

import (
	"bytes"
	"sync"
	"testing"
)

var pool2 = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func Benchmark_f2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f2()
	}
}

func f2() {
	// When getting from a Pool, you need to cast
	s := pool2.Get().(*bytes.Buffer)
	// We write to the object
	s.Write([]byte("dirty"))
	// Then put it back
	pool2.Put(s)

	return
}

```

```
$ go test -bench=f2 -benchmem
Benchmark_f2-8   	50000000	        38.2 ns/op	      14 B/op	       0 allocs/op
```

```Opt tip: Use sync.Pool is reduced your memory allocation pressure.```

### Exercise: sync.Pool
A type of data (book) needs to be written to a json file. An ISBN number is added to new book ({title, author}) and written out to a file.  Use sync.Pool to reduce allocations prior to writing.
See book1_test.go and book2_test.go

## sync.Once for Lazy Initialization

When programs have costly resources being loaded, it helps to do that only once.

In version 1 of our code, we have a template that needs to be parsed.  This example template is currently being read from memory, but there are usually many templates and they are read from the file system which can be very slow.

```code/sync-once```

In the first naive example, we load the template each time as and when required.  This is useful in that the template is loaded only when it is needed.

```
// 1.go
var t *template.Template

func f() {
	t = template.Must(template.New("").Parse(s))
	_ = t

	// do task with template
}

func main() {
	for i := 0; i < 10000; i++ {
		f()
	}
}

```

The time taken for this is about 0.637 seconds.  Can we improve on this?  
```
$ time go run 1.go

real	0m0.637s
user	0m0.712s
sys	0m0.346s
```

In version 1, we are re-parsing the template each time, which is unnecessary.  In the second version, we load the template only once at the beginning of the program.

```
// 2.go
func main() {
	// costs time at load and maybe unused
	t = template.Must(template.New("").Parse(s))
	_ = t

	for i := 0; i < 10000; i++ {
		f()
	}
}
```

This works well, but doing all our initialization at the very beginning will slow down the program's start.  It's often the case that there are many templates but all of them aren't needed or used immediately.  This is not preferred when there are multiple copies running in Kubernetes pods and we expect scaling to be very quick

```
time go run 2.go

real	0m0.365s
user	0m0.376s
sys	0m0.198s
```

In version 3 of our code, we use the sync.Once struct to ensure that the code is run once and only once and only at the instance it is first invoked, thus loading it 'lazily'.  

sync.Once is goroutine safe and will not be called simultaneously.

```
// 3.go
var t *template.Template
var o sync.Once

func g() {
	fmt.Println("within g()")
	t = template.Must(template.New("").Parse(s))
	_ = t
}

func f() {
	// only done once and when used
	o.Do(g)

	// do task with template

}

func main() {
	for i := 0; i < 10000; i++ {
		f()
	}
}
```

You can see that in our very simple program, the difference is not much.  But in typical production code, such changes could have a considerable impact.

```
time go run 3.go
within g()

real	0m0.380s
user	0m0.392s
sys	0m0.209s
```

```Opt tip: Consider lazily loading your resources using sync.Once at time of first use.```

## GOMAXPROCS
Discussion: for a program to be more efficient should you have more threads/goroutines or less?
Discussion: goroutines are kinda sorta similar to threads.  So why don't we just use threads instead of goroutines?

Threads typically take up more resources than goroutines - a minimum thread stack typically is upwards of 1MB.
A goroutine typically starts of at 2kh.  So, that's, at a vvery minimum, a reduction of 500x.  Anything else though?

A primary cost factor is contention.  Programs that has parallelism does not necessarily have higher performance because of greater contention for resources.

### What is GOMAXPROCS?
The GOMAXPROCS setting controls how many operating systems threads attempt to execute code simultaneously.  For example, if GOMAXPROCS is 4, then the program will only execute code on 4 operating system threads at once, even if there are 1000 goroutines. The limit does not count threads blocked in system calls such as I/O.

GOMAXPROCS can be set explicitly using the GOMAXPROCS environment variable or by calling runtime.GOMAXPROCS from within a program.

```code/gomaxprocs```

```
func main() {
	fmt.Println("runtime.NumCPU()=", runtime.NumCPU())
}
```

On my quad-core CPU it prints:
```
runtime.NumCPU()= 8
```

Why is it showing 8 for NumCPU for a quad-core machine? The Intel chips on my machine is hyperthreaded - for each processor core that is physically present, the operating system addresses two virtual (logical) cores and shares the workload between them when possible.

### What should be the value of GOMAXPROCS?
Upto Go 1.4, GOMAXPROCS defaulted to 1 because the 

The default setting of GOMAXPROCS in all Go releases [up to 1.4] is 1, because programs with frequent goroutine switches ran much slower when using multiple threads. It is much cheaper to switch between two goroutines in the same thread than to switch between two goroutines in different threads.

Goroutine scheduling affinity and other improvements to the scheduler have largely addressed the problem, by keeping goroutines that are concurrent but not parallel in the same thread.

For Go 1.5, the default setting of GOMAXPROCS to the number of CPUs available, as determined by runtime.NumCPU.

### Running with different GOMAXPROCS

```
GOMAXPROCS=1 go run mergesort.go v1 & go tool trace v1.trace
```

![GOMAXPROCS=1](./images/gomaxprocs/gomaxprocs-1.png)

```
GOMAXPROCS=8 go run mergesort.go v1 & go tool trace v1.trace
```

![GOMAXPROCS=8](./images/gomaxprocs/gomaxprocs-8.png)

```
GOMAXPROCS=18 go run mergesort.go v1 & go tool trace v1.trace
```

![GOMAXPROCS=18](./images/gomaxprocs/gomaxprocs-18.png)

The number is the max possible and it is not required that the Go runtime create as many logical processors as you have specified.  

### Exercise
```gocode/gomaxprocs```

Run the following and see the differences in the trace.

```
GOMAXPROCS=1 go run mergesort.go v2 && go tool trace v2.trace
GOMAXPROCS=8 go run mergesort.go v2 && go tool trace v2.trace
GOMAXPROCS=18 go run mergesort.go v2 && go tool trace v2.trace

GOMAXPROCS=1 go run mergesort.go v3 && go tool trace v3.trace
GOMAXPROCS=8 go run mergesort.go v3 && go tool trace v3.trace
GOMAXPROCS=18 go run mergesort.go v3 && go tool trace v3.trace
```

```Opt Tip: Do not assume that increasing the number of GOMAXPROCS always improves speed.```

## GOGC

The GOGC variable sets the initial garbage collection target percentage. A collection is triggered when the ratio of freshly allocated data to live data remaining after the previous collection reaches this percentage. The default is GOGC=100. Setting GOGC=off disables the garbage collector entirely. The runtime/debug package's SetGCPercent function allows changing this percentage at run time. 

GOGC controls the aggressiveness of the garbage collector.

Setting this value higher, say GOGC=200, will delay the start of a garbage collection cycle until the live heap has grown to 200% of the previous size. Setting the value lower, say GOGC=20 will cause the garbage collector to be triggered more often as less new data can be allocated on the heap before triggering a collection.

With the introduction of the low latency collector in Go 1.5, phrases like “trigger a garbage collection cycle” become more fluid, but the underlying message that values of GOGC greater than 100 mean the garbage collector will run less often, and for values of GOGC less than 100, more often


### Exercise
```gocode/gogc```

Run the following and see the differences in the trace for heap and GC.

```
GOGC=off go run mergesort.go v1 & go tool trace v1.trace
GOGC=50 go run mergesort.go v1 & go tool trace v1.trace
GOGC=100 go run mergesort.go v1 & go tool trace v1.trace
GOGC=200 go run mergesort.go v1 & go tool trace v1.trace
```

GOGC=off
![GOGC=off](./images/gogc/gogc-off.png)

GOGC=50
![GOGC=50](./images/gogc/gogc-50.png)

GOGC=100
![GOGC=100](./images/gogc/gogc-100.png)

GOGC=200
![GOGC=200](./images/gogc/gogc-200.png)

#References
* (https://www.dotconferences.com/2019/03/daniel-marti-optimizing-go-code-without-a-blindfold)[Daniel Marti's talk - Optimizing Go Code without a Blindfold]
* (https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html)[dave cheney high performance workshop]
* (https://github.com/davecheney/high-performance-go-workshop)[github - dave cheney high performance workshop]
* (https://commandercoriander.net/blog/2018/04/10/dont-lock-around-io/)[don't lock around io]
* (https://blog.gopheracademy.com/advent-2017/go-execution-tracer/)[advent 2017 - go execution tracer]
* (https://docs.google.com/document/u/1/d/1FP5apqzBgr7ahCCgFO-yoVhk4YZrNIDNf9RybngBc14/pub)[execution tracer design doc]
* https://www.alexedwards.net/blog/an-overview-of-go-tooling
* (https://www.alexedwards.net/blog/configuring-sqldb)[configuring sqldb for better performance]
* (https://www.alexedwards.net/blog/how-to-rate-limit-http-requests)[rate limit http requests]
* https://www.alexedwards.net/blog/understanding-mutexes
* https://stackimpact.com/docs/go-performance-tuning/
* https://stackimpact.com/blog/practical-golang-benchmarks/
* https://www.ardanlabs.com/blog/2017/06/design-philosophy-on-data-and-semantics.html
* https://github.com/ardanlabs/gotraining
* http://www.doxsey.net/blog/go-and-assembly
* https://medium.com/observability/debugging-latency-in-go-1-11-9f97a7910d68
* https://rakyll.org/profiler-labels/
* https://stackoverflow.com/questions/45027236/what-differentiates-exception-frames-from-other-data-on-the-return-stack
* https://www.infoq.com/presentations/self-heal-scalable-system
* https://dave.cheney.net/paste/clear-is-better-than-clever.pdf
* https://golang.org/pkg/sync/#Pool, https://dev.to/hsatac/syncpool-34pd
* http://dominik.honnef.co/go-tip/2014-01-10/#syncpool
* https://www.quora.com/In-C-what-does-buffering-I-O-or-buffered-I-O-mean
* https://stackoverflow.com/questions/1450551/buffered-vs-unbuffered-io
* http://www.agardner.me/golang/garbage/collection/gc/escape/analysis/2015/10/18/go-escape-analysis.html
* (https://docs.google.com/presentation/d/e/2PACX-1vTxoBN41dYFB8aV8c0SDET3B2htsAavXPAwR-CMyfT2LfARR2KjOt8EPIU1zn8ceSuxrL8BmkOqqL_c/pub?start=false&loop=false&delayms=3000&slide=id.g524654fd95_0_117)[Performance Optimization Sins - Aliaksandar Valialkin]
* https://blog.gopheracademy.com/advent-2018/postmortem-debugging-delve/
* https://github.com/golang/go/wiki/DesignDocuments
* (https://docs.google.com/document/d/1nr-TQHw_er6GOQRsF6T43GGhFDelrAP0NqSS_00RgZQ/edit)[Go execution modes]
* https://rakyll.org/profiler-labels/
* https://rakyll.org/pprof-ui/
* https://medium.com/@blanchon.vincent/go-should-i-use-a-pointer-instead-of-a-copy-of-my-struct-44b43b104963
* (https://www.youtube.com/watch?v=b0o-xeEoug0)[Performance tuning Go in GCP]
* https://medium.com/observability/want-to-debug-latency-7aa48ecbe8f7
* https://medium.com/dm03514-tech-blog/sre-debugging-simple-memory-leaks-in-go-e0a9e6d63d4d
* https://www.ardanlabs.com/blog/2013/07/understanding-type-in-go.html
* https://www.geeksforgeeks.org/structure-member-alignment-padding-and-data-packing/
* https://developers.redhat.com/blog/2016/06/01/how-to-avoid-wasting-megabytes-of-memory-a-few-bytes-at-a-time/
* https://go101.org/article/memory-layout.html
* https://dave.cheney.net/2015/10/09/padding-is-hard
* http://www.catb.org/esr/structure-packing/
* (https://scvalex.net/posts/29/)[Escape Analysis in Go]
* https://www.ardanlabs.com/blog/2018/01/escape-analysis-flaws.html
* https://www.ardanlabs.com/blog/2017/05/language-mechanics-on-stacks-and-pointers.html
* https://segment.com/blog/allocation-efficiency-in-high-performance-go-services/
* https://godoc.org/golang.org/x/perf/cmd/benchstat
* https://www.dotconferences.com/2019/03/daniel-marti-optimizing-go-code-without-a-blindfold
* https://www.youtube.com/watch?v=jiXnzkAzy30
* (https://gist.github.com/arsham/bbc93990d8e5c9b54128a3d88901ab90)[go cpu mem profiling benchmarks gist]
* https://hashrocket.com/blog/posts/go-performance-observations
* https://www.ardanlabs.com/blog/2017/05/language-mechanics-on-escape-analysis.html
* https://dave.cheney.net/2014/06/07/five-things-that-make-go-fast 
* https://stackoverflow.com/questions/2113751/sizeof-struct-in-go
* https://stackoverflow.com/questions/31496804/how-to-get-the-size-of-struct-and-its-contents-in-bytes-in-golang?rq=1
* https://github.com/campoy/go-tooling-workshop/tree/master/3-dynamic-analysis
* https://blog.usejournal.com/why-you-should-like-sync-pool-2c7960c023ba
* (https://rakyll.org/scheduler/)[work stealing scheduler]
* https://morsmachine.dk/go-scheduler
* https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part1.html
* https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part2.html
* https://www.ardanlabs.com/blog/2018/12/scheduling-in-go-part3.html
* https://www.welcometothejungle.co/fr/articles/languages-software-go-elixir
* https://eng.uber.com/optimizing-m3/
* https://medium.com/@fzambia/bisecting-go-performance-degradation-4d4a7ee83a63
* https://golang.org/doc/diagnostics.html
* http://jesseszwedko.com/gsp-go-debugging/#slide1
* https://fntlnz.wtf/post/gopostmortem/
* https://dave.cheney.net/2013/10/15/how-does-the-go-build-command-work
* https://medium.freecodecamp.org/how-i-investigated-memory-leaks-in-go-using-pprof-on-a-large-codebase-4bec4325e192
* https://medium.com/@cep21/using-go-1-10-new-trace-features-to-debug-an-integration-test-1dc39e4e812d
* https://medium.com/golangspec/goroutine-leak-400063aef468
* https://medium.com/@val_deleplace/go-code-refactoring-the-23x-performance-hunt-156746b522f7
* https://medium.com/@teivah/good-code-vs-bad-code-in-golang-84cb3c5da49d
* https://matoski.com/article/golang-profiling-flamegraphs/
* https://dzone.com/articles/so-you-wanna-go-fast
* https://www.slideshare.net/BadooDev/profiling-and-optimizing-go-programs
* https://about.sourcegraph.com/go/an-introduction-to-go-tool-trace-rhys-hiltner
* https://speakerdeck.com/rhysh/an-introduction-to-go-tool-trace
* https://stackimpact.com/blog/go-profiler-internals/
* https://syslog.ravelin.com/go-and-memory-layout-6ef30c730d51
* https://github.com/golang/go/wiki/Performance
* https://blog.golang.org/ismmkeynote
* https://making.pusher.com/golangs-real-time-gc-in-theory-and-practice/
* https://pusher.com/sessions/meetup/the-realtime-guild/golangs-realtime-garbage-collector
* https://blog.cloudflare.com/go-dont-collect-my-garbage/
* https://syslog.ravelin.com/further-dangers-of-large-heaps-in-go-7a267b57d487
* https://www.akshaydeo.com/blog/2017/12/23/How-did-I-improve-latency-by-700-percent-using-syncPool/
* (https://docs.google.com/document/d/1At2Ls5_fhJQ59kDK2DFVhFu3g5mATSXqqV5QrxinasI/edit)[Go 1.5 GOMAXPROCS default document]
* https://dave.cheney.net/2015/11/29/a-whirlwind-tour-of-gos-runtime-environment-variables
* (https://www.youtube.com/watch?v=ZMZpH4yT7M0)[https://engineers.sg/video/understanding-allocations-the-stack-and-the-heap-gophercon-sg-2019--3371]
* (https://blog.golang.org/ismmkeynote)[Getting to Go's Garbage Collector]
