[leetcode 10. 正则表达式匹配](https://leetcode-cn.com/problems/regular-expression-matching/)

## 动态规划解法

* [啥是动态规划？](https://www.kicoe.com/article/id/20)

假设函数 `dp(i, j) = string[:i]与pattern[:j]能否匹配`

---

```
s: aaab
p: a*b
```

对于上面的字符串`aaab`和模式`a*b`，我们可以把 dp 执行的结果列一个表（`^`表示空）

>[info] 用**`o`**表示可匹配，用**`x`**表示不可匹配

p\s|^|a|a|a|b
:-:|:-:|-|-|-|-
**^**|o|x|x|x|x
**a**|x|o|x|x|x
*****|o|o|o|o|x
**b**|x|x|x|x|o

可以看到第一行也就是 pattern 为空串时，除非 string 也为空即 dp(0, 0)位置，否则其他位置都匹配不上。

### 代码形式

```go
for i := 0; i < max_i; i++ {
	dp[i][0] = i == 0
}
```

---

```
s: abc
p: a*.
```

p\s|^|a|b|c
:-:|:-:|-|-|-
**^**|o|x|x|x
**a**|x|o|x|x
*****|o|o|x|x
**.**|x|o|o|x

如果花点时间逐行逐列填些这样的表，可以很快发现规律：

1. 如果是 pattern[j] 是普通字符，普通的对比 string[i] 与 pattern[j] 是否相等且上一次匹配 dp[i-1][j-1] 是否为 true
2. 如果是 pattern[j] 是`*`，分为两种情况
    * 匹配前一个字符0次，判断dp[i][j-2]是否为true
    * 匹配前一个字符1次及以上，判断 string[i] 与 pattern[j-1] 是否相等且 dp[i-1][j] 是否为 true

> 所有判断是否相等的逻辑都应该考虑 pattern 是`.`的话直接为 true

### 代码形式

```go
if pattern[j] == '*' {
	dp[i][j] = dp[i][j-2] ||
	 (pattern[j-1] == '.' || str[i] == pattern[j-1]) && dp[i-1][j]
} else {
	dp[i][j] = ( pattern[j] == '.' || str[i] == pattern[j] )  && dp[i-1][j-1]
}
```

### 加上限制条件后的完整代码

```go
package main

import (
	"fmt"
)

func main() {
	s := "abc"
	p := "a*."

	max_i := len(s) + 1
	max_j := len(p) + 1
	dp := make([][]bool, max_i, max_i)
	for i := 0; i < max_i; i++ {
		dp[i] = make([]bool, max_j, max_j)
		// init first line [true, false, false ...]
		dp[i][0] = i == 0
	}

	for i := 0; i < max_i; i++ {
		for j := 1; j < max_j; j++ {
			if p[j-1] == '*' {
				dp[i][j] = dp[i][j-2] || ((j >= 2 && (p[j-2] == '.' || i >= 1 && s[i-1] == p[j-2])) && i >= 1 && dp[i-1][j])
			} else {
				dp[i][j] = (j >= 1 && (p[j-1] == '.' || i >= 1 && s[i-1] == p[j-1])) && i >= 1 && dp[i-1][j-1]
			}
		}
	}

	for j := 0; j < max_j; j++ {
		for i := 0; i < max_i; i++ {
			fmt.Printf("%v\t", dp[i][j])
        }
        fmt.Printf("\n")
	}
}
```

看了下当前 leetcode 提交最快的版本也是自底向上的动态规划，不过比我这个多了些优化。
比如当 string 为空时，pattern 只有长度为偶数，且每个偶数位上都是`\*`才能匹配成功，也就是直接判断了我们上面所说的 pattern[j] 为`*`和 dp[i][j-2] 来提前返回，
这样写最大好处就是可以不需要再考虑 string 为空的情况，而 pattern 为空很好处理，就可以仅用 len(s) 和 len(p) 长度的二维数组来进行递推。

> 不过我认为 leetcode 最重要的还是思想啦，繁琐的优化细枝末节反而会丧失乐趣不是吗~
>>[danger] 虽然根本没做过几道题

### 参考

* [leetbook - 正则匹配问题](https://hk029.gitbooks.io/leetbook/%E5%8A%A8%E6%80%81%E8%A7%84%E5%88%92/010.%20Regular%20Expression%20Matching/010.%20Regular%20Expression%20Matching.html)