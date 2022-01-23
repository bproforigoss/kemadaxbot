package main

import (
	"testing"
)

func TestPrimeFactorization(t *testing.T) {

	got1 := primeFactorization(3884)
	got2 := primeFactorization(100)
	want1 := "2, 2"
	want2 := "2, 2, 5, 5"

	if got1.factorsWithCommas() != want1 {
		t.Errorf("got %q, wanted %q", got1.factorsWithCommas(), want1)
	}
	if got2.factorsWithCommas() != want2 {
		t.Errorf("got %q, wanted %q", got2.factorsWithCommas(), want2)
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

type convertTest struct {
	number   int
	expected string
}

var convertTests = []convertTest{
	convertTest{1, "egy"},
	convertTest{1999, "egyezerkilencszázkilencvenkilenc"},
	convertTest{516784, "ötszáztizenhatezer-hétszáznyolcvannégy"},
	convertTest{1111111111, "egymilliárd-egyszáztizenegymillió-egyszáztizenegyezer-egyszáztizenegy"},
}

func TestConvert(t *testing.T) {
	for _, test := range convertTests {
		if output := convert(test.number); output != test.expected {
			t.Errorf("Output %v not equal to expected %v", output, test.expected)
		}
	}

}

func BenchmarkPrimeFactors(b *testing.B) {
	primeFactorization(100)
}
func BenchmarkGenerateBigprime(b *testing.B) {
	generateBigPrime()
}

