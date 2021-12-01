package fileservice

import (
	"crypto/sha256"
	"os"
	"runtime"
	"unsafe"
)

type DB struct {
	file       *os.File
	keylengths map[[32]byte]uint64
}

func OpenFile(filename string) (*DB, error) {
	file, err := os.Open(filename)
	if err == os.ErrNotExist {
		file, err = os.Create(filename)
		var l uint64 = 0
		file.Write((*(*[8]byte)(unsafe.Pointer(&l)))[:])
		file.Seek(8, 0)
	}
	if err != nil {
		return nil, err
	}
	keylengths := map[[32]byte]uint64{}

	db := new(DB)
	db.file = file
	db.keylengths = keylengths
	runtime.SetFinalizer(db, func(db *DB) {
		db.Close()
	})

	return db, nil
}

func (db *DB) Set(key, val string) {
	h := sha256.Sum256([]byte(key))
	if _, ok := db.keylengths[h]; ok {

		var index uint64
		for i := 0; i < len(db.keylengths); i++ {
			b := make([]byte, 40)
			db.file.ReadAt(b, int64(i*40)+8)
			c := b[32:]
			index += *(*uint64)(unsafe.Pointer(&c))
			if string(b[:33]) == string(h[:]) {
				db.file.Seek(int64(index+uint64(len(db.keylengths))*40)+8, 0)

				i := 0
				for {
					chunk := make([]byte, 1000000)
					l, err := db.file.Read(chunk)
					if err != nil {
						break
					}
					i++
				}
				return
			}

		}
		db.writeHeader(h, uint64(len(val)))
	} else {
		db.writeHeader(h, uint64(len(val)))
		db.appendData(val)
	}
}

func (db *DB) Get(key string) string {
	return ""
}

func (db *DB) Delete(key string) {
	h := sha256.Sum256([]byte(key))
	if _, ok := db.keylengths[h]; ok {
		delete(db.keylengths, h)
	}
}

func (db *DB) writeHeader(key [32]byte, length uint64) {

	if _, ok := db.keylengths[key]; ok {
		for i := 0; i < len(db.keylengths); i++ {
			b := make([]byte, 40)
			db.file.ReadAt(b, int64(i*40)+8)
			if string(b[:33]) == string(key[:]) {
				db.file.WriteAt((*(*[8]byte)(unsafe.Pointer(&length)))[:], int64(i*40)+8+32)
			}
		}
	} else {
		l := uint64(len(db.keylengths) + 1)
		db.file.WriteAt((*(*[8]byte)(unsafe.Pointer(&l)))[:], 0)
		db.file.WriteAt((*(*[8]byte)(unsafe.Pointer(&length)))[:], int64(l*40)+8+32)
	}
	db.keylengths[key] = length
}

func (db *DB) appendData(val string) {
	db.file.Seek(0, 2)
	db.file.Write([]byte(val))
}

func (db *DB) Close() {}
