* Returning a local variable in c would cause errors. But it is possible in Go.
* "Note that, unlike in C, it’s perfectly OK to return the address of a local variable; the storage associated with the variable survives after the function returns."
* 

Run:
go run -gcflags '-m -l' 1.go

References:
[Escape Analysis in Go](https://scvalex.net/posts/29/)
