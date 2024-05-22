package primality

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAKSSteps_One(t *testing.T) {
	composite := basePowerCheck(31) // prime
	require.False(t, composite)
	composite = basePowerCheck(49) // 7^2
	require.True(t, composite)
	composite = basePowerCheck(121) // 11^2
	require.True(t, composite)
}

func TestAKSSteps_Two(t *testing.T) {
	assert.Equal(t, 8, findCorrectOrder(7))
	assert.Equal(t, 25, findCorrectOrder(31))
}

func TestAKSSteps_FastPowerMod(t *testing.T) {
	require.Equal(t, 1, fastPowerMod(4, 3, 7))
	require.Equal(t, 0, fastPowerMod(4, 4, 8))
	require.Equal(t, 27, fastPowerMod(29, 17, 31))
}

func TestAKSSteps_PolynomialExpansion(t *testing.T) {
	p1 := polynomialExpansion(3, 0)
	require.Equal(t, []int{1, 3, 3, 1}, p1)

	p2 := polynomialExpansion(4, 2)
	require.Equal(t, []int{16, 32, 24, 8, 1}, p2)

	p3 := polynomialExpansion(5, -3)
	require.Equal(t, []int{-243, 405, -270, 90, -15, 1}, p3)
}

func TestAKSSteps_PolynomialRemainder(t *testing.T) {
	p1 := polynomialExpansion(6, 3)
	p2 := polynomialExpansion(5, 2)
	rem := polynomialRemainder(p1, p2)
	require.Equal(t, []int{473, 786, 495, 140, 15}, rem)

	p1 = polynomialExpansion(5, 3)
	p2 = polynomialExpansion(3, -1)
	rem = polynomialRemainder(p1, p2)
	require.Equal(t, []int{384, 0, 640}, rem)
}

func TestAKSSteps_PolynomialSubtraction(t *testing.T) {
	p1 := polynomialExpansion(6, 3)
	p2 := polynomialExpansion(5, 2)
	sub := polynomialSubtraction(p1, p2)
	require.Equal(t, []int{697, 1378, 1135, 500, 125, 17, 1}, sub)

	sub = polynomialSubtraction(p2, p1)
	require.Equal(t, []int{-697, -1378, -1135, -500, -125, -17, -1}, sub)
}

func TestAKSSteps_PolymomialScalarMultiply(t *testing.T) {
	p1 := polynomialExpansion(4, 2)
	scalar := 3
	mul := polynomialMultiplyScalar(p1, scalar)
	require.Equal(t, []int{48, 96, 72, 24, 3}, mul)

	p2 := polynomialExpansion(3, -2)
	scalar = 6
	mul = polynomialMultiplyScalar(p2, scalar)
	require.Equal(t, []int{-48, 72, -36, 6}, mul)
}
