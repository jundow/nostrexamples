package fp

import (
	"math/big"
)

type FpEC struct {
	Order   big.Int
	x       big.Int
	y       big.Int
	neutral bool
}
