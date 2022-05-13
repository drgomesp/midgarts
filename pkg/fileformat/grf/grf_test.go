package grf_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/drgomesp/midgarts/pkg/fileformat/grf"
	"github.com/stretchr/testify/assert"
)

const (
	dataPath = "./../../data"
)

func TestEntryHeaders(t *testing.T) {
	var tests = []struct {
		Name             string
		ExpectedFileName string
		ExpectedEntries  map[string]*grf.Entry
	}{
		{
			Name:             "load file with raw data",
			ExpectedFileName: "raw",
			ExpectedEntries: map[string]*grf.Entry{
				"raw": {
					Header: grf.EntryHeader{
						CompressedSize:        74,
						CompressedSizeAligned: 74,
						UncompressedSize:      74,
						Flags:                 0x01,
						Offset:                0,
					},
					Data: bytes.NewBuffer(nil),
				},
				"corrupted": {
					Header: grf.EntryHeader{
						CompressedSize:        132,
						CompressedSizeAligned: 123,
						UncompressedSize:      20,
						Flags:                 0x03,
						Offset:                34,
					},
					Data: bytes.NewBuffer(nil),
				},
				"compressed": {
					Header: grf.EntryHeader{
						CompressedSize:        16,
						CompressedSizeAligned: 16,
						UncompressedSize:      74,
						Flags:                 0x01,
						Offset:                74,
					},
					Data: bytes.NewBuffer(nil),
				},
				"compressed-des-header": {
					Header: grf.EntryHeader{
						CompressedSize:        16,
						CompressedSizeAligned: 16,
						UncompressedSize:      74,
						Flags:                 0x05,
						Offset:                90,
					},
					Data: bytes.NewBuffer(nil),
				},
				"compressed-des-full": {
					Header: grf.EntryHeader{
						CompressedSize:        16,
						CompressedSizeAligned: 16,
						UncompressedSize:      74,
						Flags:                 0x03,
						Offset:                106,
					},
					Data: bytes.NewBuffer(nil),
				},
				"big-compressed-des-full": {
					Header: grf.EntryHeader{
						CompressedSize:        361,
						CompressedSizeAligned: 368,
						UncompressedSize:      658,
						Flags:                 0x03,
						Offset:                122,
					},
					Data: bytes.NewBuffer(nil),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			grfFile, err := grf.Load(fmt.Sprintf("%s/%s", dataPath, "with-files.grf"))
			assert.NoError(t, err)
			assert.Equal(t, tt.ExpectedEntries, grfFile.GetEntries("/"))
		})
	}
}

func TestEntryContents(t *testing.T) {
	var tests = []struct {
		FilePath        string
		Name            string
		EntryName       string
		ExpectedDataStr string
	}{
		{
			FilePath:        fmt.Sprintf("%s/%s", dataPath, "with-files.grf"),
			Name:            "load file without compression or encryption",
			EntryName:       "raw",
			ExpectedDataStr: "client client client client client client client client client client client client client client client",
		},
		{
			FilePath:        fmt.Sprintf("%s/%s", dataPath, "with-files.grf"),
			Name:            "load file with compression and no encryption",
			EntryName:       "compressed",
			ExpectedDataStr: "client client client client client client client client client client client client client client client",
		},
		{
			FilePath:        fmt.Sprintf("%s/%s", dataPath, "with-files.grf"),
			Name:            "load file with compression and partial encryption",
			EntryName:       "compressed-des-header",
			ExpectedDataStr: "client client client client client client client client client client client client client client client",
		},
		{
			FilePath:        fmt.Sprintf("%s/%s", dataPath, "with-files.grf"),
			Name:            "load file with compression and full encryption",
			EntryName:       "compressed-des-full",
			ExpectedDataStr: "client client client client client client client client client client client client client client client",
		},
		{
			FilePath:        fmt.Sprintf("%s/%s", dataPath, "with-files.grf"),
			Name:            "load big file with compression and full encryption",
			EntryName:       "big-compressed-des-full",
			ExpectedDataStr: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed venenatis bibendum venenatis. Aliquam quis velit urna. Suspendisse nec posuere sem. Donec risus quam, vulputate sed augue ultricies, dignissim hendrerit purus. Nulla euismod dolor enim, vel fermentum ex ultricies ac. Donec aliquet vehicula egestas. Sed accumsan velit ac mauris porta, id imperdiet purus aliquam. Phasellus et faucibus erat. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Pellentesque vel nisl efficitur, euismod augue eu, consequat dui. Maecenas vestibulum tortor purus, egestas posuere tortor imperdiet eget. Nulla sit amet placerat diam.",
		},
	}

	for _, tt := range tests {
		grfFile, err := grf.Load(tt.FilePath)
		assert.NoError(t, err)

		t.Run(tt.Name, func(t *testing.T) {
			entry, err := grfFile.GetEntry(tt.EntryName)
			assert.NoError(t, err)
			assert.Equal(t, tt.ExpectedDataStr, string(entry.Data.Bytes()))
		})
	}
}
