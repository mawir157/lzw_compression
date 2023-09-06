package main

import (
    "errors"
    "fmt"
    "strings"
    "time"
    "unsafe"
)

import utils "avos/lzw/internal/jmtutils"

func Encode(s string) (encoded []int, lut map[string]int, err error) {
    if (!strings.ContainsRune(s, '#')) {
        return nil, nil, errors.New("No escape character in string")
    }

    lut = utils.InitLut()
    if !utils.ValidateLut(lut) {
        return nil, nil, errors.New("Something went wrong")
    }
    encoded = make([]int, 0)

    loop_string := s
    current_value := 0
    for ;(len(loop_string) != 0); {
        // if the next rune is '#' we're done
        if (utils.Take(loop_string, 1) == "#") {
            encoded = append(encoded, lut["#"])
            break
        }
        for n := 1; n <= len(loop_string); n++ {
            head := utils.Take(loop_string, n)
            if val, ok := lut[head]; ok {
                // this string is in the LUT, so cache it and try to find a
                // longer one
                current_value = val
            } else { // if is not in the lut we are done...
                // ...so add the encoded value to the slice...
                encoded = append(encoded, current_value)
                // ...and add the new string to the lut...
                lut[head] = len(lut)
                // and drop the word from the start of the string
                loop_string = utils.Drop(loop_string, (n - 1))
                break
            }
        }
    }

    return encoded, lut, nil
}

func DecodeWithoutLut(encoded []int) (string, error) {
    lut := utils.InitLut()

    return Decode(encoded, lut)
}

func Decode(encoded []int, lut map[string]int) (string, error) {
    if !utils.ValidateLut(lut) {
        return "", errors.New("Something went wrong")
    }

    if encoded[0] > len(utils.WHITELIST) {
        return "", errors.New("Something went wrong")
    }

    conjecture := ""
    var decoded strings.Builder

    for _, v := range encoded {
        if r, ok := utils.ReverseLookUp(lut, v); ok {
            // update my conjecture
            conjecture += utils.Take(r, 1)
            decoded.WriteString(r)
            // if we have not seen this conecture before
            if _, ok := lut[conjecture]; !ok {
                lut[conjecture] = len(lut)
                conjecture = r
            }
        } else {
            // if we get here, then value is missing from the look up table this
            // means that the value is conjecture + the first rune of conjecture
            conjecture += utils.Take(conjecture, 1)
            decoded.WriteString(conjecture)
            lut[conjecture] = len(lut)
        }
    }

    return decoded.String(), nil
}

func main() {
    start := time.Now()
    test_string, err := utils.ReadStrFile("../moby_dick.txt")
    elapsed := time.Since(start)

//     test_string := utils.SanitiseInput(
// `Call me Ishmael. Some years ago-never mind how long precisely-having little or 
// no money in my purse, and nothing particular to interest me on shore, I thought 
// I would sail about a little and see the watery part of the world. It is a way I 
// have of driving off the spleen and regulating the circulation. Whenever I find 
// myself growing grim about the mouth; whenever it is a damp, drizzly November in 
// my soul; whenever I find myself involuntarily pausing before coffin warehouses, 
// and bringing up the rear of every funeral I meet; and especially whenever my 
// hypos get such an upper hand of me, that it requires a strong moral principle to 
// prevent me from deliberately stepping into the street, and methodically knocking 
// people's hats off-then, I account it high time to get to sea as soon as I can. 
// This is my substitute for pistol and ball. With a philosophical flourish Cato 
// throws himself upon his sword; I quietly take to the ship. There is nothing 
// surprising in this. If they but knew it, almost all men in their degree, 
// some time or other, cherish very nearly the same feelingstowards the ocean with 
// me.#`)

    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("File read, time elapsed", elapsed)
    }

    start = time.Now()
    e, _, err := Encode(test_string)
    elapsed = time.Since(start)
    if err != nil {
        fmt.Println(err)
        return
    } else {
        fmt.Println("File encoded, time elapsed", elapsed)
    }

    start = time.Now()
    s, err := DecodeWithoutLut(e)
    elapsed = time.Since(start)
    if err == nil {
        fmt.Println("File decoded, time elapsed", elapsed)
        if (s == test_string) {
            fmt.Println("File decoded correctly")
        } else {
            fmt.Println("[WARNING] file decoded INCORRECTLY")
        }
        
    } else {
        fmt.Println(err)
        return
    }

    string_bytes := len(test_string) * int(unsafe.Sizeof(test_string))
    fmt.Println("Original string", string_bytes, "bytes")

    code_bytes := len(e) * int(unsafe.Sizeof(e))
    fmt.Println("Compressed string", code_bytes, "bytes")

    fmt.Println("Compression factor", float64(string_bytes) / float64(code_bytes))
}
