package iterable

func FromString(s string) Iterable[rune] {
	var runes []rune
	for _, r := range s {
		runes = append(runes, r)
	}
	return FromSlice(runes)
}

func FromStringBytes(s string) Iterable[byte] {
	return IterableFun[byte](func() Iterator[byte] {
		pos := 0
		return Fun[byte](func() (byte, bool) {
			if pos >= len(s) {
				return 0, false
			}
			b := s[pos]
			pos++
			return b, true
		})
	})
}
