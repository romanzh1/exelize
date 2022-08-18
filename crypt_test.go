// Copyright 2016 - 2022 The excelize Authors. All rights reserved. Use of
// this source code is governed by a BSD-style license that can be found in
// the LICENSE file.
//
// Package excelize providing a set of functions that allow you to write to and
// read from XLAM / XLSM / XLSX / XLTM / XLTX files. Supports reading and
// writing spreadsheet documents generated by Microsoft Excel™ 2007 and later.
// Supports complex components by high compatibility, and provided streaming
// API for generating or reading data from a worksheet with huge amounts of
// data. This library needs Go version 1.15 or later.

package excelize

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	f, err := OpenFile(filepath.Join("test", "encryptSHA1.xlsx"), Options{Password: "password"})
	assert.NoError(t, err)
	assert.EqualError(t, f.SaveAs(filepath.Join("test", "BadEncrypt.xlsx"), Options{Password: "password"}), ErrEncrypt.Error())
	assert.NoError(t, f.Close())
}

func TestEncryptionMechanism(t *testing.T) {
	mechanism, err := encryptionMechanism([]byte{3, 0, 3, 0})
	assert.Equal(t, mechanism, "extensible")
	assert.EqualError(t, err, ErrUnsupportedEncryptMechanism.Error())
	_, err = encryptionMechanism([]byte{})
	assert.EqualError(t, err, ErrUnknownEncryptMechanism.Error())
}

func TestHashing(t *testing.T) {
	assert.Equal(t, hashing("unsupportedHashAlgorithm", []byte{}), []byte(nil))
}

func TestGenISOPasswdHash(t *testing.T) {
	for hashAlgorithm, expected := range map[string][]string{
		"MD4":     {"2lZQZUubVHLm/t6KsuHX4w==", "TTHjJdU70B/6Zq83XGhHVA=="},
		"MD5":     {"HWbqyd4dKKCjk1fEhk2kuQ==", "8ADyorkumWCayIukRhlVKQ=="},
		"SHA-1":   {"XErQIV3Ol+nhXkyCxrLTEQm+mSc=", "I3nDtyf59ASaNX1l6KpFnA=="},
		"SHA-256": {"7oqMFyfED+mPrzRIBQ+KpKT4SClMHEPOZldliP15xAA=", "ru1R/w3P3Jna2Qo+EE8QiA=="},
		"SHA-384": {"nMODLlxsC8vr0btcq0kp/jksg5FaI3az5Sjo1yZk+/x4bFzsuIvpDKUhJGAk/fzo", "Zjq9/jHlgOY6MzFDSlVNZg=="},
		"SHA-512": {"YZ6jrGOFQgVKK3rDK/0SHGGgxEmFJglQIIRamZc2PkxVtUBp54fQn96+jVXEOqo6dtCSanqksXGcm/h3KaiR4Q==", "p5s/bybHBPtusI7EydTIrg=="},
	} {
		hashValue, saltValue, err := genISOPasswdHash("password", hashAlgorithm, expected[1], int(sheetProtectionSpinCount))
		assert.NoError(t, err)
		assert.Equal(t, expected[0], hashValue)
		assert.Equal(t, expected[1], saltValue)
	}
}