// Copyright 2026 The Gc Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gc

import (
	"encoding/binary"
	"go/types"
	"testing"
)

// The linux/ppc row's Go-side shapes must be what the Go toolchain compiles
// the generated code to. The oracle is go/types' own gc sizes for "ppc",
// which only a Go that knows GOARCH=ppc (the go-ppc port) provides.
func TestNewABILinuxPPC(t *testing.T) {
	abi, err := NewABI("linux", "ppc")
	if err != nil {
		t.Fatal(err)
	}
	if abi.ByteOrder != binary.BigEndian {
		t.Fatalf("linux/ppc byte order = %v, want big endian", abi.ByteOrder)
	}
	// C _Bool on Adelie ppc is size 1 (cabi.out).
	if got, want := abi.Types[Bool], (AbiType{1, 1, 1}); got != want {
		t.Errorf("linux/ppc Bool = %v, want %v", got, want)
	}
	sizes := types.SizesFor("gc", "ppc")
	if sizes == nil {
		t.Skip("this Go toolchain does not know GOARCH=ppc")
	}
	for kind, typ := range map[Kind]types.Type{
		Complex128:    types.Typ[types.Complex128],
		Complex64:     types.Typ[types.Complex64],
		Float32:       types.Typ[types.Float32],
		Float64:       types.Typ[types.Float64],
		Int:           types.Typ[types.Int],
		Int16:         types.Typ[types.Int16],
		Int32:         types.Typ[types.Int32],
		Int64:         types.Typ[types.Int64],
		Int8:          types.Typ[types.Int8],
		Interface:     types.NewInterfaceType(nil, nil),
		Pointer:       types.NewPointer(types.Typ[types.Int]),
		Slice:         types.NewSlice(types.Typ[types.Int]),
		String:        types.Typ[types.String],
		Uint:          types.Typ[types.Uint],
		Uint16:        types.Typ[types.Uint16],
		Uint32:        types.Typ[types.Uint32],
		Uint64:        types.Typ[types.Uint64],
		Uint8:         types.Typ[types.Uint8],
		Uintptr:       types.Typ[types.Uintptr],
		UnsafePointer: types.Typ[types.UnsafePointer],
	} {
		want := AbiType{sizes.Sizeof(typ), int(sizes.Alignof(typ)), int(sizes.Alignof(typ))}
		if got := abi.Types[kind]; got != want {
			t.Errorf("linux/ppc %v = %v, want %v (go/types)", kind, got, want)
		}
	}
}
