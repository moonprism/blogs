a simple markdown parser for JavaScript, ⚡ [example](https://moonprism.github.io/markdown.js/)

* [支持语法](#h-syntax)
	* [inline](#h-s-inline)
	* [heading](#h-s-heading)
	* [list](#h-s-list)
	* [blockquotes](#h-s-blockquotes)
	* [code](#h-s-code)
	* [table](#h-s-table)
	* [html](#h-s-html)
* [安装与使用](#h-install)
	1. [npm](#h-i-n)
	2. [markdown.min.js](#h-i-m)
* [配置](#h-config)

该项目不是一个纯粹的解析器，没有AST。
只是用正则把markdown提取出token再生成相应的HTML。实现有参考 [marked](https://github.com/markedjs/marked)，但代码和正则表达式都是自己写的，总共三百几十行，所以生成的文件非常小：<a href="https://moonprism.github.io/markdown.js/markdown.min.js"><image style="border-radius: 0;box-shadow:none;padding-top:0;position: relative;top:4px" src="https://img.badgesize.io/moonprism/markdown.js/master/dist/markdown.min.js?compression=gzip&style=flat-square&color=blue"></a>

> 刚提交时只有一百多行

>[success] 时隔三年的重写，把该有的功能都加上了

## 语法 {#h-syntax}

### Inline {#h-s-inline}

* `` `inline code` ``
* \**italicize*\*
* \*\***bold**\*\*
* \[link text](address)
* \!\[image alt text](src)
* \*\*\* or ---
* \~\~~~strikethrough~~\~\~
* \ escape inline

### Heading {#h-s-heading}
```md
## heading {#heading-id}
# h1
## h2
### h3
```
---
# h1
## h2
### h3
---

### List {#h-s-list}

```md
1. Step 1
    * Item a
    * Item b
        1. b1
    * Item c
2. Step 2
    1. first

- [x] task1
- [ ] task2
```

---

1. Step 1
    * Item a
    * Item b
        1. b1
    * Item c
2. Step 2
    1. first

- [x] task1
- [ ] task2

---

### Blockquote {#h-s-blockquotes}

```md
> blockquote
next line

>[info] info

>[danger] danger
> >[success] success
>>>[warning] warning
```

---
> blockquote
next line

>[info] info

>[danger] danger
> >[success] success
>>>[warning] warning

---

### Code {#h-s-code}
---
\```go
```go
package main

func main() {
    println('mdzz')
}
```
\```
---
### Table {#h-s-table}

```md
| left align | right align | center |
| :------| ------: | :------: |
| AND | 0 | 1 |
| 0 | 0 | 0 |
| 1 | 0 | 1 |
```
---
| left align | right align | center |
| :------| ------: | :------: |
| AND | 0 | 1 |
| 0 | 0 | 0 |
| 1 | 0 | 1 |
---

### Html {#h-s-html}

```html
<svg width="99" height="99">
  <circle cx="50" cy="50" r="40" stroke="black" stroke-width="2" fill="#d89cf6"/>
</svg>
```
---

<svg width="99" height="99"><circle cx="50" cy="50" r="40" stroke="black" stroke-width="2" fill="#d89cf6"/></svg>

---

```html
<iframe 
    frameborder="no"
    border="0"
    marginwidth="0"
    marginheight="0"
    width=330
    height=86
    src="//music.163.com/outchain/player?type=2&id=30482372&auto=0&height=66">
</iframe>
```
---

<iframe 
		frameborder="no"
		border="0"
		marginwidth="0"
		marginheight="0"
		width=330
		height=86
		src="//music.163.com/outchain/player?type=2&id=30482372&auto=0&height=66">
</iframe>

---

## 安装与使用 {#h-install}

### NPM {#h-i-n}

```js
$npm install moonprism-markdown --save

import markdown from 'moonprism-markdown'
let html = markdown('# hello world')
```

### Min.js {#h-i-m}

download [markdown.min.js](https://moonprism.github.io/markdown.js/markdown.min.js)，and import file in your page.

```html
<script type="text/javascript" src="./markdown.min.js"></script>
<script type="text/javascript">
    var html = markdown('# hello world')
</script>
```

## 配置 {#h-config}
```js
markdown('text', {
    imageCDN: 'https://cdn.xx',
    linkTargetBlank: true,
    lineParse: function(str) {return str},
    codeParse: function(str) {return str},
	debug: false
})
```