package idutil

import (
	"crypto/rand"
	"github.com/sony/sonyflake"
	"github.com/speps/go-hashids/v2"
	"istomyang.github.com/like-iam/component-base/util/iputil"
)

var sf = sonyflake.NewSonyflake(sonyflake.Settings{
	MachineID: func() (uint16, error) {
		ip := iputil.GetLocalIP()
		return uint16(ip[2])<<8 + uint16(ip[3]), nil
	},
})

// GetUniqId return distributed unique id.
func GetUniqId() (uint64, error) {
	return sf.NextID()
}

const (
	AlphabetL = "abcdefghijklmnopqrstuvwxyz"
	AlphabetU = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Number    = "0123456789"
)

// GetInstanceId use prefix, database's id and length to generate a hash code like: user-2v69o5.
// generally, we use 6 length, with AlphabetL + Number, we has (26+10)^6 = 21,7678,2336 ids can use.
func GetInstanceId(id uint64, prefix string, length int) (string, error) {
	hd := hashids.NewData()
	hd.Alphabet = AlphabetL + Number
	hd.MinLength = 6
	hd.Salt = "ap089b"
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}

	e, err := h.Encode([]int{int(id >> 32), int(id)})
	if err != nil {
		return "", err
	}

	return prefix + e, nil
}

// GetRandString generates a rand string with char set and length.
func GetRandString(charset string, n int) (string, error) {
	var randomness = make([]byte, n)
	read, err := rand.Read(randomness)
	if err != nil || read != n {
		return "", err
	}

	var r = make([]rune, n)
	var csr = []rune(charset)

	for i, rn := range randomness {
		r[i] = csr[int(rn)%len(csr)]
	}

	return string(r), nil
}
