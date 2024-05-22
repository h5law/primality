package primality

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrimalityTest_AKS(t *testing.T) {
	primes := make([]int, 0, 168)
	for i := 1; i <= 1000; i++ {
		if aks := AKS(i); aks {
			primes = append(primes, i)
		}
	}
	require.Len(t, primes, 168)
}

func TestPrimalityTest_MillerRabin(t *testing.T) {
	primes := make([]int, 0, 168)
	for i := 1; i <= 1000; i++ {
		if mr := MillerRabin(big.NewInt(int64(i)), 25, true); mr {
			primes = append(primes, i)
		}
	}
	require.Len(t, primes, 168)
}

func TestPrimalityTests_Equal(t *testing.T) {
	primes1 := make([]int, 0, 168)
	for i := 1; i <= 1000; i++ {
		if mr := MillerRabin(big.NewInt(int64(i)), 25, true); mr {
			primes1 = append(primes1, i)
		}
	}
	assert.Len(t, primes1, 168)

	primes2 := make([]int, 0, 168)
	for i := 1; i <= 1000; i++ {
		if aks := AKS(i); aks {
			primes2 = append(primes2, i)
		}
	}
	assert.Len(t, primes2, 168)

	require.Equal(t, primes1, primes2)

}
