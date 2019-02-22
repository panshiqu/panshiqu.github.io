package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path"
)

var files = []string{"010_1.jpeg", "017_1.png", "024_1.png", "029_1.png", "034_1.png", "039_1.png", "011_1.jpeg", "020_1.png", "025_1.png", "030_1.png", "035_1.png", "040_1.png", "012_1.jpeg", "021_1.png", "026_1.jpeg", "031_1.png", "036_1.png", "041_1.png", "012_2.jpeg", "022_1.png", "027_1.png", "032_1.png", "037_1.png", "042_1.jpeg", "012_3.jpeg", "023_1.png", "028_1.jpeg", "033_1.jpeg", "038_1.jpeg", "066_1.png"}

func main() {
	for _, v := range files {
		f, err := os.Open(v)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		_, ext, err := image.Decode(f)
		if err != nil {
			log.Fatal(err)
		}

		if path.Ext(v)[1:] == ext {
			continue
		}

		fmt.Println(v, ":", ext)
	}
}
