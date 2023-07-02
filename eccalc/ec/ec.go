package fp

import (
	"math/big"
)

type FpEC struct {
	x       big.Int
	y       big.Int
	neutral bool
}
