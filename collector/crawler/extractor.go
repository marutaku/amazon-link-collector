package crawler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/marutaku/amazon-link-collector/collector/utils"
)

// 取得したAmazonのURLを短くする
func ShortenAmazonLink(url string) (string, error) {
	host, err := utils.ExtractHostname(url)
	if err != nil {
		return "", err
	}
	if !strings.Contains(host, "www.amazon.co.jp") && !strings.Contains(host, "www.amazon.com") {
		return "", fmt.Errorf("link is not amazon")
	}
	// /dp or /gp/product で始まるURLを抽出
	re := regexp.MustCompile(`/dp/\w+|/gp/product/\w+`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 0 {
		return "https://" + host + matches[0], nil
	}
	// URL内部に10桁の数字or大文字アルファベットがあればそれを抽出
	re = regexp.MustCompile(`[0-9A-Z]{10}`)
	matches = re.FindStringSubmatch(url)
	if len(matches) > 0 {
		return "https://" + host + "/dp/" + matches[0], nil
	}
	fmt.Printf("link is invalid: %s\n", url)
	return "", nil
}

func ExtractAmazonLink(body string) ([]string, error) {
	re := regexp.MustCompile(`(\"|')https:\/\/www\.amazon(\.co\.jp\/|\.com)[\w%-\/=?]+(\"|')`)
	matches := re.FindAllString(body, -1)
	amazonLinks := []string{}
	for _, match := range matches {
		matchReplaced := strings.Replace(match, "\"", "", -1)
		matchReplaced = strings.Replace(matchReplaced, "'", "", -1)
		shortenLink, err := ShortenAmazonLink(matchReplaced)
		if err != nil {
			return nil, err
		}
		amazonLinks = append(amazonLinks, shortenLink)
	}
	return amazonLinks, nil
}
