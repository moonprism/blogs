package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type artInfo struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

var initFlag = flag.Bool("i", false, "强制同步所有文章")
var articleId = flag.Int("a", 0, "指定更新某文章")

func main() {
	flag.Parse()


	articlesBody, err := request("https://kicoe.com/api/articles")
	handleError(err)

	articles := make(map[string][]artInfo)

	err = json.Unmarshal(articlesBody, &articles)
	handleError(err)

	if *articleId != 0 {
		for date, arts := range articles {
			for _, art := range arts {
				if art.Id == *articleId {
					mdFileName := date+"/"+art.Title+".md"
					err := syncArticle(art.Id, mdFileName)
					handleError(err)
				}
			}
		}
		return
	}

	keys := make([]string, 0, len(articles))
	for k := range articles {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	var readme string

	for _, date := range keys {
		// 生成归档markdown & 解析文章内容
		readme += fmt.Sprintf("### %s\n\n", date)
		if (!exists(date)) {
			os.Mkdir(date, os.ModePerm)
		}
		for _, art := range articles[date] {
			url := strings.ReplaceAll(url.QueryEscape(art.Title), "+", "%20")
			readme += fmt.Sprintf("* [%s](https://github.com/moonprism/blogs/blob/master/%s/%s.md)\n", art.Title, date, url)
			mdFileName := date+"/"+art.Title+".md"
			if *initFlag || !exists(mdFileName) {
				err := syncArticle(art.Id, mdFileName)
				handleError(err)
			}
		}
		readme += "\n"
	}

	file, err := os.OpenFile("README.md", os.O_TRUNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	handleError(err)
	defer file.Close()

	_, err = file.Write([]byte(readme))
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func syncArticle(artId int, fileName string) error {
	println("start sync article "+strconv.Itoa(artId)+" to "+fileName)
	articleBody, err := request(fmt.Sprintf("https://kicoe.com/api/article/%d", artId))
	handleError(err)

	r, err := regexp.Compile(`([^\\]|^)!\[(.*?)\]\((.*?)\)`)
	// TODO cdn备份博客图片
	rep := []byte("${1}![${2}](https://kicoe-blog.oss-cn-shanghai.aliyuncs.com/${3})")
	articleBody = r.ReplaceAll(articleBody, rep)

	file, err := os.OpenFile(fileName, os.O_TRUNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	handleError(err)
	defer file.Close()

	_, err = file.Write(articleBody)
	return err
}

func request(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
