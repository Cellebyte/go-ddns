package certbot

import (
	"bytes"
	"encoding/pem"
	"fmt"
)

func DERChainToSingleBuffer(derChain [][]byte) ([]byte, error) {
	var buf bytes.Buffer
	for _, der := range derChain {
		if err := pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
			// pem.Encode never returns a non-nil error, but keep this for completeness.
			return buf.Bytes(), fmt.Errorf("pem encode: %w", err)
		}
	}
	return buf.Bytes(), nil
}

func SingleBufferToDERChain(buf []byte) ([][]byte, error) {
	var out [][]byte
	for {
		var block *pem.Block
		block, buf = pem.Decode(buf)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" {
			// ignore other blocks or return an error if you prefer
			continue
		}
		out = append(out, block.Bytes)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no certificates found")
	}
	return out, nil
}
