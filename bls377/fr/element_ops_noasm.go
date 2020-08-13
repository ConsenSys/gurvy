// +build !amd64

// Copyright 2020 ConsenSys AG
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

// Code generated by goff (v0.3.1) DO NOT EDIT

// Package fr contains field arithmetic operations
package fr

// /!\ WARNING /!\
// this code has not been audited and is provided as-is. In particular,
// there is no security guarantees such as constant time implementation
// or side-channel attack resistance
// /!\ WARNING /!\

import "math/bits"

func mul(z, x, y *Element) {
	_mulGeneric(z, x, y)
}

func square(z, x *Element) {
	_squareGeneric(z, x)
}

// FromMont converts z in place (i.e. mutates) from Montgomery to regular representation
// sets and returns z = z * 1
func fromMont(z *Element) {
	_fromMontGeneric(z)
}

func add(z, x, y *Element) {
	var carry uint64

	z[0], carry = bits.Add64(x[0], y[0], 0)
	z[1], carry = bits.Add64(x[1], y[1], carry)
	z[2], carry = bits.Add64(x[2], y[2], carry)
	z[3], _ = bits.Add64(x[3], y[3], carry)

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[3] < 1345280370688173398 || (z[3] == 1345280370688173398 && (z[2] < 6968279316240510977 || (z[2] == 6968279316240510977 && (z[1] < 6461107452199829505 || (z[1] == 6461107452199829505 && (z[0] < 725501752471715841))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 725501752471715841, 0)
		z[1], b = bits.Sub64(z[1], 6461107452199829505, b)
		z[2], b = bits.Sub64(z[2], 6968279316240510977, b)
		z[3], _ = bits.Sub64(z[3], 1345280370688173398, b)
	}
}

func double(z, x *Element) {
	var carry uint64

	z[0], carry = bits.Add64(x[0], x[0], 0)
	z[1], carry = bits.Add64(x[1], x[1], carry)
	z[2], carry = bits.Add64(x[2], x[2], carry)
	z[3], _ = bits.Add64(x[3], x[3], carry)

	// if z > q --> z -= q
	// note: this is NOT constant time
	if !(z[3] < 1345280370688173398 || (z[3] == 1345280370688173398 && (z[2] < 6968279316240510977 || (z[2] == 6968279316240510977 && (z[1] < 6461107452199829505 || (z[1] == 6461107452199829505 && (z[0] < 725501752471715841))))))) {
		var b uint64
		z[0], b = bits.Sub64(z[0], 725501752471715841, 0)
		z[1], b = bits.Sub64(z[1], 6461107452199829505, b)
		z[2], b = bits.Sub64(z[2], 6968279316240510977, b)
		z[3], _ = bits.Sub64(z[3], 1345280370688173398, b)
	}
}

func sub(z, x, y *Element) {
	var b uint64
	z[0], b = bits.Sub64(x[0], y[0], 0)
	z[1], b = bits.Sub64(x[1], y[1], b)
	z[2], b = bits.Sub64(x[2], y[2], b)
	z[3], b = bits.Sub64(x[3], y[3], b)
	if b != 0 {
		var c uint64
		z[0], c = bits.Add64(z[0], 725501752471715841, 0)
		z[1], c = bits.Add64(z[1], 6461107452199829505, c)
		z[2], c = bits.Add64(z[2], 6968279316240510977, c)
		z[3], _ = bits.Add64(z[3], 1345280370688173398, c)
	}
}

func neg(z, x *Element) {
	if x.IsZero() {
		z.SetZero()
		return
	}
	var borrow uint64
	z[0], borrow = bits.Sub64(725501752471715841, x[0], 0)
	z[1], borrow = bits.Sub64(6461107452199829505, x[1], borrow)
	z[2], borrow = bits.Sub64(6968279316240510977, x[2], borrow)
	z[3], _ = bits.Sub64(1345280370688173398, x[3], borrow)
}