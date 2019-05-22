package main

import (
	"fmt"
	"html/template"
	"sync"
)

var s = `
<h1>{{.PageTitle}}<h1>
<ul>
    {{range .Todos}}
        {{if .Done}}
            <li class="done">{{.Title}}</li>
        {{else}}
            <li>{{.Title}}</li>
        {{end}}
    {{end}}
</ul>
`

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
