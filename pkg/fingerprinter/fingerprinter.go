//nolint:gosec,revive // it's ok for fingerprint
package fingerprinter

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
)

func SHA1FromStrings(strs ...string) string {
	bytesCnt := 0
	for _, str := range strs {
		bytesCnt += len(str)
	}

	buf := bytes.NewBuffer(make([]byte, 0, bytesCnt+len(strs)-1))
	for i, str := range strs {
		buf.WriteString(str)

		if i < len(strs)-1 {
			buf.WriteByte('.')
		}
	}

	hash := sha1.Sum(buf.Bytes())

	return hex.EncodeToString(hash[:])
}
