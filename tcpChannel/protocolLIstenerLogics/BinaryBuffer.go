package protocolLIstenerLogics

//BIG TODO MAKE THIS GENERIC LIKE IN THE GENERICS YOUTUBE VIDEO
import (
	"golang.org/x/exp/constraints"
	_ "golang.org/x/exp/constraints"
)

type CustomData interface {
	constraints.Complex
}
type TestGenNumb[T constraints.Unsigned] struct {
	num T
}

/************************* NETWORK MESSAGES **************************/

//ALWAYS SEND A DEFAULT MESSAGE HANDLERID, IN CASE THE DEVELOPER DOES NOT WANT OT USE THE NETWORKMESSAGE STRUCT
