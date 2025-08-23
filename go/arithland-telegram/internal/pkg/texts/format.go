package texts

import (
	"fmt"
	"github.com/mshafiee/jalali"
	"strings"
	"time"
)

func Format(template string, data map[string]string) string {
	args := make([]string, 0, len(data)*2)
	for o, n := range data {
		args = append(args, fmt.Sprintf("{%s}", o), n)
	}

	return strings.NewReplacer(args...).Replace(template)
}

func FormatFloat(num float64) string {
	return ReplaceDigitsToFarsi(fmt.Sprintf("%.2f", num))
}

func FormatInt(num int64) string {
	return ReplaceDigitsToFarsi(fmt.Sprintf("%d", num))
}

func FormatBoolAsEmoji(v bool) string {
	if v {
		return "✅"
	}

	return "❌"
}

func FormatTime(t time.Time) string {
	return ReplaceDigitsToFarsi(jalali.ToJalali(t).In(jalali.Tehran()).Format("%d %B %Y، ساعت %H:%M"))
}
