package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	url := `アイウエオ"https://www.amazon.co.jp/%E3%82%AA%E3%83%BC%E3%83%87%E3%82%A3%E3%82%AA%E3%83%95%E3%82%A1%E3%83%B3-LAN%E3%82%B1%E3%83%BC%E3%83%96%E3%83%AB-Cat6-%E3%82%AE%E3%82%AC%E3%83%93%E3%83%83%E3%83%88%E5%AF%BE%E5%BF%9C-%E3%83%95%E3%83%A9%E3%83%83%E3%83%88%E3%82%BF%E3%82%A4%E3%83%97/dp/B086VY2N25/ref=sr_1_9?keywords=%E7%9F%AD%E3%81%84lan%E3%82%B1%E3%83%BC%E3%83%96%E3%83%AB&qid=1703692884&sr=8-9&th=1"アイウエオ`
	if strings.Contains(url, "www.amazon.co.jp") || strings.Contains(url, "www.amazon.com") {
		re := regexp.MustCompile(`(\"|')https:\/\/www\.amazon(\.co\.jp\/|\.com)[\w%-\/=?]+(\"|')`)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 0 {
			fmt.Println(matches[0])
		} else {
			fmt.Println("Amazonの商品ページで実行してください")
		}
	} else {
		fmt.Println("Amazonの商品ページで実行してください")
	}
}
