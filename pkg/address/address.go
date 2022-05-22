package address

import (
	"crypto/sha256"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

/*
  Create a blockchain address by following the steps in the link below
  Quote: https://www.oreilly.com/library/view/mastering-bitcoin/9781491902639/ch04.html
*/

func New(pubx, prix []byte) string {
	// 1. hash by sha256
	h1 := sha256.New()
	h1.Write(pubx)
	h1.Write(prix)
	digest1 := h1.Sum(nil)

	// 2. ripemd160 hash
	h2 := ripemd160.New()
	h2.Write(digest1)
	digest2 := h2.Sum(nil)

	// 3. add version byte in front of ripemd160 hash
	vd3 := make([]byte, 21)
	vd3[0] = 0x00
	copy(vd3[1:], digest2[:])

	// 4. perform sha-256 hash on the extended ripemd160 result.
	h4 := sha256.New()
	h4.Write(vd3)
	digest4 := h4.Sum(nil)

	// 5. perform sha-256 hash on the result of the previous SHA-256 hash.
	h5 := sha256.New()
	h5.Write(digest4)
	digest5 := h5.Sum(nil)

	// 6. take the first 4 bytes of the second SHA-256 hash for checksum.
	chsum := digest5[:4]

	// 7. Add the 4 checksum bytes from 7 at the end of extended ripemd160 hash from 4
	dc7 := make([]byte, 25)
	copy(dc7[:21], vd3[:])
	copy(dc7[21:], chsum[:])

	// 8. Convert the result from a byte string into base58
	address := base58.Encode(dc7)
	return address
}
