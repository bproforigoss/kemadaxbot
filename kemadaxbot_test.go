package main

import (
	"testing"
)

func TestPrimeFactorization(t *testing.T) {

	got, _ := primeFactors(600000000)
	want := []string{"2", "2", "2", "2", "2", "2", "2", "2", "2", "3", "5", "5", "5", "5", "5", "5", "5", "5"}

	if len(got) != len(want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("got %q, wanted %q", got, want)
		}
	}
}

type IsPrimeTest struct {
	number   int
	expected bool
}

var IsPrimeTests = []IsPrimeTest{
	IsPrimeTest{2, true},
	IsPrimeTest{4, false},
	IsPrimeTest{53, true},
	IsPrimeTest{1000000, false},
}

func TestIsPrime(t *testing.T) {

	for _, test := range IsPrimeTests {
		if output := IsPrime(test.number); output != test.expected {
			t.Errorf("Output %v not equal to expected %v", output, test.expected)
		}
	}
}

func TestConvert(t *testing.T) {

	got := convert(57412)
	want := "ötvenhétezer-négyszáztizenkettő"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func BenchmarkPrimeFactors(b *testing.B) {
	primeFactors(100)
}
func BenchmarkGenerateBigprime(b *testing.B) {
	generateBigPrime()
}
