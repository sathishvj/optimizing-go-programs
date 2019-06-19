Let's compare what the compiled output of these two fairly similar programs are.

```
// a.go
3 func a(a []int) {
4	n := 6
5	_ = a[n]
6 }
```

```
// b.go
3 func b(b [5]int) {
4	n := len(b) - 1
5	_ = b[n]
6 }
```

```
$ go tool compile -S a.go > a.co
$ go tool compile -S b.go > b.co
$ vimdiff a.co b.co
```

```
"".a STEXT nosplit size=39 args=0x18 locals=0x8
	(a.go:3)	TEXT	"".a(SB), NOSPLIT|ABIInternal, $8-24
	(a.go:3)	SUBQ	$8, SP
	(a.go:3)	MOVQ	BP, (SP)
	(a.go:3)	LEAQ	(SP), BP
	(a.go:3)	FUNCDATA	$0, gclocals·1a65...
	(a.go:3)	FUNCDATA	$1, gclocals·69c1...
	(a.go:3)	FUNCDATA	$3, gclocals·33cd...
	(a.go:5)	PCDATA	$2, $0
	(a.go:5)	PCDATA	$0, $1
	(a.go:5)	MOVQ	"".a+24(SP), AX
	(a.go:5)	CMPQ	AX, $6
	(a.go:5)	JLS	32
	(a.go:6)	PCDATA	$2, $-2
	(a.go:6)	PCDATA	$0, $-2
	(a.go:6)	MOVQ	(SP), BP
	(a.go:6)	ADDQ	$8, SP
	(a.go:6)	RET
	(a.go:5)	PCDATA	$2, $0
	(a.go:5)	PCDATA	$0, $1
	(a.go:5)	CALL	runtime.panicindex(SB)
	(a.go:5)	UNDEF
	0x0000 48 83 ec 08 48 89 2c 24 48 8d 2c 24 48 8b 44 24  H...H.,$H.,$H.D$
	0x0010 18 48 83 f8 06 76 09 48 8b 2c 24 48 83 c4 08 c3  .H...v.H.,$H....
	0x0020 e8 00 00 00 00 0f 0b                             .......
	rel 33+4 t=8 runtime.panicindex+0
```

```
// b.co
"".b STEXT nosplit size=1 args=0x28 locals=0x0
	(b.go:3)	TEXT	"".b(SB), NOSPLIT|ABIInternal, $0-40
	(b.go:3)	FUNCDATA	$0, gclocals·33cd...
	(b.go:3)	FUNCDATA	$1, gclocals·33cd...
	(b.go:3)	FUNCDATA	$3, gclocals·33cd...
	(b.go:6)	RET
```

There seems to be way more happening in a.go than in b.go - about 20+ lines more, which seems surprising.

A little too much though.  That's probably because of optimizations by the compiler.  Let's remove those with the -N option.

```
$ go tool compile -S -N a.go > a.co
$ go tool compile -S -N b.go > b.co
$ vimdiff a.co b.co
```

```
"".a STEXT nosplit size=49 args=0x18 locals=0x10
	(a.go:3)	TEXT	"".a(SB), NOSPLIT|ABIInternal, $16-24
	(a.go:3)	SUBQ	$16, SP
	(a.go:3)	MOVQ	BP, 8(SP)
	(a.go:3)	LEAQ	8(SP), BP
	(a.go:3)	FUNCDATA	$0, gclocals·1a65...
	(a.go:3)	FUNCDATA	$1, gclocals·69c1...
	(a.go:3)	FUNCDATA	$3, gclocals·33cd...
	(a.go:4)	PCDATA	$2, $0
	(a.go:4)	PCDATA	$0, $0
	(a.go:4)	MOVQ	$6, "".n(SP)
	(a.go:5)	PCDATA	$0, $1
	(a.go:5)	CMPQ	"".a+32(SP), $6
	(a.go:5)	JHI	32
	(a.go:5)	JMP	42
	(a.go:6)	PCDATA	$2, $-2
	(a.go:6)	PCDATA	$0, $-2
	(a.go:6)	MOVQ	8(SP), BP
	(a.go:6)	ADDQ	$16, SP
	(a.go:6)	RET
	(a.go:5)	PCDATA	$2, $0
	(a.go:5)	PCDATA	$0, $1
	(a.go:5)	CALL	runtime.panicindex(SB)
	(a.go:5)	UNDEF
	0x0000 48 83 ... 
	0x0010 04 24 ...
	0x0020 48 8b ...
	0x0030 0b      
	rel 43+4 t=8 runtime.panicindex+0
```

```
"".b STEXT nosplit size=34 args=0x28 locals=0x10
	(b.go:3)	TEXT	"".b(SB), NOSPLIT|ABIInternal, $16-40
	(b.go:3)	SUBQ	$16, SP
	(b.go:3)	MOVQ	BP, 8(SP)
	(b.go:3)	LEAQ	8(SP), BP
	(b.go:3)	FUNCDATA	$0, gclocals·33cd...
	(b.go:3)	FUNCDATA	$1, gclocals·33cd...
	(b.go:3)	FUNCDATA	$3, gclocals·33cd...
	(b.go:4)	PCDATA	$2, $0
	(b.go:4)	PCDATA	$0, $0
	(b.go:4)	MOVQ	$4, "".n(SP)
	(b.go:5)	JMP	24
	(b.go:6)	PCDATA	$2, $-2
	(b.go:6)	PCDATA	$0, $-2
	(b.go:6)	MOVQ	8(SP), BP
	(b.go:6)	ADDQ	$16, SP
	(b.go:6)	RET
	0x0000 48 83 ...
	0x0010 04 24 ...
	0x0020 10 c3
```

