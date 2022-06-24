package chatboterrors

import (
	"errors"
	"fmt"
)

type Example struct {
	number int
}

func (e *Example) Error() string {
	return fmt.Sprintf("wrong parameter, only number less than %d is accepted", e.number)
}
func newExample(number int) *Example {
	return &Example{
		number: number,
	}

}

var ErrParamaterIsNotNumber error = errors.New("wrong parameter, only number is accepted as parameter")
var ErrParamaterIsTooLarge error = newExample(500)
var ErrParamaterIsToosmall error = errors.New("wrong parameter, parameter must be greater than 0")
var ErrParamaterIsInvalidUrl error = errors.New("wrong parameter, invalid URL")
var ErrParamaterIsTooLargeSetReplicaCount error = newExample(50)
var ErrParamaterIsToosmallSetReplicaCount error = errors.New("wrong parameter, parameter must be greater than 1")
var ErrWrongInsufficientNumberOfParameter error = errors.New("too many or too less parameter provided to this command")
var ErrParamaterIsToosmallPrimeFactorization error = errors.New("wrong parameter, parameter must be greater than 2")
var ErrParamaterIsPrime error = errors.New("wrong parameter, parameter must be greater than 2")
var ErrParamaterIsTooLargeConvert error = newExample(999999999999)
var ErrParamaterIsToosmallConvert error = errors.New("wrong parameter, parameter must be greater than 1")
