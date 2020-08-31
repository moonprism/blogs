[leetcode - 正则表达式匹配](https://leetcode-cn.com/problems/regular-expression-matching/description/)

## 描述

给定一个字符串 (s) 和一个字符模式 (p)。实现支持 `.` 和 `\*` 的正则表达式匹配。

* `.` 匹配任意单个字符。
* `\*` 匹配零个或多个前面的元素。

匹配应该覆盖整个字符串 (s) ，而不是部分字符串

## 解题思路

这个题目主要是模拟正则表达式中的`\*`匹配。

正则表达式中`\*`表示匹配其前一个字符0~n遍。可以先假设`\*`只匹配前一个字符0遍，当匹配发生不等的时候将`\*`匹配的字符量+1

比如
```go
s = "abc"
p = "a*b"
```
先让模式p中a\*匹配0个a，这时候'a'和'b'无法匹配，回溯p到a\*匹配一个a，这样a,b都可以匹配，但接下来s里的'c'无法匹配了。再回溯p到a\*匹配两个a发现第二个a根本无法匹配，这时候说明再多匹配a\*次数也无济于事，可以判断为不匹配。

如果模式p中存在多个`\*`，可以先让所有`\*`都匹配0次，从后面的开始累加。当发生无法匹配时回到前一个`\*`累加并重新匹配。使用的数据结构很明显就是栈

### 测试代码

```go
package main

import "fmt"

func main() {
	ss := []byte("abcss")
	pp := []byte("abc.*")
	// 保存*节点的栈结构
	type stack_all struct {
		ssi int
		spi int
	}
      // 字符串与模式下标
	var si, pi int
	// 将遇到的*信息入栈  用数组模拟
	var stack []stack_all
	for ;pi<len(pp) || si<len(ss); {
		// 找到*入栈
		if pi+1<len(pp) && pp[pi+1]=='*' {
			stack = append(stack, stack_all{
				ssi: si,
				spi: pi,
			})
			pi = pi+2
			continue
		}
		// 当前字符匹配成功
		if si<len(ss) && pi<len(pp) && (ss[si]==pp[pi] || pp[pi]=='.') {
			pi++
			si++
		} else if len(stack)>0 && stack[len(stack)-1].ssi<len(ss) && stack[len(stack)-1].spi<len(pp)  {
			// 打印栈信息
			fmt.Println(stack);
			// 后面数字匹配失败,但++之后*还可以匹配
			if ss[stack[len(stack)-1].ssi] == pp[stack[len(stack)-1].spi] || pp[stack[len(stack)-1].spi]=='.' {
                // 将下标重置到这个*（栈顶）之前位置重新匹配
				stack[len(stack)-1].ssi++
				si = stack[len(stack)-1].ssi
				pi = stack[len(stack)-1].spi+2
				continue
			} else {
				// *不可以匹配则出栈
				stack = stack[0: len(stack)-1]
				continue
			}
		} else {
			fmt.Println("F")
			return
		}
	}
	fmt.Println("T")
}

```

> 提交之后执行效率很一般，但第一次独立解出hard级题目还是挺有成就感的。