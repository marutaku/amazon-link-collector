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
	re := regexp.MustCompile(`/dp/\w+|/gp/product/\w+`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 0 {
		return "https://" + host + matches[0], nil
	} else {
		return "", fmt.Errorf("link is invalid")
	}
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
