package lazy

import (
	"crypto/aes"
	"math"
	"net"
	"strings"
	"time"

	"math/rand"
)

// a common go file to get information
func getMacAddrs() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}

func Get_first_mac_addr() string {
	a, _ := getMacAddrs()
	return a[0]
}

func Difference_str_arr(a, b []string) []string {

	ok := true
	var list []string

	for _, v := range b {
		for _, x := range a {
			if v == x {
				ok = false
				break
			}
		}
		if ok {
			list = append(list, v)
		}
		ok = true
	}
	return list
}

func RemoveDuplicateStrings(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func RemoveItem(slice []string, item string) []string {
	new := []string{}
	for _, v := range slice {
		if item != v {
			new = append(new, v)
		}
	}
	return new

}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func stringSliceIndex(s []string, str string) int {
	for i, v := range s {
		if v == str {
			return i
		}
	}
	return -1
}

func waitUntil(condition func() bool, execute func(), duration time.Duration) {
	for {
		if condition() {
			execute()
			break
		}
		time.Sleep(duration)
	}
}

func RandString(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits := 6                    // 6 bits to represent a letter index
	letterIdxMask := 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax := 63 / letterIdxBits    // # of letter indices fitting in 63 bits
	src := rand.NewSource(time.Now().UnixNano())
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(int(cache) & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

func EncryptAES(key []byte, plaintext []byte) ([]byte, error) {
	// create cipher
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// allocate space for ciphered data
	out := make([]byte, len(plaintext))

	// encrypt
	c.Encrypt(out, plaintext)
	// return hex string
	return out, nil
}

func DecryptAES(key []byte, ciphertext []byte) ([]byte, error) {

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	return pt, nil
}

func Hash64(s []byte) uint64 {
	var h uint64 = 0
	for pos, char := range s {
		h += uint64(char) * uint64(math.Pow(31, float64(len(s)-pos+1)))
	}
	return h
}
