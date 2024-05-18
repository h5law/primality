package primality

// MaskOfEratosthenes uses a bitmask produced by a modified sieve of
// Eratosthenes in order to compare the desired integer against the
// corresponding index in the mask
func MaskOfEratosthenes(n int) bool {
	if n < 2 || ((n > 2) && n&1 == 0) {
		return false
	}
	bz := SieveOfEratosthenes(uint64(n))
	return bz[n] != 0
}
