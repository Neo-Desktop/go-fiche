package main

import (
	"fmt"
	"strconv"
)

// slugMap contains a sequence of valid characters for slug generation
const slugMap = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"

// generateSlug takes an integer and generates a string based on slugMap
func generateSlug(seed int64) string {
	stringSeed := fmt.Sprintf("%d", seed)
	evenLength := (len(stringSeed) / 2) * 2
	digitHold := int64(0)
	out := ""
	for i := 0; i < evenLength; i += 2 {
		digitHold, _ = strconv.ParseInt(stringSeed[i:i+2], 10, 8)
		digitHold = digitHold % int64(len(slugMap))
		out += slugMap[digitHold : digitHold+1]
	}
	return out
}
