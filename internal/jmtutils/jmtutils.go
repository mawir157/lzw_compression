package jmtutils

import (
    "os"
    "strings"
)

const WHITELIST =
"#ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789 ,.-;:?'"

func Take(s string, n int) string {
    rs := []rune(s)
  return string(rs[:n])
}

func Drop(s string, n int) string {
    rs := []rune(s)
  return string(rs[n:])
}

func SanitiseInput(s string) string {
    sanitised := ""
    for _, r := range s {
        if strings.ContainsRune(WHITELIST, r) {
            sanitised += string(r)
        }
    }
    return sanitised
}

func InitLut() (lut map[string]int) {
    lut = make(map[string]int)
    for i, r := range WHITELIST {
        lut[string(r)] = i
    }

    return
}

func ValidateLut(lut map[string]int) bool {
    for _, r := range WHITELIST {
        if _, ok := lut[string(r)]; !ok {
            return false
        }
    }

    return true
}

func ReverseLookUp(lut map[string]int, value int) (key string, found bool) {
    for k, v := range lut {
        if (v == value) {
            return k, true
        }
    }
    return "", false
}

func ReadStrFile(fname string) (str string, err error) {
    b, err := os.ReadFile(fname)
    if err != nil {
        return "", err
    }

    s := string(b)
    s += "#"

    return SanitiseInput(s), nil
}