Even without the optimizations, there are more instructions that the CPU has to run in the case of a.go {n:=6} more than b.go {n:=len(b)-1}.

There are some interesting differences between the two.  The {n:=6} version has a compare statement (CMPQ) and panic statements (runtime.panicindex) while the other version does not have them.

Let's also compile both with another option and see if we get any clues there.

```
$ go tool compile -d=ssa/check_bce/debug=1 a.go
a.go:5:7: Found IsInBounds

$ go tool compile -d=ssa/check_bce/debug=1 b.go
```

So, the compile tool shows no output with this option for b.go while a.go says "Found IsInBounds" at line number 5 (\_ = a[n]).

### Bounds Check Elimination (bce)
From Wikipedia: bounds-checking elimination is a compiler optimization useful in programming languages or runtimes that enforce bounds checking, the practice of checking every index into an array to verify that the index is within the defined valid range of indexes. Its goal is to detect which of these indexing operations do not need to be validated at runtime, and eliminating those checks.

When arrays and slices are being accessed, grow provides safety by checking that the index is valid.  This implies additional instructions.  A language like C does not have this check; instead it is upto the programmer to add it if required or not do it at their own risk.

Go provides the check but is able to eliminate in certain cases when it is able to prove that the index being accessed is within the allowed range.

In the function ```func a(a []int) { n := 6; _ = a[n] }```, Go is not able to prove at compile time that the index 6 will be in the slice that is passed.  However, in the function ```func b(b [5]int) { n := len(b) - 1; _ = b[n] }```, it is guaranted that the index will be within the length of the array of size 5.  Thus Go is able to optimize by eliminating the bounds check.

Exercise: What if we passed a slice into b.go instead of an array.  Is there a bounds check still?  Why or why not?
See c.go

```
3 func c(b []int) {
4     n := len(b) - 1
5     _ = b[n]
6 }
```

```
$ go tool compile -d=ssa/check_bce/debug=1 c.go
c.go:5:7: Found IsInBounds
```

What is the bce output of the case below?  will the compiler be able to eliminate the bounds check?

```
// d.go
func d(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 9
	}
}
```

```
$ go tool compile -d=ssa/check_bce/debug=1 d.go
```

When it is definite that the index will not receive a value outside of its size (on either end), then bce can happen.

### Providing bce Hints

Example 1

```
// e.go
3 func e(b []byte, n int) {
4     for i := 0; i < n; i++ {
5         b[i] = 9
6     }
7 }
```

```
$ go tool compile -d=ssa/check_bce/debug=1 d.go
d.go:5:8: Found IsInBounds
```

Give that this is running inside a loop, the bce will run as many times.  Is there a way to reduce this?  Probably something outside the loop and prior?

```
// f.go
3 func f(b []byte, n int) {
4     _ = b[n-1]
5     for i := 0; i < n; i++ {
6         b[i] = 9
7     }
8 }
```

```
$ go tool compile -d=ssa/check_bce/debug=1 e.go
e.go:4:7: Found IsInBounds
```

Having done the check once outside, we are able to eliminate the remaining checks in the loop.


How about this one?  There are 4 bounds checks.  Can we reduce them?

Example 2

```
// g.go
func g1(b []byte, v uint32) {
	b[0] = byte(v + 48) // Found IsInBounds
	b[1] = byte(v + 49) // Found IsInBounds
	b[2] = byte(v + 50) // Found IsInBounds
	b[3] = byte(v + 51) // Found IsInBounds
}
```

```
// g.go
func g2(b []byte, v uint32) {
	b[3] = byte(v + 51) // Found IsInBounds
	b[0] = byte(v + 48)
	b[1] = byte(v + 49)
	b[2] = byte(v + 50)
}
```

Example 3

```
// h.go
func h1(b []byte, n int) {
	b[n+0] = byte(1) // Found IsInBounds
	b[n+1] = byte(2) // Found IsInBounds
	b[n+2] = byte(3) // Found IsInBounds
	b[n+3] = byte(4) // Found IsInBounds
	b[n+4] = byte(5) // Found IsInBounds
	b[n+5] = byte(6) // Found IsInBounds
}
```

```
func h2(b []byte, n int) {
	b = b[n : n+6] // Found IsSliceInBounds
	b[0] = byte(1)
	b[1] = byte(2)
	b[2] = byte(3)
	b[3] = byte(4)
	b[4] = byte(5)
	b[5] = byte(6)
}
```

Example 4

```
func i1(a, b, c []byte) {
	for i := range a {
		a[i] = b[i] + c[i] // 5:11 Found IsInBounds and 5:12  Found IsInBounds
	}
}
```

```
func i2(a, b, c []byte) {
	_ = b[len(a)-1] // Found IsInBounds
	_ = c[len(a)-1] // Found IsInBounds
	for i := range a {
		a[i] = b[i] + c[i]
	}
}
```

