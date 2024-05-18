package primality

// SieveOfEratosthenes is an implementation of the ancient sieve of Eratosthenes
// to generate a list of primes up to the provided integer.
func SieveOfEratosthenes(n uint64) []uint64 {
	if n < 2 {
		return nil // 1 is not prime
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
	// Initialise a slice to store the primes we generate
	primes := make([]uint64, 0, n+1)
	// Primes are remaining false values indexes
	for i, v := range b {
		if !v {
			primes = append(primes, uint64(i))
		}
	}
	return primes
}
