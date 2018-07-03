package totp

import (
	"log"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

var keySecret string
var keyAlgorithm = otp.AlgorithmSHA256
var keyPeriod uint

func Initialise(secret string, period time.Duration) {
	keySecret = secret
	keyPeriod = uint(period / time.Second)
	log.Printf("TOTP: secret=%v, period=%v\n", secret, keyPeriod)
}

func Validate(digits string, now time.Time) bool {
	valid, err := totp.ValidateCustom(digits, keySecret, now, totp.ValidateOpts{
		Algorithm: keyAlgorithm,
		Period:    keyPeriod,
		Digits:    otp.DigitsSix,
	})
	if err != nil {
		log.Printf("Invalid TOTP: %v\n", err)
	}
	return valid
}
