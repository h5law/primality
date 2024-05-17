package primality

import "math"

// Check n != a^b for a,b > 1, returning true if it is otherwise false if not.
func basePowerCheck(n uint64) bool {
	for i := 2.0; i < math.Log2(float64(n)); i++ {
		a := math.Pow(float64(n), 1.0/i)
		// Check if the rounded value a is equal to its integer value.
		if int(math.Floor(a+0.5)) == int(a) {
			return true // n is a composite number.
		}
	}
	return false
}

// Find the order such that Ord_r(n) > (log_2(n))^2
func findCorrectOrder(n uint64) uint64 {
	maxK := math.Floor(math.Pow(math.Log2(float64(n)), 2))
	maxR := max(3.0, math.Ceil(math.Pow(math.Log2(float64(n)), 5)))
	nextR := true
	i := 2.0
	for i = 2.0; nextR && i < maxR; i++ {
		nextR = false
		for j := 1.0; !nextR && j <= maxK; j++ {
			nextR = math.Mod(math.Pow(float64(n), j), i) == 1 ||
				math.Mod(math.Pow(float64(n), j), i) == 0
		}
	}
	return uint64(i - 1.0)
}

func gcd(a, b uint64) uint64 {
	for a != b {
		if a > b {
			a -= b
		} else {
			b -= a
		}
	}
	return a
}

func gcdChecker(n, r uint64) bool {
	// Check all the gcd values in the interval [2, Ord_r(n)]
	for i := r; r > 1; r-- {
		d := gcd(i, n)
		// Check if i and n are coprime
		if d > 1 && d < n {
			return true // n is a composite number
		}
	}
	return false
}

func eulerTotient(n uint64) uint64 {
	res := n // Initialize result as n
	p := uint64(2)
	for p*p <= n {
		if n%p == 0 {
			for n%p == 0 {
				n /= p
			}
			res -= res / p
		}
		p++
	}
	if n > 1 {
		res -= res / n
	}
	return res
}

func polyModChecker()
