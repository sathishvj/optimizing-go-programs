#Contents
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
	- sync pools
	- sync once and lazy initializations
* specific optimizations
	- string vs buffer
	- heavy work in mutexes (https://commandercoriander.net/blog/2018/04/10/dont-lock-around-io/)
	- buffered vs unbuffered output
	- use int keys instead of string keys

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
