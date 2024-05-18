package primality

// MaskOfEratosthenes uses a bitmask produced by a modified sieve of
// Eratosthenes in order to compare the desired integer against the
// corresponding index in the mask, determining if it is prime or not
func MaskOfEratosthenes(n uint64) bool {
	if n < 2 || ((n > 2) && n&1 == 0) {
		return false // 1 and evens > 2 are not prime
	}
	// Initialise a boolean slice with all false values
	b := make([]bool, n+1)
	b[0], b[1] = true, true
	// Set all multiples of numbers up to sqrt(n) to true
	for i := uint64(2); i*i <= n; i++ {
		if b[i] {
			continue
		}
		for j := i * i; j <= n; j += i {
			b[j] = true
		}
	}
	// Initialise a slice to store the prime bitmask we generate
	primes := make([]byte, n+1)
	// Primes are remaining false values indexes
	for i, v := range b {
		if !v {
			primes[i] = 0x1
		}
	}
	return primes[n] != 0x0
}
