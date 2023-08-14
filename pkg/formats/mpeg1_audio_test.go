package formats

import (
	"testing"

	"github.com/pion/rtp"
	"github.com/stretchr/testify/require"
)

func TestMPEG1AudioAttributes(t *testing.T) {
	format := &MPEG1Audio{}
	require.Equal(t, "MPEG-1/2 Audio", format.Codec())
	require.Equal(t, 90000, format.ClockRate())
	require.Equal(t, true, format.PTSEqualsDTS(&rtp.Packet{}))
}

func TestMPEG1AudioDecEncoder(t *testing.T) {
	format := &MPEG1Audio{}

	enc, err := format.CreateEncoder()
	require.NoError(t, err)

	pkts, err := enc.Encode([][]byte{{
		0xff, 0xfb, 0x14, 0x64, 0x00, 0x0f, 0xf0, 0x00,
		0x00, 0x69, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00,
		0x0d, 0x20, 0x00, 0x00, 0x01, 0x00, 0x00, 0x01,
		0xa4, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x34,
		0x80, 0x00, 0x00, 0x04, 0x4c, 0x41, 0x4d, 0x45,
		0x33, 0x2e, 0x31, 0x30, 0x30, 0x55, 0x55, 0x55,
		0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
		0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
		0x55, 0xc0, 0x65, 0xf4, 0xa0, 0x31, 0x8f, 0xce,
		0x8d, 0x46, 0xfc, 0x8c, 0x73, 0xb9, 0x34, 0x3e,
		0xb5, 0x03, 0x39, 0xc0, 0x04, 0x01, 0x98, 0x44,
		0x38, 0xe0, 0x98, 0x10, 0x9b, 0xa8, 0x0f, 0xa8,
	}}, 0)
	require.NoError(t, err)
	require.Equal(t, format.PayloadType(), pkts[0].PayloadType)

	dec, err := format.CreateDecoder()
	require.NoError(t, err)

	byts, _, err := dec.Decode(pkts[0])
	require.NoError(t, err)
	require.Equal(t, [][]byte{{
		0xff, 0xfb, 0x14, 0x64, 0x00, 0x0f, 0xf0, 0x00,
		0x00, 0x69, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00,
		0x0d, 0x20, 0x00, 0x00, 0x01, 0x00, 0x00, 0x01,
		0xa4, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x34,
		0x80, 0x00, 0x00, 0x04, 0x4c, 0x41, 0x4d, 0x45,
		0x33, 0x2e, 0x31, 0x30, 0x30, 0x55, 0x55, 0x55,
		0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
		0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55,
		0x55, 0xc0, 0x65, 0xf4, 0xa0, 0x31, 0x8f, 0xce,
		0x8d, 0x46, 0xfc, 0x8c, 0x73, 0xb9, 0x34, 0x3e,
		0xb5, 0x03, 0x39, 0xc0, 0x04, 0x01, 0x98, 0x44,
		0x38, 0xe0, 0x98, 0x10, 0x9b, 0xa8, 0x0f, 0xa8,
	}}, byts)
}
