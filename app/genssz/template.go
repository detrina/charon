// Copyright © 2022-2023 Obol Labs Inc. Licensed under the terms of a Business Source License 1.1

package main

var tmpl = `package {{.Package}}

// Code generated by genssz. DO NOT EDIT.

import (
	ssz "github.com/ferranbt/fastssz"

	"github.com/obolnetwork/charon/app/errors"
	"github.com/obolnetwork/charon/app/z"
)

{{range .Types}}
{{ $abbr := .Abbr }}

// HashTreeRootWith ssz hashes the {{.Name}} object with a hasher
func ({{$abbr}} {{.Name}}) HashTreeRootWith(hw ssz.HashWalker) (err error) {
	indx := hw.Index()

{{range .Fields}}
	// Field {{.Index}}: '{{.Name}}' ssz:"{{.SSZTag}}"
	{{- if .IsUint64}}
	hw.PutUint64(uint64({{$abbr}}.{{.Name}}))
	{{else if .IsByteList}}
	err = putByteList(hw, []byte({{$abbr}}.{{.Name}}[:]), {{.Size}}, "{{.Name}}")
	if err != nil {
		return err
	}
	{{else if .IsBytesN}}
	err = putBytesN(hw, []byte({{$abbr}}.{{.Name}}[:]), {{.Size}})
	if err != nil {
		return err
	}
	{{else if .IsComposite}}
	err = {{$abbr}}.{{.Name}}.HashTreeRootWith(hw)
	if err != nil {
		return err
	}
	{{else if .IsCompositeList}}
	{
		listIdx := hw.Index()
		for _, item := range {{$abbr}}.{{.Name}} {
			err = item.HashTreeRootWith(hw)
			if err != nil {
				return err
			}
		}

		hw.MerkleizeWithMixin(listIdx, uint64(len({{$abbr}}.{{.Name}})), uint64({{.Size}}))
	}
	{{end}}
{{end}}

	hw.Merkleize(indx)

	return nil
}
{{end}}

// putByteList appends a ssz byte list.
// See reference: github.com/attestantio/go-eth2-client/spec/bellatrix/executionpayload_encoding.go:277-284.
func putByteList(h ssz.HashWalker, b []byte, limit int, field string) error {
	elemIndx := h.Index()
	byteLen := len(b)
	if byteLen > limit {
		return errors.Wrap(ssz.ErrIncorrectListSize, "put byte list", z.Str("field", field))
	}
	h.AppendBytes32(b)
	h.MerkleizeWithMixin(elemIndx, uint64(byteLen), uint64(limit+31)/32)

	return nil
}

// putByteList appends b as a ssz fixed size byte array of length n.
func putBytesN(h ssz.HashWalker, b []byte, n int) error {
	if len(b) > n {
		return errors.New("bytes too long", z.Int("n", n), z.Int("l", len(b)))
	}

	h.PutBytes(leftPad(b, n))

	return nil
}

// leftPad returns the byte slice left padded with zero to ensure a length of at least l.
func leftPad(b []byte, l int) []byte {
	for len(b) < l {
		b = append([]byte{0x00}, b...)
	}

	return b
}

`