package xrand

import "math/rand"

func SliceElement[t any](s []t) t {
	if len(s) == 0 {
		panic("SliceElement called on empty slice")
	}
	return s[rand.Intn(len(s))]
}
