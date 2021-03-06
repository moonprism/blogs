我一般把vim当做编辑器而不是IDE来使用，在IDEA系列的IDE（当然vscode也有类似）中都可以搜索到vim的插件ideaVIm，安装这个插件后可以按照vimrc的格式在用户目录编写vim配置文件 `~/.ideavimrc` （windows下是在用户目录下`_ideavimrc`）。

## 基础

vim有两种重要的模式，命令模式和输入模式。

刚启动vim的时候默认进入命令模式，这时候键盘上的按键被识别为各种命令。

按下 `a` `i` `o` 等命令将会进入普通编辑器的输入模式

## 光标移动

| - | - |
| ----------- | ------------------- |
| `gd` | 跳转定义 |
| `Ctrl` + `i` \ `o` | 前进\后退，可以配合gd命令使用 |
| `h` `j` `k` `l` | 左下上右移动光标 |
| `Ctrl` + `u` \ `d` | 向上\下移动半页 |
| `^` \ `$` | 移动到行的开头\末尾 |
| `w` | 移动到下一个单词的开头 |
| `e` | 移动到下一个单词的末尾 |
| `b` | 移动到上一个单词开头 |

## 更新数据

| - | - |
| --------- | ------------------------ |
| `i` | 进入输入模式 insert |
| `a` | 进入输入模式 append |
| `o` | 移动光标到下一行开头并进入输入模式 |
| `x` | 删除当前字符 |
| `X` | 删除上一个字符, 前面跟数字表示执行多少次 |
| `dd` | 剪切整行，前面可以跟数字 |
| `r` | 替换当前字符 |
| `ciw` | 删除当前词并进入插入模式 |
| `yy` | copy整行，前面可以跟数字 |
| `p` | 粘贴 |
| `u` | 撤销操作 |
| `Ctrl` + `r` | 反撤销 |

## 选中

| - | - |
| --------- | ----------------- |
| `v` | 进入选中模式 |
| `Ctrl` + `v` | 进入选中块模式 |
| `V` | 进入选中行模式 |

>[info] 选中模式下可以用光标移动的命令来选择字符统一操作

添加注释还是建议直接用IDE里的快捷键，先使用 `V` 命令选中所有你想要添加注释的行，再按 `Ctrl` + `/` 添加注释。

### 关于 `ciw`

这是一个vim中很常用的命令组合，指 change in word 更改当前词，同样可以使用 `ci'` `ci"` 等命令更新 `''` `""` 中的内容。change也可以换成别的操作命令，比如：`diw`删除当前词，`yiw` copy当前词等。

### 录制宏

```
qa // 选择a寄存器开始录制
...
q // 结束录制
@a // 执行
```

## 简单配置

可以注意到上面列举出来的常用命令很多都要用到 `Ctrl`, 而一般的键盘的 `Ctrl` 都在左下角很难操作到的地方，这其实是因为vim设计的时候键盘布局 `Ctrl` 是在现在 `tap` 键的位置。所以对于这个问题最好的方法是类似HHKB将基本没什么用处的大写锁定 `Caps` 键改成 `Ctrl`（大写输入可以用 `Shift` 键配合输入）。

当然也有人喜欢将 `Caps`改成更常用的`Esc`， `Esc` 在vim中常常用来在输入模式下返回命令模式，这里推荐使用vim的配置在输入模式下快速连按jk触发返回命令模式：
```
inoremap jk <Esc>
```

上面列举的常用命令中还有一组很难按到的就是 `&` `^` 控制光标移动到行的开头和末尾，可以配置

```
noremap <C-h> ^
noremap <C-l> $
```
通过 `Ctrl` + `h` 和 `Ctrl` + `l` 来代替它们。

### vscode

上面的配置都是使用 vimrc 的格式来写的，如果是使用 vscode 的话，安装vim插件后需要到 settings.json 中设置：

```json
  "vim.insertModeKeyBindings": [
    {
      "before": ["j", "k"],
      "after": ["<Esc>"]
    }
  ],
  "vim.normalModeKeyBindingsNonRecursive": [
    {
      "before": ["<C-h>"],
      "after": ["^"]
    },
    {
      "before": ["<C-l>"],
      "after": ["$"]
    }
  ]
```