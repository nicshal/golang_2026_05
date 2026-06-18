package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	mlt := -1
	var currentRune rune
	var err error

	sb := strings.Builder{}
	for _, rn := range str {
		if unicode.IsDigit(rn) {
			mlt, err = getMlt(mlt, currentRune, rn)
			if err != nil {
				return "", ErrInvalidString
			}
		} else {
			add(&mlt, currentRune, &sb)
			currentRune = rn
			mlt = -1
		}
	}

	add(&mlt, currentRune, &sb)
	return sb.String(), nil
}

func add(mlt *int, currentRune rune, sb *strings.Builder) {
	if currentRune == 0 {
		return
	}
	if *mlt < 0 {
		*mlt = 1
	}
	(*sb).WriteString(strings.Repeat(string(currentRune), *mlt))
}

func getMlt(mlt int, currentRune rune, rn rune) (int, error) {
	if mlt < 0 {
		if currentRune == 0 {
			return 0, ErrInvalidString
		}
		mlt, _ = strconv.Atoi(string(rn))
		return mlt, nil
	}
	return 0, ErrInvalidString
}
