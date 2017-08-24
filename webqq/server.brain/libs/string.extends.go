package libs

//start：正数 - 在字符串的指定位置开始,超出字符串长度强制把start变为字符串长度
//       负数 - 在从字符串结尾的指定位置开始
//       0 - 在字符串中的第一个字符处开始
//length:正数 - 从 start 参数所在的位置返回
//       负数 - 从字符串末端返回

func Substr(str string, start, length int) string {
    if length == 0 {
        return ""
    }
    rune_str := []rune(str)
    len_str := len(rune_str)

    if start < 0 {
        start = len_str + start
    }
    if start > len_str {
        start = len_str
    }
    end := start + length
    if end > len_str {
        end = len_str
    }
    if length < 0 {
        end = len_str + length
    }
    if start > end {
        start, end = end, start
    }
    return string(rune_str[start:end])
}
