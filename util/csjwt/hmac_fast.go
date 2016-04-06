// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package csjwt

import (
	"crypto"
	"crypto/hmac"
)

// SigningMethodHMACFast implements the HMAC-SHA family of pre-warmed signing methods.
// Less allocations, bytes and a little bit faster but maybe the underlying
// mutex can become the bottleneck.
type SigningMethodHMACFast struct {
	Name Algorithm
	ht   hmacTank
}

func newHMACFast(a Algorithm, h crypto.Hash, key Key) (*SigningMethodHMACFast, error) {
	sm := &SigningMethodHMACFast{
		Name: a,
	}
	if key.Error != nil {
		return nil, key.Error
	}
	if len(key.hmacPassword) == 0 {
		return nil, ErrInvalidKey
	}
	// Can we use the specified hashing method?
	if !h.Available() {
		return nil, ErrHashUnavailable
	}
	sm.ht = newHMACTank(h, key.hmacPassword)
	return sm, nil
}

func NewHMACFast256(key Key) (*SigningMethodHMACFast, error) {
	return newHMACFast(HS256, crypto.SHA256, key)
}

func NewHMACFast384(key Key) (*SigningMethodHMACFast, error) {
	return newHMACFast(HS384, crypto.SHA384, key)
}

func NewHMACFast512(key Key) (*SigningMethodHMACFast, error) {
	return newHMACFast(HS512, crypto.SHA512, key)
}

func (m *SigningMethodHMACFast) Alg() string {
	return m.Name.String()
}

// Verify the signature of HSXXX tokens.  Returns nil if the signature is valid.
func (m *SigningMethodHMACFast) Verify(signingString, signature []byte, _ Key) error {

	// Decode signature, for comparison
	sig, err := DecodeSegment(signature)
	if err != nil {
		return err
	}

	// This signing method is symmetric, so we validate the signature
	// by reproducing the signature from the signing string and key, then
	// comparing that against the provided signature.
	hasher := m.ht.get()
	defer m.ht.put(hasher)

	if _, err := hasher.Write(signingString); err != nil {
		return err
	}

	if !hmac.Equal(sig, hasher.Sum(nil)) {
		return ErrSignatureInvalid
	}

	// No validation errors.  Signature is good.
	return nil
}

// Sign implements the Sign method from SigningMethod interface.
func (m *SigningMethodHMACFast) Sign(signingString []byte, _ Key) ([]byte, error) {

	hasher := m.ht.get()
	defer m.ht.put(hasher)

	if _, err := hasher.Write(signingString); err != nil {
		return nil, err
	}

	return EncodeSegment(hasher.Sum(nil)), nil
}