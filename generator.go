package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	fs, err := ioutil.ReadDir("csdn/list")
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer

	buf.WriteString(`---
layout: default
---

# blog list

* * *

## 2018

`)

	index := 72 // 共72篇博客

	year := "2018" // 最新为2018年写成

	for _, v := range fs {
		if v.IsDir() {
			continue
		}

		f, err := ioutil.ReadFile(fmt.Sprintf("csdn/list/%s", v.Name()))
		if err != nil {
			log.Fatal(err)
		}

		for {
			n1 := bytes.Index(f, []byte(`title="编辑">`))
			if n1 == -1 {
				break
			}

			n2 := bytes.IndexByte(f[n1:], '<')
			if n2 == -1 {
				break
			}

			n3 := bytes.Index(f[n1+n2:], []byte(`class="">201`))
			if n3 == -1 {
				break
			}

			n4 := bytes.IndexByte(f[n1+n2+n3:], '<')
			if n4 == -1 {
				break
			}

			t := f[n1+n2+n3+9 : n1+n2+n3+n4]

			if y := string(t[:4]); year != y {
				buf.WriteString(fmt.Sprintf("## %s\n\n", y))
				year = y
			}

			if err := ioutil.WriteFile(fmt.Sprintf("blog/%03d.md", index), []byte(fmt.Sprintf(`---
layout: default
---

# %s
_%s_

* * *

content
`, f[n1+15:n1+n2], fmt.Sprintf("%s-%s-%s %s:%s:%s", t[:4], t[7:9], t[12:14], t[18:20], t[21:23], t[24:]))), 0644); err != nil {
				log.Fatal(err)
			}

			buf.WriteString(fmt.Sprintf("[%s](./blog/%03d.html).\n\n", f[n1+15:n1+n2], index))

			f = f[n1+n2+n3+n4:]

			if l := bytes.Index(f, []byte(`<a href="`)); l != -1 {
				if r := bytes.IndexByte(f[l+9:], '"'); r != -1 {
					fmt.Printf("\"https://panshiqu.github.io/blog/%03d.html\": \"%s\",\n", index, f[l+9:l+9+r])
				}
			}

			index--

			if index == 65 {
				index = 64 // 回收站占位
			}
		}
	}

	if err := ioutil.WriteFile("index.md", buf.Bytes()[:buf.Len()-1], 0644); err != nil {
		log.Fatal(err)
	}
}
