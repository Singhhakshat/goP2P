package crypto // 1. Every file belongs to a package. This folder name is 'crypto'.

import (
	"crypto/rand" // 2. Standard library for secure random numbers.
	"fmt"         // 3. For "formatted" I/O (printing and string formatting).
)

func GenerateCode() string { // func <function name> <return type>
	b := make([]byte, 16) // 4. Create a byte slice of length 16.
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
