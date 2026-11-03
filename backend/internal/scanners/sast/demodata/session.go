package democorp

import (
	"fmt"
	"math/rand"
)

func NewSessionToken() string {
	return fmt.Sprintf("sess_%d", rand.Intn(1000000000))
}
