以前用 chromedp 写过一个漫画爬虫，用于爬取图片生成一些 vol.moe 上没有的 kindle mobi 漫画。本来是设计做成一个简单命令行的，但职业习惯写成 bs 架构，解析那块又没做好解耦，导致越写越复杂，最后在没有前端的情况下可以用 swag 调 api 生成几本想看的草草弃坑。

总之 chromedp 给我留下一股印象: 强大好用，但就是常常会在意想不到的地方出错（其实就是平时浏览网页也会遇到的各种小问题），处理这些错误需要小心谨慎。最后内心断言这工具就只是个玩具罢了。

>[info] ... 只能说，以前工作还是太安逸了

## chromedp

[chromedp](https://github.com/chromedp/chromedp) 是一个使用 go 语言编写的 chrome 包。可以通过[Chrome DevTools Protocol](https://chromedevtools.github.io/devtools-protocol/)操作chrome，获取其页面和cookie等。

## 试用

官网的栗子，其中代码稍微改了下，传了个空配置运行时打开浏览器方便调试。

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	opts := []chromedp.ExecAllocatorOption{}
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// navigate to a page, wait for an element, click
	var example string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://golang.org/pkg/time/`),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.WaitVisible(`body > footer`),
		// find and click "Expand All" link
		chromedp.Click(`#pkg-examples > div`, chromedp.NodeVisible),
		// retrieve the value of the textarea
		chromedp.Value(`#example_After .play .input textarea`, &example),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Go's time.After example:\n%s", example)
}
```

本地安装了 chrome 浏览器的话，运行以上代码，就可以看到浏览器打开页面后关闭，同时将获取的元素信息输出到终端了。

## pixiv

> [pixiv](https://www.pixiv.net/) 是喜欢插画、对插画感兴趣的朋友们的交流平台。在这里，您可以投稿自己的插画作品，并与其他用户轻松展开交流。

其实以前在博客写过如何用 PHP 模拟 pixiv 登录、爬取收藏图片，因为几个版本后代码不可用（fiddler 抓包的时候发现登录又多了一个 google 验证器字段）与整站图片丢失不想找回把而文章删了。

## 选择器

chromedp 中的选择器是标准的 selector，获取非常方便，比如下面的 pixiv 登录页面:

![](https://kicoe-blog.oss-cn-shanghai.aliyuncs.com/OWhruZbOsIIViCBEfELy.jpg)

在浏览器控制台中 copy selector 就可以复制该元素的选择器写法：

```
#LoginComponent > form
```

## 获取 PHPSESSID

```go
func main() {
	...
	// 改写Run代码
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://accounts.pixiv.net/login`),
		chromedp.WaitVisible(`#LoginComponent > form`),
		// 输入框
		chromedp.SendKeys(`#LoginComponent > form > div.input-field-group > div:nth-child(1) > input[type=text]`, username),
		chromedp.SendKeys(`#LoginComponent > form > div.input-field-group > div:nth-child(2) > input[type=password]`, password),
		// 点击登录按钮
		chromedp.Click(`#LoginComponent > form > button`, chromedp.NodeVisible),
		// 等待跳转(登录成功)
		chromedp.WaitNotPresent(`#LoginComponent`),
		// 获取指定cookie
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetAllCookies().Do(ctx)
			if err != nil {
				return err
			}

			for _, cookie := range cookies {
				if (cookie.Name == "PHPSESSID") {
					log.Println(cookie.Value)
				}
			}

			return nil
		}),
	)
}
```

相比于普通爬虫，通过浏览器的写法还是要简单许多的（相应的需要耗费更多计算资源。

## headless & docker

```yml
# docker-compose.yml
services:
  chrome:
    image: chromedp/headless-shell:latest
```