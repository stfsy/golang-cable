package benchmarks

import (
	"regexp"
	"testing"
)

func BenchmarkMatchStringString(b *testing.B) {
	pattern := "^/person/add$"
	for i := 0; i < b.N; i++ {
		path := "/person/add"
		regexp.MatchString(pattern, path)
	}
}

func BenchmarkCompileFindAllStrings(b *testing.B) {
	pattern, _ := regexp.Compile("^/person/add$")
	for i := 0; i < b.N; i++ {
		path := "/person/add"
		pattern.FindAllString(path, -1)
	}
}

func BenchmarkCompileFindStringTwoTimesWTrailingSlash(b *testing.B) {
	pattern, _ := regexp.Compile("^/person/add(/)?$")
	pattern2, _ := regexp.Compile("^/person/add/$")
	for i := 0; i < b.N; i++ {
		path := "/person/add"
		pattern.FindString(path)
		pattern2.FindString(path)
	}
}

func BenchmarkCompileFindStringOptionalTrailingSlash(b *testing.B) {
	pattern, _ := regexp.Compile("^/person/add(/)?$")
	for i := 0; i < b.N; i++ {
		path := "/person/add"
		pattern.FindString(path)
	}
}

func BenchmarkCompileFindString(b *testing.B) {
	pattern, _ := regexp.Compile("^/person/add$")
	for i := 0; i < b.N; i++ {
		path := "/person/add"
		pattern.FindString(path)
	}
}

func BenchmarkCompileFindongesttring(b *testing.B) {
	pattern, _ := regexp.Compile("^/person/add$")
	pattern.Longest()
	for i := 0; i < b.N; i++ {
		path := "/person/add"
		pattern.FindString(path)
	}
}
