package canvas

import (
	"strconv"
	"time"
)

func StrStrOrNil(strPtr *string) string {
	if strPtr == nil {
		return "nil"
	}

	return *strPtr
}

func IntStrOrNil(intPtr *int) string {
	if intPtr == nil {
		return "nil"
	}

	return strconv.Itoa(*intPtr)
}

func BoolStrOrNil(boolPtr *bool) string {
	if boolPtr == nil {
		return "nil"
	}

	if *boolPtr {
		return "true"
	} else {
		return "false"
	}
}

func TimeStrOrNil(timePtr *time.Time) string {
	if timePtr == nil {
		return "nil"
	}

	return timePtr.String()
}
