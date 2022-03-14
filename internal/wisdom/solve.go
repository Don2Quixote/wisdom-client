package wisdom

import (
	"math/big"
)

// solve defines how to find prime factors of number -
// with uint64s or with big ints and returns result as []*big.Int.
func solve(number *big.Int) []*big.Int {
	if number.IsUint64() {
		factors := PrimeFactors(number.Uint64())

		bigFactors := make([]*big.Int, 0, len(factors))
		for _, factor := range factors {
			bigFactors = append(bigFactors, (&big.Int{}).SetUint64(factor))
		}

		return bigFactors
	}

	return PrimeFactorsBig(number)
}

// Get all prime factors of a given uint64 number n.
func PrimeFactors(n uint64) (pfs []uint64) {
	// Get the number of 2s that divide n.
	for n%2 == 0 {
		pfs = append(pfs, 2)
		n /= 2
	}

	// n must be odd at this point. so we can skip one element
	// (note i += 2).
	for i := uint64(3); i*i <= n; i += 2 {
		// while i divides n, append i and divide n.
		for n%i == 0 {
			pfs = append(pfs, i)
			n /= i
		}
	}

	// This condition is to handle the case when n is a prime number
	// greater than 2.
	if n > 2 {
		pfs = append(pfs, n)
	}

	return
}

// Get all prime factors of a given big.Int number n.
// This does a lot of allocations compared to PrimeFactors for uint64.
func PrimeFactorsBig(n *big.Int) (pfs []*big.Int) {
	// Get the number of 2s that divide n.
	for {
		nc := (&big.Int{}).Set(n)

		if nc.Mod(nc, bigInt2).Cmp(bigInt0) != 0 {
			break
		}

		pfs = append(pfs, bigInt2)
		n.Div(n, bigInt2)
	}

	// n must be odd at this point. so we can skip one element
	// (note i.Add(i, bigInt2)).
	for i := big.NewInt(3); ; i.Add(i, bigInt2) {
		doubleI := (&big.Int{}).Mul(i, i)

		if doubleI.Cmp(n) == 1 {
			break
		}

		for {
			nc := (&big.Int{}).Set(n)

			if nc.Mod(nc, i).Cmp(bigInt0) != 0 {
				break
			}

			pfs = append(pfs, (&big.Int{}).Set(i))
			n.Div(n, i)
		}
	}

	// This condition is to handle the case when n is a prime number
	// greater than 2.
	if n.Cmp(bigInt2) == 1 {
		pfs = append(pfs, n)
	}

	return pfs
}
