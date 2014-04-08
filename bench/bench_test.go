package main

import (
	"fmt"
	"testing"
)

func mylcs(x, y string) (lcs string) {
	lenx := len(x)
	leny := len(y)
	if lenx == 0 || leny == 0 {
		return ""
	}
	i := 0
	for {
		i++
		if lenx < i || leny < i {
			break
		}

		if x[lenx-i] == y[leny-i] {
			lcs = string(x[lenx-i]) + lcs
		} else {
			break
		}
	}

	var lcsx, lcsy string
	if lenx >= i-1 && leny >= i {
		lcsx = mylcs(x[:lenx-i+1], y[:leny-i])
	}
	if lenx >= i && leny >= i-1 {
		lcsy = mylcs(x[:lenx-i], y[:leny-i+1])
	}

	if len(lcsx) > len(lcsy) {
		return lcsx + lcs
	}

	return lcsy + lcs
}

func recursivelcs(a, b string) string {
	aLen := len(a)
	bLen := len(b)
	if aLen == 0 || bLen == 0 {
		return ""
	} else if a[aLen-1] == b[bLen-1] {
		return recursivelcs(a[:aLen-1], b[:bLen-1]) + string(a[aLen-1])
	}
	x := recursivelcs(a, b[:bLen-1])
	y := recursivelcs(a[:aLen-1], b)
	if len(x) > len(y) {
		return x
	}
	return y
}

func dynamiclcs(a, b string) string {
	aLen := len(a)
	bLen := len(b)
	lengths := make([][]int, aLen+1)
	for i := 0; i <= aLen; i++ {
		lengths[i] = make([]int, bLen+1)
	}
	// row 0 and column 0 are initialized to 0 already

	for i := 0; i < aLen; i++ {
		for j := 0; j < bLen; j++ {
			if a[i] == b[j] {
				lengths[i+1][j+1] = lengths[i][j] + 1
			} else if lengths[i+1][j] > lengths[i][j+1] {
				lengths[i+1][j+1] = lengths[i+1][j]
			} else {
				lengths[i+1][j+1] = lengths[i][j+1]
			}
		}
	}

	for _, r := range lengths {
		fmt.Println(r)
	}

	// read the substring out from the matrix
	s := make([]byte, 0, lengths[aLen][bLen])
	for x, y := aLen, bLen; x != 0 && y != 0; {
		if lengths[x][y] == lengths[x-1][y] {
			x--
		} else if lengths[x][y] == lengths[x][y-1] {
			y--
		} else {
			s = append(s, a[x-1])
			x--
			y--
		}
	}
	// reverse string
	r := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		r[i] = s[len(s)-1-i]
	}
	return string(r)
}

// func lcs_greedy(x,y) {
// 	 symbols := map[byte]int{}
// 		r:=0
// 		p:=0
// 		p1,L:=0,idx
// 		m:=x.length
// 		n:=y.length
// 		S := new Buffer(m<n?n:m)

// 	func popsym(index int){
// 		 s := x[index]
// 		pos := symbols[s]+1
// 		pos = y.indexOf(s,pos>r?pos:r);
// 		if(pos===-1){pos=n;}
// 		symbols[s]=pos;
// 		return pos;
// 	}

// 	p1 = popsym(0);
// 	for(i=0;i < m;i++){
// 		p = (r===p)?p1:popsym(i);
// 		p1 = popsym(i+1);
// 		idx=(p > p1)?(i++,p1):p;
// 		if(idx===n){p=popsym(i);}
// 		else{
// 			r=idx;
// 			S[L++]=x.charCodeAt(i);
// 		}
// 	}
// 	return S.toString('utf8',0,L);
// }

func BenchmarkMyLCS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mylcs("abcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefg", "abckefga")
	}
}

func BenchmarkRecursiveLCS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		recursivelcs("abcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefg", "abckefga")
	}
}

func BenchmarkDynamicLCS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dynamiclcs("abcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefg", "abckefga")
	}
}

func TestFoo(t *testing.T) {
	mylcs("ABCDEFG", "ABCDOFK")
}

func TestMyLCS(t *testing.T) {
	println(mylcs("abcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefg", "abckefgaabckefgaabckefgaabckefgaabckefgaabckefgaabckefgaabckefgaabckefgaabckefgaabckefga"))
}

func TestRecursiveLCS(t *testing.T) {
	println(recursivelcs("abcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefgabcdefg", "abckefgaabckefgaabckefgaabckefgaabckefgaabckefgaabckefgaabckefgaabckefgaabckefgaabckefga"))
}

func TestDynamicLCS(t *testing.T) {
	println(dynamiclcs("abcdefg", "abckefg`"))
}
