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
		"nested": map[string]interface{}{
			"egg": map[string]interface{}{
				"name":  "Eggbert",
				"value": 0.5,
			},
			"ham": map[string]interface{}{
				"name":  "Hambert",
				"value": 0.75,
			},
		},
		"listTest (compound)": []map[string]interface{}{
			{
				"createdOn": int64(1264099775885),
				"name":      "Compound Tag #0",
			},
			{
				"createdOn": int64(1264099775885),
				"name":      "Compound Tag #1",
			},
		},
		"empty": []int64{},
		"listTest (long)": []int64{
			11,
			12,
			13,
			14,
			15,
		},
		"byteTest":   byte(127),
		"shortTest":  int16(32767),
		"intTest":    int32(2147483647),
		"longTest":   int64(9223372036854775807),
		"byteArr":    []byte{1, 2, 3, 4, 5}, // Will be printed as base64 encoded because thats how JSON works
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

	// Decode compressed file stream
	err := nbt.DecodeCompressedStream(f, &m)
	if err != nil {
		panic(err)
	}

	// Pretty print the map using JSON
	t, _ := json.MarshalIndent(m, "", "  ")
	fmt.Print(string(t))
}
