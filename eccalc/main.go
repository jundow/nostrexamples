package main

import (
	"fmt"
	"math/big"

	"eccalc/fp"
)

func main() {
	z := big.NewInt(0)
	ret := big.NewInt(0)

	p64 := int64(29)
	p := big.NewInt(p64)

	for i := int64(0); i < p64; i++ {
		x := big.NewInt(i)
		for j := int64(0); j < p64; j++ {
			y := big.NewInt(j)
			fp.FpMul(x, y, p, z)
			fmt.Print(z)
			fmt.Print("\t")
		}
		fmt.Println()
	}

	for i := int64(1); i < p64; i++ {
		x := big.NewInt(i)
		fp.FpInv(x, p, z)
		fmt.Println(x, z, fp.FpMul(x, z, p, ret))
	}

	fmt.Println("######")

	pow := big.NewInt(0)
	m := big.NewInt(3)
	for i := int64(0); i < p64; i++ {
		pow.SetInt64(i)
		fp.FpPow(m, pow, p, ret)
		fmt.Println("ANSWER", m, pow, p, ret)
	}

}
