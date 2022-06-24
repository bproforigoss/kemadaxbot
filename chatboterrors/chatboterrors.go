package chatboterrors

import "errors"

var ErrParamaterIsNotNumber error = errors.New("wrong parameter, only number is accepted as parameter")
var ErrParamaterIsTooLarge error = errors.New("wrong parameter, only number less than 500 is accepted")
var ErrParamaterIsToosmall error = errors.New("wrong parameter, parameter must be greater than 0")
var ErrParamaterIsInvalidUrl error = errors.New("wrong parameter, invalid URL")
var ErrParamaterIsTooLargeSetReplicaCount error = errors.New("wrong parameter, only number less than 50 is accepted")
var ErrParamaterIsToosmallSetReplicaCount error = errors.New("wrong parameter, parameter must be greater than 1")
var ErrWrongInsufficientNumberOfParameter error = errors.New("too many or too less parameter provided to this command")
var ErrParamaterIsToosmallPrimeFactorization error = errors.New("wrong parameter, parameter must be greater than 2")
var ErrParamaterIsPrime error = errors.New("wrong parameter, parameter must be greater than 2")
var ErrParamaterIsTooLargeConvert error = errors.New("wrong parameter, only number less than 999.999.999.999 is accepted")
var ErrParamaterIsToosmallConvert error = errors.New("wrong parameter, parameter must be greater than 1")
