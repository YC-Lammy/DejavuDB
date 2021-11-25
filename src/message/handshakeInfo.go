package message

import (
	"errors"
	"strconv"
	"strings"
)

type Handshakeinfo struct {
	Role string
	Pass string
	Host string
	Port string

	ID uint64 //0 if new born
}

func (h Handshakeinfo) ToBytes() []byte {
	return []byte(strings.Join([]string{
		h.Role, h.Pass, h.Host, h.Port,
		strconv.FormatUint(h.ID, 10)}, ";"))
}

func (h *Handshakeinfo) FromBytes(b []byte) (err error) {
	s := strings.Split(string(b), ";")
	if len(s) != 5 {
		return errors.New("Hanshake requires 5 fields")
	}
	h.Role = s[0]
	h.Pass = s[1]
	h.Host = s[2]
	h.Port = s[3]

	h.ID, err = strconv.ParseUint(s[4], 10, 64)
	return
}
