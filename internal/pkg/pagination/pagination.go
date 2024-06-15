package pagination

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ParsePageToken(pageToken string) (lastValue int64) {
	pToken := strings.Split(pageToken, "_")

	if len(pToken) != 2 {
		return
	}

	sortValue := pToken[1]
	lastValue, _ = strconv.ParseInt(sortValue, 10, 64)

	return
}

func CreatePageToken(data interface{}, limit int) (nextToken string) {
	voData := reflect.ValueOf(data)
	if voData.Kind() != reflect.Slice {
		return
	}

	if !(voData.IsValid()) {
		return
	}

	if voData.Len() < limit {
		return
	}

	lastData := voData.Index(voData.Len() - 1)

	nextLastValue := lastData.FieldByName("UpdatedAt")

	nextToken = fmt.Sprintf("%s_%d", "last", nextLastValue)

	return
}
