package strategy

import (
	"crypto/md5"
	"math/big"
)

func hashFn(key string) *big.Int {
	h := md5.New()

	h.Write([]byte(key))
	md := h.Sum(nil)

	keymodbackends := new(big.Int)

	mdBigInt := new(big.Int).SetBytes(md)

	return keymodbackends.Mod(mdBigInt, big.NewInt(int64(19)))
}
