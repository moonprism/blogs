package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type apiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type artInfoData struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type artInfoResponse struct {
	apiResponse
	Data        map[string][]artInfoData `json:"data"`
}

var initFlag = flag.Bool("i", false, "强制同步所有文章")
var articleId = flag.Int("a", 0, "指定更新某文章")

func main() {
	flag.Parse()

	content, err := request("https://kicoe.com/api/articles")
	handleError(err)

	var resp artInfoResponse
	err = json.Unmarshal(content, &resp)
	handleError(err)

	articles := resp.Data

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

	var resp apiResponse
	err = json.Unmarshal(articleBody, &resp)
	handleError(err)

	r, err := regexp.Compile(`([^\\]|^)!\[(.*?)\]\((.*?)\)`)

	fs := r.FindAllSubmatch(articleBody, -1)

	for i, _ := range fs {
		url := fmt.Sprintf("https://kicoe-blog.oss-cn-shanghai.aliyuncs.com/%s", string(fs[i][3]))
		syncImage(url)
	}

	rep := []byte("${1}![${2}](https://raw.githubusercontent.com/moonprism/blogs/master/image/${3})")
	articleBody = r.ReplaceAll(articleBody, rep)

	file, err := os.OpenFile(fileName, os.O_TRUNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	handleError(err)
	defer file.Close()

	_, err = file.Write([]byte(resp.Data.(string)))
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

func syncImage(url string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	handleError(err)

	req.Header.Set("Referer", "https://kicoe.com/")

	resp, err := client.Do(req)
	handleError(err)
	defer resp.Body.Close()

	path := "image/"
	s := strings.Split(url, "/")
	fileName := s[len(s)-1]

	out, err := os.OpenFile(path+fileName, os.O_TRUNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	handleError(err)
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	handleError(err)
}
