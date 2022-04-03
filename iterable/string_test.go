package iterable

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFromString(t *testing.T) {
	require.Equal(t, []rune{'你', '好'}, ToSlice(FromString("你好")))
	require.Equal(t, []rune{'👩', 8205, '💻'}, ToSlice(FromString("👩‍💻")))
}

func TestFromStringBytes(t *testing.T) {
	require.Equal(t, []byte{0xe4, 0xbd, 0xa0, 0xe5, 0xa5, 0xbd}, ToSlice(FromStringBytes("你好")))
	require.Equal(t, []byte{0xf0, 0x9f, 0x91, 0xa9, 0xe2, 0x80, 0x8d, 0xf0, 0x9f, 0x92, 0xbb}, ToSlice(FromStringBytes("👩‍💻")))

}
