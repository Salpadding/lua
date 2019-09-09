package ast

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

func joinComma(i interface{}) string {
	return join(toGeneral(i), ", ")
}

func join(li []interface{}, sep string) string {
	res := make([]string, len(li))
	for i := range res {
		str, ok := li[i].(fmt.Stringer)
		if !ok {
			return ""
		}
		res[i] = str.String()
	}
	return strings.Join(res, sep)
}
func toGeneral(args interface{}) []interface{} {
	s := reflect.ValueOf(args)
	if s.Kind() != reflect.Slice {
		panic("toGeneral given a non-slice type")
	}
	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}
	return ret
}

func indent(length int, in string) string {
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
