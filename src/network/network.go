package network

import (
	"crypto/aes"
	"crypto/cipher"
	"dejavuDB/src/config"
	"dejavuDB/src/lazy"
	"dejavuDB/src/message"
	"errors"
	"net"
	"unsafe"
)

var HandShakeRequiredError = errors.New("Handshake required before read write")

type Conn struct {
	Conn net.Conn

	handshaking bool
	handshaked  bool

	aeskey cipher.Block
}

func NewConn(conn net.Conn) *Conn {
	return &Conn{
		Conn: conn,
	}
}

func (c *Conn) Write(message []byte) (int, error) {
	b, err := AESencrypt(c.aeskey, message)
	if err != nil {
		return 0, err
	}
	return c.writeRaw(b)
}

func (c *Conn) writeRaw(msg []byte) (int, error) {
	if !c.handshaked && !c.handshaking {
		return 0, HandShakeRequiredError
	}
	l := uint64(len(msg))

	c.Conn.Write((*(*[8]byte)(unsafe.Pointer(&l)))[:])
	return c.Conn.Write(msg)
}

func (c *Conn) ReadMesaage() ([]byte, error) {
	b, err := c.readRaw()
	if err != nil {
		return nil, err
	}
	return AESdecrypt(c.aeskey, b)
}

func (c *Conn) readRaw() ([]byte, error) {
	if !c.handshaked && !c.handshaking {
		return nil, HandShakeRequiredError
	}
	var length uint64
	var lenbuf = make([]byte, 8)
	c.Conn.Read(lenbuf)
	leng := [8]byte{}
	copy(leng[:], lenbuf)
	length = *(*uint64)(unsafe.Pointer(&leng))
	message := make([]byte, length)
	_, err := c.Conn.Read(message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (c *Conn) SendHandshake() error {
	c.handshaking = true
	defer func() { c.handshaking = false }()

	private, public := genKey()
	s, err := ExportRsaPublicKeyAsPemStr(public)
	if err != nil {
		return err
	}
	c.writeRaw([]byte(s))

	msg, err := c.readRaw()
	if err != nil {
		return err
	}

	ae := RSA_OAEP_Decrypt(string(msg), *private)
	aesk, err := aes.NewCipher(ae)

	if err != nil {
		return err
	}

	c.aeskey = aesk

	v := message.Handshakeinfo{
		Role: config.Role,
		Pass: config.Password,
		Host: config.Host,
		Port: config.Port,

		ID: config.ID,
	}
	c.Write(v.ToBytes())
	msg, err = c.ReadMesaage()
	if err != nil {
		return err
	}
	if string(msg) != "ok" {
		return errors.New(string(msg))
	}
	c.handshaked = true
	return nil
}

func (c *Conn) RecvHandshake() (*message.Handshakeinfo, error) {
	c.handshaking = true
	defer func() { c.handshaking = false }()

	k, err := c.readRaw()
	if err != nil {
		return nil, err
	}
	key, err := ParseRsaPublicKeyFromPemStr(string(k))
	if err != nil {
		return nil, err
	}

	a := lazy.RandString(32)

	c.writeRaw([]byte(RSA_OAEP_Encrypt(a, *key)))
	aesk, err := aes.NewCipher([]byte(a))
	if err != nil {
		return nil, err
	}
	c.aeskey = aesk

	msg, err := c.ReadMesaage()
	if err != nil {
		c.Write([]byte(err.Error()))
		return nil, err
	}

	handshake := &message.Handshakeinfo{}
	err = handshake.FromBytes(msg)
	if err != nil {
		c.Write([]byte(err.Error()))
		return nil, err
	}

	if handshake.Pass != config.Password {
		c.Write([]byte("password incorrect"))
		return nil, errors.New("password incorrect")
	}

	c.Write([]byte("ok"))
	c.handshaked = true

	return handshake, nil
}

func (c *Conn) ClientHandshake() error {
	c.handshaking = true
	defer func() { c.handshaking = false }()

	k, err := c.readRaw()
	if err != nil {
		return err
	}
	key, err := ParseRsaPublicKeyFromPemStr(string(k))
	if err != nil {
		return err
	}

	a := lazy.RandString(32)

	c.writeRaw([]byte(RSA_OAEP_Encrypt(a, *key)))
	aesk, err := aes.NewCipher([]byte(a))
	if err != nil {
		return err
	}
	c.aeskey = aesk
	c.handshaked = true
	return nil
}
