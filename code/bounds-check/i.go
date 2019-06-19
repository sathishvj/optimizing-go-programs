package main

func i1(a, b, c []byte) {
	for i := range a {
		a[i] = b[i] + c[i] // 5:11 Found IsInBounds and 5:12  Found IsInBounds
	}
}

func i2(a, b, c []byte) {
	_ = b[len(a)-1] // Found IsInBounds
	_ = c[len(a)-1] // Found IsInBounds
	for i := range a {
		a[i] = b[i] + c[i]
	}
}
