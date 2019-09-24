package common

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

func JoinComma(i interface{}) string {
	return Join(ToGeneral(i), ", ")
}

func Join(li []interface{}, sep string) string {
	res := make([]string, len(li))
	for i := range res {
		str, ok := li[i].(fmt.Stringer)
		str1, ok2 := li[i].(string)
		if !ok && !ok2 {
			return ""
		}
		if ok {
			res[i] = str.String()
		}
		if ok2 {
			res[i] = str1
		}
	}
	return strings.Join(res, sep)
}
func ToGeneral(args interface{}) []interface{} {
	s := reflect.ValueOf(args)
	if s.Kind() != reflect.Slice {
		panic("ToGeneral given a non-slice type")
	}
	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}
	return ret
}

func Indent(length int, in string) string {
	buf := bytes.NewBufferString(in)
	var (
		res  bytes.Buffer
		err  error
		line string
	)
	for err == nil {
		line, err = buf.ReadString('\n')
		for i := 0; i < length; i++ {
			res.WriteRune(' ')
		}
		res.WriteString(line)
	}
	return res.String()
}
