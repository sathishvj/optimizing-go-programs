## Contents
* Testing and Benchmarking
* go tool pprof
* go tool trace
	- how to read the views
	- tagging sections
* environment variables
	- (#GOMAXPROCS)[GOMAXPROCS]
	- (#GOGC)[GOGC]
* go memory analysis
	- (#Stack-and-Heap)[stack and heap]
	- (#Escape-Analysis)[escape analysis]
* concurrency
	- (#sync.Pools)[sync Pools]
	- (#sync.Once-for-Lazy-Initialization)[sync once and lazy initializations]
* go slices
	- how do slices work internally. allocation and reuse.
* specific optimizations
	- string vs buffer
	- heavy work in mutexes (https://commandercoriander.net/blog/2018/04/10/dont-lock-around-io/)
	- buffered vs unbuffered output
	- use int keys instead of string keys

* Performance Tuning Patterns

## Testing

Unit testing is important enough to be a standard library.

```code/testing```

```
func BeginsWith(s, pat string) bool {
	return strings.HasPrefix(s, pat)
}

func Test_BeginsWith(t *testing.T) {
	tc := []struct {
		s, pat string
		exp    bool
	}{
		{"GoLang", "Go", true},
		{"GoLang", "Java", false},
		{"GoLang is awesome", "awe", false},
		{"awesome is GoLang. - Yoda", "awe", true},
	}

	for _, tt := range tc {
		if BeginsWith(tt.s, tt.pat) != tt.exp {
			t.Fail()
		}
	}
}
```

```
$ go test -v
=== RUN   Test_BeginsWith
--- PASS: Test_BeginsWith (0.00s)
PASS
```

Testing validates your code.  It checks for correctness.

```Opt tip: unit testing first, always.```

p.s. When you run benchmarks, tests are run first.

## Benchmarking

Benchmarking checks for optimization.

```code/testing```

```
func Benchmark_BeginsWith(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BeginsWith("GoLang", "Go")
	}
}
```

```
$ go test -v -bench=. -benchmem
=== RUN   Test_BeginsWith
--- PASS: Test_BeginsWith (0.00s)
goos: darwin
goarch: amd64
Benchmark_BeginsWith-8   	500000000	         3.69 ns/op	       0 B/op	       0 allocs/op
PASS
```

Benchmarking functions don't always care about the result (that is checked by unit testing).  However, the speed/allocations/blocking of a function could be dependent on the inputs - so test different inputs. 

## Coverage

The Go tooling also gives you automatic coverage results.  Less code is faster code.  Tested and covered code is more reliable code.

```code/cover```

```
go test -covermode=count -coverprofile=count.out fmt
go tool cover -html=count.out
```

For current folder:
```
go test -covermode=count -coverprofile=count.out
go tool cover -html=count.out
```


## Profiling

Package pprof writes runtime profiling data in the format expected by the pprof visualization tool.

The first step to profiling a Go program is to enable profiling. Support for profiling benchmarks built with the standard testing package is built into go test.

```
func isGopher(email string) (string, bool) {
	re := regexp.MustCompile("^([[:alpha:]]+)@golang.org$")
	match := re.FindStringSubmatch(email)
	if len(match) == 2 {
		return match[1], true
	}
	return "", false
}

func Benchmark_isGopher(b *testing.B) {

	tcs := []struct {
		in    string
		exp   bool
		expId string
	}{
		{
			"a@golang.org",
			true,
			"a",
		},
	}

	for i := 0; i < b.N; i++ {
		isGopher(tcs[0].in)
	}
}
```

```
go test -bench=. -cpuprofile=cpu.pprof

go tool pprof cpu.pprof

go-torch --binaryname web.test -b cpu.pprof
```

More recently (1.10?), pprof got its own UI.

```
$ go get github.com/google/pprof
```

The tool launches a web UI if -http flag is provided. For example, in order to launch the UI with an existing profile data, run the following command:


```
pprof -http=:8080 cpu.pprof
```

There is also a standard HTTP interface to profiling data. Adding the following line will install handlers under the /debug/pprof/ URL to download live profiles:

```
import _ "net/http/pprof"
See the net/http/pprof package for more details.
```

## M, P, G

OS Layout

![OS Layout](./images/tracing/1-OS-process-and-its-threads.png)

Goroutines on a Thread

![Goroutines on a Thread](./images/tracing/2-goroutines-on-a-thread.png)

Goroutines on Blocking Thread

![Goroutines on Blocking Thread](./images/tracing/3-goroutines-on-a-blocking-thread.png)

Concurrency and Parallelism

![Concurrency and Parallelism](./images/tracing/4-concurrency-and-parallelism.png)

## Tracing

https://blog.gopheracademy.com/advent-2017/go-execution-tracer/

Ever wondered how are your goroutines being scheduled by the go runtime? Ever tried to understand why adding concurrency to your program has not given it better performance? The go execution tracer can help answer these and other questions to help you diagnose performance issues, e.g, latency, contention and poor parallelization.

Data is collected by the tracer without any kind of aggregation or sampling. In some busy applications this may result in a large file.

While the CPU profiler does a nice job to telling you what function is spending most CPU time, it does not help you figure out what is preventing a goroutine from running or how are the goroutines being scheduled on the available OS threads. That’s precisely where the tracer really shines.

### Ways to get a Trace

* Using the runtime/trace pkg
This involved calling trace.Start and trace.Stop and was covered in our “Hello, Tracing” example.

* Using -trace=<file> test flag
This is useful to collect trace information about code being tested and the test itself.

```code/tracing```
```
go test -trace=a.out && go tool trace a.out
```

* Using debug/pprof/trace handler
This is the best method to collect tracing from a running web application.

### View Trace

```
go tool trace trace_file.out
```

![View Trace](./images/tracing/view-trace.png)

1. Timeline
Shows the time during the execution and the units of time may change depending on the navigation. One can navigate the timeline by using keyboard shortcuts (WASD keys, just like video games).
2. Heap
Shows memory allocations during the execution, this can be really useful to find memory leaks and to check how much memory the garbage collection is being able to free at each run.
3. Goroutines
Shows how many goroutines are running and how many are runnable (waiting to be scheduled) at each point in time. A high number of runnable goroutines may indicate scheduling contention, e.g, when the program creates too many goroutines and is causing the scheduler to work too hard.
4. OS Threads
Shows how many OS threads are being used and how many are blocked by syscalls.
5. Virtual Processors
Shows a line for each virtual processor. The number of virtual processors is controlled by the GOMAXPROCS environment variable (defaulting to the number of cores).
6. Goroutines and events
Displays where/what goroutine is running on each virtual processor. Lines connecting goroutines represent events. In the example image, we can see that the goroutine “G1 runtime.main” spawned two different goroutines: G6 and G5 (the former is the goroutine responsible for collecting the trace data and the latter is the one we started using the “go” keyword).
A second row per processor may show additional events such as syscalls and runtime events. This also includes some work that the goroutine does on behalf of the runtime (e.g assisting the garbage collector).


### View Goroutine

![View Trace](./images/tracing/view-goroutine.png)

This information includes:

* Its “name” (Title)
* When it started (Start)
* Its duration (Wall Duration)
* The stack trace when it started
* The stack trace when it finished
* Events generated by this goroutine

### Tracing Example

```code/tracing```

```
func main() {
	f, _ := os.OpenFile(version+".trace", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	trace.Start(f)
	defer trace.Stop()

	mergesortv1(s)

}

```

```
go run mergesort.go v1 && go tool trace v1.trace
```

### Tracing Conclusion
The tracer is a powerful tool for debugging concurrency issues, e.g, contentions and logical races. But it does not solve all problems: it is not the best tool available to track down what piece of code is spending most CPU time or allocations. The go tool pprof is better suited for these use cases.

The tool really shines when you want to understand the behavior of a program over time and to know what each goroutine is doing when NOT running. Collecting traces may have some overhead and can generate a high amount of data to be inspected.

## GOMAXPROCS

Discussion: for a program to be more efficient should you have more threads/goroutines or less?

Discussion: goroutines are kinda sorta similar to threads.  So why don't we just use threads instead of goroutines?

Threads typically take up more resources than goroutines - a minimum thread stack typically is upwards of 1MB.
A goroutine typically starts of at 2kh.  So, that's, at a very minimum, a reduction of 500x.  Anything else though?

Context switching in Linux is about 1000ns while in go it is about 200ns - https://eli.thegreenplace.net/2018/measuring-context-switching-and-memory-overheads-for-linux-threads/

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

```Opt Tip: This helps you analyze your GC patterns but I can't find any posts that recommend this as a good performance tuning strategy.```


## Stack and Heap

Discussion: where is the stack memory shown in a trace diagram? 

ref: https://scvalex.net/posts/29/

### Stack Frame
ref: http://www.cs.uwm.edu/classes/cs315/Bacon/Lecture/HTML/ch10s07.html

The stack frame, also known as activation record is the collection of all data on the stack associated with one subprogram call.

The stack frame generally includes the following components:

* The return address
* Argument variables passed on the stack
* Local variables (in HLLs)
* Saved copies of any registers modified by the subprogram that need to be restored

 The Stack
 ---------

```
|      f()      |
|               |
+---------------+
|  func f(){    |  \
|       g()     |   } Stack frame of calling function f()
|  }            |  /
+---------------+
|  func g() {   |  \
|     a := 10   |   } Stack frame of called function: g()
|  }            |  /
+---------------+
================= // invalid below this
```

As the function call returns, the stack unwinds leaving previous stack frames invalid.

```
|      f()      |
|               |
+---------------+
|  func f(){    |  \
|       g()     |   } Stack frame of calling function f()
|  }            |  /
+---------------+
================= // invalid below this
|  func g() {   |  \
|     a := 10   |   } Stack frame of called function: g()
|               |  /
+---------------+
```

All local variables are no more accessible.  In C, returning a pointer to a local variable would cause a segmentation fault.

```
// online c editor - https://onlinegdb.com/HySykSJoE

#include <stdio.h>

int* f() {
    int a;
    a = 10;
    return &a;
}

void main()
{
    int* p = f();
    printf("p is: %x\n", p);   // p is 0
    printf("*p is: %d\n", *p); // segmentation fault

	// 
}
```

## Escape Analysis

In C, returning the reference of a local variable causes a segfault because that memory is no more valid.

```
// online c editor - https://onlinegdb.com/HySykSJoE

#include <stdio.h>

int* f() {
    int a;
    a = 10;
    return &a;
}

void main()
{
    int* p = f();
    printf("p is: %x\n", p);   // p is 0
    printf("*p is: %d\n", *p); // segmentation fault

	//
}
```

In Go, it is allowed to return the reference of a local variable.  

```
package main

import (
	"fmt"
)

func f() *int {
	x := 10
	return &x
}

func main() {
	fmt.Println(*f()) // prints 10
}
```

How is that possible?

From Effective Go: "Note that, unlike in C, it’s perfectly OK to return the address of a local variable; the storage associated with the variable survives after the function returns."

"When possible, the Go compilers will allocate variables that are local to a function in that function’s stack frame. However, if the compiler cannot prove that the variable is not referenced after the function returns, then the compiler must allocate the variable on the garbage-collected heap to avoid dangling pointer errors. In the current compilers, if a variable has its address taken, that variable is a candidate for allocation on the heap. However, a basic escape analysis recognizes some cases when such variables will not live past the return from the function and can reside on the stack."

*Can we figure out when variables escape to the heap?*

```
// go build -gcflags='-m' 1.go
// go build -gcflags='-m -l' 1.go to avoid inlining
// go build -gcflags='-m -l -m' 1.go for verbose comments.
```

```
func f() {
	var i = 5
	i++
	_ = i
}

func main() {
	f()
}
```

```
$ go build -gcflags='-m -l -m' 1.go
```

```
func f_returns() int {
	var i = 5
	i++
	return i
}

func main() {
	f_returns()
}
```

```
$ go build -gcflags='-m -l -m' 1.go
```

```
func f_returns_ptr() *int {
	var i = 5
	i++
	return &i
}

func main() {
	f_returns_ptr()
}
```

```
$ go build -gcflags='-m -l -m' 1.go
# command-line-arguments
./1.go:24:9: &i escapes to heap
./1.go:24:9: 	from ~r0 (return) at ./1.go:24:2
./1.go:22:6: moved to heap:
```

Once the variable is on the heap, there is pressure on the Garbage Collector.

Garbage collection is a convenient feature of Go - automatic memory management makes code cleaner and memory leaks less likely. However, GC also adds overhead as the program periodically needs to stop and collect unused objects. The Go compiler is smart enough to automatically decide whether a variable should be allocated on the heap, where it will later need to be garbage collected, or whether it can be allocated as part of the stack frame of the function which declared it. Stack-allocated variables, unlike heap-allocated variables, don’t incur any GC overhead because they’re destroyed when the rest of the stack frame is destroyed - when the function returns.

To perform escape analysis, Go builds a graph of function calls at compile time, and traces the flow of input arguments and return values.

However, if there are variables to be shared, it is appropriate for it to be on the heap.

```Opt tip: If you’ve profiled your program’s heap usage and need to reduce GC time, there may be some wins from moving frequently allocated variables off the heap. ```

See: https://segment.com/blog/allocation-efficiency-in-high-performance-go-services/
See: http://www.agardner.me/golang/garbage/collection/gc/escape/analysis/2015/10/18/go-escape-analysis.html

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

## Go Performance Patterns
When application performance is a critical requirement, the use of built-in or third-party packages and methods should be considered carefully. The cases when a compiler can optimize code automatically are limited. The Go Performance Patterns are benchmark- and practice-based recommendations for choosing the most efficient package, method or implementation technique.

Some points may not be applicable to a particular program; the actual performance optimization benefits depend almost entirely on the application logic and load.

### Parallelize CPU work
When the work can be parallelized without too much synchronization, taking advantage of all available cores can speed up execution linearly to the number of physical cores.

### Make multiple I/O operations asynchronous
Network and file I/O (e.g. a database query) is the most common bottleneck in I/O-bound applications. Making independent I/O operations asynchronous, i.e. running in parallel, can improve downstream latency. Use sync.WaitGroup to synchronize multiple operations.

### Avoid memory allocation in hot code
Object creation not only requires additional CPU cycles, but will also keep the garbage collector busy. It is a good practice to reuse objects whenever possible, especially in program hot spots. You can use sync.Pool for convenience. See also: Object Creation Benchmark

### Favor lock-free algorithms
Synchronization often leads to contention and race conditions. Avoiding mutexes whenever possible will have a positive impact on efficiency as well as latency. Lock-free alternatives to some common data structures are available (e.g. Circular buffers).

### Use read-only locks
The use of full locks for read-heavy synchronized variables will unnecessarily make reading goroutines wait. Use read-only locks to avoid it.

### Use buffered I/O
Disks operate in blocks of data. Accessing disk for every byte is inefficient; reading and writing bigger chunks of data greatly improves the speed. See also: File I/O Benchmark

### Use StringBuffer or StringBuilder instead of += operator
A new string is allocated on every assignment, which is inefficient and should be avoided. See also: String Concatenation Benchmark.

### Use compiled regular expressions for repeated matching
It is inefficient to compile the same regular expression before every matching. While obvious, it is often overlooked. See also: Regexp Benchmark.

### Preallocate slices
Go manages dynamically growing slices intelligently; it allocates twice as much memory every time the current capacity is reached. During re-allocation, the underlying array is copied to a new location. To avoid copying the memory and occupying garbage collection, preallocate the slice fully whenever possible. See also: Slice Appending Benchmark.

### Use Protocol Buffers or MessagePack instead of JSON and Gob
JSON and Gob use reflection, which is relatively slow due to the amount of work it does. Although Gob serialization and deserialization is comparably fast, though, and may be preferred as it does not require type generation. See also: Serialization Benchmark.

### Use int keys instead of string keys for maps
If the program relies heavily on maps, using int keys might be meaningful, if applicable. See also: Map Access Benchmark.

# References
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
