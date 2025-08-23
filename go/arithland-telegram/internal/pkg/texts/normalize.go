package texts

import (
	"github.com/lovelydeng/gomoji"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"unicode"
)

var (
	punctuationsRegexp  = regexp.MustCompile(`[\p{P}\p{S}]`)
	whitespacesRegexp   = regexp.MustCompile(`[\p{Zs}\p{Zl}\p{Zp}\x{200B}\x{200C}\x{200D}\x{2060}]+`)
	similarCharReplacer = strings.NewReplacer(
		"ي", "ی",
		"ك", "ک",
		"ة", "ه",
		"أ", "ا",
		"إ", "ا",
		"٤", "۴",
		"٥", "۵",
		"٦", "۶",
	)
	digitToEnglishReplacer = strings.NewReplacer(
		// Farsi
		"۰", "0",
		"۱", "1",
		"۲", "2",
		"۳", "3",
		"۴", "4",
		"۵", "5",
		"۶", "6",
		"۷", "7",
		"۸", "8",
		"۹", "9",

		// Arabic
		"٠", "0",
		"١", "1",
		"٢", "2",
		"٣", "3",
		"٤", "4",
		"٥", "5",
		"٦", "6",
		"٧", "7",
		"٨", "8",
		"٩", "9",
	)
	digitToFarsiReplacer = strings.NewReplacer(
		// English
		"0", "۰",
		"1", "۱",
		"2", "۲",
		"3", "۳",
		"4", "۴",
		"5", "۵",
		"6", "۶",
		"7", "۷",
		"8", "۸",
		"9", "۹",

		// Arabic
		"٠", "۰",
		"١", "۱",
		"٢", "۲",
		"٣", "۳",
		"٤", "۴",
		"٥", "۵",
		"٦", "۶",
		"٧", "۷",
		"٨", "۸",
		"٩", "۹",
	)
)

func NormalizeTitle(s string) string {
	s = RemoveDiacritics(s)
	s = RemovePunctuations(s)
	s = RemoveEmojis(s)
	s = UnifyWhitespaces(s)
	s = UnifySimilarChars(s)
	s = ReplaceDigitsToFarsi(s)
	s = TrimSpace(s)

	return s
}

func NormalizeDescription(s string) string {
	s = RemoveDiacritics(s)
	s = RemoveEmojis(s)
	s = UnifyWhitespaces(s)
	s = UnifySimilarChars(s)
	s = ReplaceDigitsToFarsi(s)
	s = TrimSpace(s)

	return s
}

func NormalizeValue(s string) string {
	s = RemoveDiacritics(s)
	s = RemovePunctuations(s)
	s = RemoveEmojis(s)
	s = UnifyWhitespaces(s)
	s = UnifySimilarChars(s)
	s = ReplaceDigitsToEnglish(s)
	s = TrimSpace(s)

	return s
}

func RemoveDiacritics(s string) string {
	s, _, _ = transform.String( //nolint: errcheck
		transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC), s,
	)

	return s
}

func RemovePunctuations(s string) string {
	return punctuationsRegexp.ReplaceAllString(s, "")
}

func RemoveEmojis(s string) string {
	return gomoji.RemoveEmojis(s)
}

func UnifyWhitespaces(s string) string {
	return whitespacesRegexp.ReplaceAllString(s, " ")
}

func UnifySimilarChars(s string) string {
	return similarCharReplacer.Replace(s)
}

func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

func ReplaceDigitsToEnglish(s string) string {
	return digitToEnglishReplacer.Replace(s)
}

func ReplaceDigitsToFarsi(s string) string {
	return digitToFarsiReplacer.Replace(s)
}
