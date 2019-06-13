package main

import (
	"math/rand"
	"strconv"
	"testing"
)

var NumItems int = 1000000

func BenchmarkMapStringKeys(b *testing.B) {
	m := make(map[string]string)
	k := make([]string, 0)

	for i := 0; i < NumItems; i++ {
		key := strconv.Itoa(rand.Intn(NumItems))
		//key += ` is the key value that is being used. `
		key += ` is the key value that is being used and a shakespeare sonnet. ` + sonnet106
		m[key] = "value" + strconv.Itoa(i)
		k = append(k, key)
	}

	i := 0
	l := len(m)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, ok := m[k[i]]; ok {
		}

		i++
		if i >= l {
			i = 0
		}
	}
}

func BenchmarkMapIntKeys(b *testing.B) {
	m := make(map[int]string)
	k := make([]int, 0)

	for i := 0; i < NumItems; i++ {
		key := rand.Intn(NumItems)
		m[key] = "value" + strconv.Itoa(i)
		k = append(k, key)
	}

	i := 0
	l := len(m)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, ok := m[k[i]]; ok {
		}

		i++
		if i >= l {
			i = 0
		}
	}
}

var sonnet106 = `When in the chronicle of wasted time
I see descriptions of the fairest wights,
And beauty making beautiful old rhyme
In praise of ladies dead, and lovely knights,
Then, in the blazon of sweet beauty’s best,
Of hand, of foot, of lip, of eye, of brow,
I see their antique pen would have express’d
Even such a beauty as you master now.
So all their praises are but prophecies
Of this our time, all you prefiguring;
And, for they look’d but with divining eyes,
They had not skill enough your worth to sing:
For we, which now behold these present days,
Had eyes to wonder, but lack tongues to praise.`
