package str

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

// Empty 字符串是否为空
func Empty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// SubString  begin从0开始，截取字符串，只支持英文和数字
func SubString(str string, begin, length int) string {
	lth := len(str)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length

	if end > lth {
		end = lth
	}
	return str[begin:end]
}

// SubStringFull  begin从0开始，截取字符串，支持中文
func SubStringFull(str string, begin, length int) string {
	rs := []rune(str)
	lth := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length

	if end > lth {
		end = lth
	}
	return string(rs[begin:end])
}

// Ucfirst 字符串第一个字母大写
func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// HexReverse 16进制大小端反转。1234ABCD --> CDAB3412
func HexReverse(hex string) string {
	str := ""
	rHex := []rune(hex)
	hexLen := len(rHex)
	for i := 0; i < hexLen-1; i += 2 {
		str = fmt.Sprintf("%c%c%s", rHex[i], rHex[i+1], str)
	}

	return str
}

// Hex2Dec  十六进制字符串转十进制数字（INT）
func Hex2Dec(hex string) int {
	i, _ := strconv.ParseInt(hex, 16, 64)
	return int(i)
}

// Hex2DecU  十六进制字符串转十进制数字（无符号型）
func Hex2DecU(hex string) uint {
	u, _ := strconv.ParseUint(hex, 16, 64)
	return uint(u)
}

// String2Int ，字符串转数字，忽略错误，如有错，返回0
func String2Int(s string) int32 {
	i, err := strconv.Atoi(s)
	if err != nil {
		i = 0
	}
	return int32(i)
}

// String2Int64  字符串转数据
func String2Int64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		i = 0
	}
	return i
}

// String2UInt ，字符串转成UINT，有错返回0
func String2UInt(s string) uint32 {
	return uint32(String2Int(s))
}

// FormatHexData 格式化16进制，两个字符串加一个空格
func FormatHexData(str string) string {
	len := len(str)
	var buffer bytes.Buffer
	for j := 0; j <= len; j += 2 {
		buffer.WriteString(SubString(str, j, 2))
		buffer.WriteString(" ")
	}
	s := strings.TrimSpace(buffer.String())
	return s
}

// Convert2U32 interface转uint32
func Convert2U32(i interface{}) uint32 {
	var num uint32 = 0
	switch i.(type) {
	case float64:
		num = uint32(i.(float64))
	case string:
		num = String2UInt(i.(string))
	case int64:
		num = uint32(i.(int64))
	case int:
		num = uint32(i.(int))
	default:
		break
	}
	return num
}

// Convert2Int32 interface转int32
func Convert2Int32(i interface{}) int32 {
	var num int32 = 0
	switch i.(type) {
	case float64:
		num = int32(i.(float64))
	case string:
		num = String2Int(i.(string))
	case int64:
		num = int32(i.(int64))
	case int:
		num = int32(i.(int))
	default:
		break
	}
	return num
}

// Md5 32位MD5
func Md5(strs string) string {
	w := md5.New()
	io.WriteString(w, strs)
	//将str写入到w中
	return fmt.Sprintf("%x", w.Sum(nil))
}

// U32toString  uin32转string
func U32toString(num uint32) string {
	return strconv.Itoa(int(num))
}

func String2UInt64(s string) uint64 {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		i = 0
	}
	return i
}

func Convert2U64(i interface{}) uint64 {
	var num uint64 = 0
	switch i.(type) {
	case float64:
		num = uint64(i.(float64))
	case string:
		num = String2UInt64(i.(string))
	case int64:
		num = uint64(i.(int64))
	case int:
		num = uint64(i.(int))
	default:
		break
	}
	return num
}
