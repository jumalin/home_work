package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

type unpackBuilder struct {
	*strings.Builder
}

const (
	Empty          = ""
	EncodingSymbol = 92
	OneSymbol      = 49
)

func Unpack(packedString string) (string, error) {
	packed := []rune(packedString)
	unpackedBuilder := unpackBuilder{&strings.Builder{}}
	for i := 0; i < len(packed); i++ {
		var err error

		if packed[i] == EncodingSymbol {
			err = writeEncoded(&i, &packed, &unpackedBuilder)
		} else {
			err = write(&i, &packed, &unpackedBuilder)
		}
		if err != nil {
			return Empty, err
		}
	}

	return unpackedBuilder.String(), nil
}

func write(i *int, packed *[]rune, unpackedBuilder *unpackBuilder) error {
	sym := (*packed)[*i]

	if unicode.IsDigit(sym) {
		return ErrInvalidString
	}

	if repeated(i, packed) {
		_, err := unpackedBuilder.writeRepeat(string(sym), (*packed)[*i+1])
		*i++
		return err
	}

	_, err := unpackedBuilder.write(string(sym))
	return err
}

func writeEncoded(i *int, packed *[]rune, unpackedBuilder *unpackBuilder) error {
	size := len(*packed)

	if nothingToEncode(i, size) {
		return ErrInvalidString
	}

	sym := (*packed)[*i+1]
	if cannotRepeatEncoded(i, size) {
		_, err := unpackedBuilder.write(string(sym))
		*i++
		return err
	}

	if count := (*packed)[*i+2]; encodedRepeated(count) {
		_, err := unpackedBuilder.writeRepeat(string(sym), count)
		*i += 2
		return err
	}

	_, err := unpackedBuilder.write(string(sym))
	*i++

	return err
}

func (b *unpackBuilder) writeRepeat(literalToAdd string, literalCount rune) (int, error) {
	count, _ := strconv.Atoi(string(literalCount))
	return b.WriteString(strings.Repeat(literalToAdd, count))
}

func (b *unpackBuilder) write(literalToAdd string) (int, error) {
	return b.writeRepeat(literalToAdd, OneSymbol)
}

func repeated(i *int, packed *[]rune) bool {
	return *i+1 < len(*packed) && unicode.IsDigit((*packed)[*i+1])
}

func encodedRepeated(count rune) bool {
	return unicode.IsDigit(count)
}

func cannotRepeatEncoded(i *int, size int) bool {
	return *i+2 == size
}

func nothingToEncode(i *int, size int) bool {
	return *i == size-1
}
