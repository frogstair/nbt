# NBT for Go

This is a simple NBT encoder/decoder that does not use any custom structs/types that supports compression

## Supported types of data

For now it can only encode/decode tags whose root node is a **composite**, other types are coming in the near future. It can also only encode/decode maps, and support for structs and struct tags will also be added later

## Features

- Easy to represent data with map syntax
- Compression with gzip
- Fast
- No external libraries
- Very compact
- A more compact wrapper to improve readability

## Precautions

It is necessary to explicitly specify a type for each **numerical value** because of how NBT works.

## Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"

    "github.com/frogstair/nbt"
)

func main() {

    // Sample data to work with
    data := map[string]interface{}{
        "nested": nbt.C{
            "egg": nbt.C{
                "name":  "Eggbert",
                "value": 0.5,
            },
            "ham": nbt.C{
                "name":  "Hambert",
                "value": 0.75,
            },
        },
        "listTest (compound)": []nbt.C{
            {
                "createdOn": int64(1264099775885),
                "name":      "nbt.Compound Tag #0",
            },
            {
                "createdOn": int64(1264099775885),
                "name":      "nbt.Compound Tag #1",
            },
        },
        "empty": []int64{},
        // No need to specify type here, because go will force every number to be an int64
        "listTest (long)": []int64{
            11,
            12,
            13,
            14,
            15,
        },
        "byteTest":  byte(127),
        "shortTest": int16(32767),
        "intTest":   int32(2147483647),
        "longTest":  int64(9223372036854775807),
        // Will be printed as base64 encoded because thats how JSON works
        "byteArr": []byte{1, 2, 3, 4, 5},
        // Support for unicode
        "stringTest": "HELLO WORLD THIS IS A TEST STRING! こんにちは世界〜",
    }

    // Encode the data as bytes and compress it with gzip
    b := nbt.EncodeCompress(data, "")
    // Write data to a file
    ioutil.WriteFile("./bigtest.nbt", b, 0644)

    // Open the file with the data
    f, _ := os.Open("bigtest.nbt")
    // Make an empty container
    m := make(map[string]interface{})

    // Decode compressed file stream. Any reader that contains the NBT data works
    err := nbt.DecodeCompressedStream(f, &m)
    if err != nil {
        panic(err)
    }

    // Pretty print the map using JSON
    t, _ := json.MarshalIndent(m, "", "  ")
    fmt.Print(string(t), "\n")
}


```

This example can also be found in the `examples` folder
