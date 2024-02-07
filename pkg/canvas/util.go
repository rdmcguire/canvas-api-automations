package canvas

import (
	"fmt"
	"time"
)

func StrOrNil[T any](v *T) string {
	if v == nil {
		return "nil"
	}

	switch vT := any(*v).(type) {
	case float32, float64:
		return fmt.Sprintf("%.4f", vT)

	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", vT)

	case bool:
		if vT {
			return "true"
		}
		return "false"

	case string:
		return vT

	case time.Time:
		return vT.String()
	}

	return "unsupported type ptr"
}
