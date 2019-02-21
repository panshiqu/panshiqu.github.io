package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func getContent(f1, f2 []byte, url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	n := bytes.Index(body, f1)
	if n == -1 {
		log.Fatal("n==-1")
	}
	var buf bytes.Buffer
	for i := n; i < len(body); i++ {
		if body[i] == '<' {
			if bytes.Compare(body[i:i+len(f2)], f2) == 0 {
				break
			}
			r := bytes.IndexByte(body[i:], '>')
			if r == -1 {
				log.Fatal("r==-1")
			}
			i += r
		} else {
			l := bytes.IndexByte(body[i:], '<')
			if l == -1 {
				log.Fatal("l==-1")
			}
			buf.Write(body[i : i+l])
			i += l - 1
		}
	}
	return buf.Bytes()
}

var m = map[string]string{
	"https://panshiqu.github.io/blog/072.html": "https://blog.csdn.net/panshiqu/article/details/83007702",
	"https://panshiqu.github.io/blog/071.html": "https://blog.csdn.net/panshiqu/article/details/82878303",
	"https://panshiqu.github.io/blog/070.html": "https://blog.csdn.net/panshiqu/article/details/81978291",
	"https://panshiqu.github.io/blog/069.html": "https://blog.csdn.net/panshiqu/article/details/79179148",
	"https://panshiqu.github.io/blog/068.html": "https://blog.csdn.net/panshiqu/article/details/78675617",
	"https://panshiqu.github.io/blog/067.html": "https://blog.csdn.net/panshiqu/article/details/74572133",
	"https://panshiqu.github.io/blog/066.html": "https://blog.csdn.net/panshiqu/article/details/70160325",
	"https://panshiqu.github.io/blog/064.html": "https://blog.csdn.net/panshiqu/article/details/58610958",
	"https://panshiqu.github.io/blog/063.html": "https://blog.csdn.net/panshiqu/article/details/55049159",
	"https://panshiqu.github.io/blog/062.html": "https://blog.csdn.net/panshiqu/article/details/53907204",
	"https://panshiqu.github.io/blog/061.html": "https://blog.csdn.net/panshiqu/article/details/53788067",
	"https://panshiqu.github.io/blog/060.html": "https://blog.csdn.net/panshiqu/article/details/53393693",
	"https://panshiqu.github.io/blog/059.html": "https://blog.csdn.net/panshiqu/article/details/53407649",
	"https://panshiqu.github.io/blog/058.html": "https://blog.csdn.net/panshiqu/article/details/53349967",
	"https://panshiqu.github.io/blog/057.html": "https://blog.csdn.net/panshiqu/article/details/53217523",
	"https://panshiqu.github.io/blog/056.html": "https://blog.csdn.net/panshiqu/article/details/53202989",
	"https://panshiqu.github.io/blog/055.html": "https://blog.csdn.net/panshiqu/article/details/53156990",
	"https://panshiqu.github.io/blog/054.html": "https://blog.csdn.net/panshiqu/article/details/53116535",
	"https://panshiqu.github.io/blog/053.html": "https://blog.csdn.net/panshiqu/article/details/53037373",
	"https://panshiqu.github.io/blog/052.html": "https://blog.csdn.net/panshiqu/article/details/53008575",
	"https://panshiqu.github.io/blog/051.html": "https://blog.csdn.net/panshiqu/article/details/52934554",
	"https://panshiqu.github.io/blog/050.html": "https://blog.csdn.net/panshiqu/article/details/52779627",
	"https://panshiqu.github.io/blog/049.html": "https://blog.csdn.net/panshiqu/article/details/52749117",
	"https://panshiqu.github.io/blog/048.html": "https://blog.csdn.net/panshiqu/article/details/52748837",
	"https://panshiqu.github.io/blog/047.html": "https://blog.csdn.net/panshiqu/article/details/52681948",
}

func main() {
	for k, v := range m {
		i, err := strconv.Atoi(k[len(k)-8 : len(k)-5])
		if err != nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile(fmt.Sprintf("cmp1/%d", i), getContent([]byte("<hr />"), []byte("<footer"), k), 0644); err != nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile(fmt.Sprintf("cmp2/%d", i), getContent([]byte("<div id=\"content_views\" class=\"markdown_views prism-atom-one-dark\">"), []byte("</article>"), v), 0644); err != nil {
			log.Fatal(err)
		}
	}
}
