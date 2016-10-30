package kr

import (
	"bytes"
	"testing"
)

func TestGenWrapEncDec(t *testing.T) {
	ps, err := GeneratePairingSecret()
	if err != nil {
		t.Fatal(err)
	}
	sessionKey, err := RandNBytes(32)
	if err != nil {
		t.Fatal(err)
	}

	encryptedKey, err := sodiumBoxSeal(sessionKey, ps.WorkstationPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	wrappedKeyMessage := append([]byte{HEADER_WRAPPED_KEY}, encryptedKey...)
	remaining, didUnwrap, err := ps.UnwrapKeyIfPresent(wrappedKeyMessage)
	if err != nil {
		t.Fatal(err)
	}
	if remaining != nil {
		t.Fatal()
	}
	if !didUnwrap {
		t.Fatal()
	}
	if !bytes.Equal(sessionKey, *ps.SymmetricSecretKey) {
		t.Fatal("SymmetricSecretKey wrong")
	}

	msg, err := RandNBytes(129)
	if err != nil {
		t.Fatal(err)
	}
	ctxt, err := ps.EncryptMessage(msg)
	if err != nil {
		t.Fatal(err)
	}

	remainingCtxt, didUnwrap, err := ps.UnwrapKeyIfPresent(ctxt)
	if remainingCtxt == nil {
		t.Fatal("should have remaining ciphertext")
	}
	if didUnwrap {
		t.Fatal("was not wrapped key")
	}
	if err != nil {
		t.Fatal(err)
	}

	ptxt, err := ps.DecryptMessage(*remainingCtxt)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(*ptxt, msg) {
		t.Fatal("decrypt failed")
	}
}
