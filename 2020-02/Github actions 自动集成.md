## github action

actions 是各种持续集成的操作片段，比如lint、test、ssh部署等，用户可以编写自己的脚本来实现并分享这些脚本让其他用户使用。

> 可以在[官方市场](https://github.com/marketplace?type=actions)中查找到他人分享的action

## workflow

在`.github/workflow`目录下的每个yml文件都将识别为一个工作流，github将会自动运行这些文件。

```yml
name: workflow name
# 指定触发条件
on:
  push:
    branches:
      - master # 推送master分支执行工作流
# 设置任务
jobs:
  build:
    name: Build
    runs-on: ubuntu-latest # 运行容器
    steps:

    - name: Set up Go 1.13
      uses: actions/setup-go@v1 # 安装golang环境
      with:
        go-version: 1.13
      id: go

    - name: Check out code into the Go module directory
      uses: actions/checkout@v2 # 检出代码

    - name: Get dependencies
      run: |
        go get -v -t -d ./...
        if [ -f Gopkg.toml ]; then
            curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
            dep ensure
        fi

    - name: Build
      run: go build -v .
```

以上就是官方对于golang代码库默认的workflow，其中的触发条件并不需要是git本身的钩子，甚至可以这样写：

```
on:
  schedule:
  # 星期一到星期五每天两点运行
  - cron: 0 2 * * 1-5
```

编译出可执行文件后，可以自动打tag发release，便于版本回滚。如果是用docker部署的话也可以直接发布容器并通知服务器更新。

## 自动编译并发布docker hub

```yml
jobs:
  publish:
    name: Publish
    runs-on: ubuntu-latest # 运行容器
    steps:
      - uses: actions/checkout@master # 检出代码
      - name: Build
        run: make build # 运行Makefile里的编译任务出可执行文件
      - name: Docker Publish
        uses: elgohr/Publish-Docker-Github-Action@v5
        with:
          name: kicoe/test
          username: kicoe
          password: ${{ secrets.DOCKER_HUB_PWD }}
          dockerfile: docker/Dockerfile
```

可以看到我们引入了一个非官方的action [Publish-Docker-Github-Action](https://github.com/elgohr/Publish-Docker-Github-Action)来发布到docker hub，具体用法可以参阅其文档。

### 敏感信息

在上面的action配置中，password不需要直接写入，只要在github个人项目的settings中配置secrets，就可以通过 `${{ secets.xxx }}` 访问私有信息了。

![](https://kicoe-blog.oss-cn-shanghai.aliyuncs.com/YRXXEjfobGDeLkUBmLRo.jpg)

## 远程ssh部署

首先生成rsa的加密对

```shell
ssh-keygen -t rsa
```

将生成的公钥添加到服务器的`~/.ssh/authorized_keys`文件中，私钥则存入github项目中的secrets。

```yml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Deploy
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.ECS_HOST }}
          username: ${{ secrets.ECS_USER }}
          key: ${{ secrets.ECS_ACCESS_TOKEN }}
          script: |
            pwd
            ls -al
```

当执行到此action时，可以在github中看到如下日志

![](https://kicoe-blog.oss-cn-shanghai.aliyuncs.com/ppuYJZODkYoRPMcXjuGO.jpg)

## blog read-static

下面是应用于自己博客的，前端代码自动部署workflow。在master分支下`web/public/static`中内容有变化的话执行该action。

大部分是使用他人造好的轮子，首先按照`gulpfile.js`使用gulp压缩静态文件，登录ECS将压缩好的文件上传。全是非脚本的静态文件就不考虑用ln了。

```yml
name: read-static

on:
  push:
    branches:
      - master
    paths:
      - 'read/public/static/**'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@master

      - name: gulp
        uses: elstudio/actions-js-build/build@v2
        with:
          wdPath: './read'

      - name: deploy
        uses: appleboy/scp-action@master
        with:
          host: ${{ secrets.ECS_HOST }}
          username: ${{ secrets.ECS_USER }}
          key: ${{ secrets.ECS_ACCESS_TOKEN }}
          source: "read/public/dist/*"
          target: "/blog/dist/"

      - name: cp
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.ECS_HOST }}
          username: ${{ secrets.ECS_USER }}
          key: ${{ secrets.ECS_ACCESS_TOKEN }}
          script: |
            cp -r /blog/dist/read/public/dist /blog/read/public/
```

[read-static.yml](https://github.com/moonprism/blog/blob/master/.github/workflows/read-static.yml)

[actions log](https://github.com/moonprism/blog/actions)