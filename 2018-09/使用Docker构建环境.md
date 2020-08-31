>[info] 最近因为无法更新，又将笔记本上的deepin系统重装了一遍。

# Docker

Docker是基于Linux容器虚拟化(LXC)技术实现的封装，通过提供的接口可以很方便的进行虚拟环境容器的管理，运行等。


## 镜像

> Docker 镜像是一个特殊的文件系统，除了提供容器运行时所需的程序、库、资源、配置等文件外，还包含了一些为运行时准备的一些配置参数（如匿名卷、环境变量、用户等）。镜像不包含任何动态数据，其内容在构建之后也不会被改变。

Docker镜像命令
```sh
# 列出所有镜像
docker image ls
# 删除镜像
docker image rm [image_name]
```

## 容器

> 镜像（Image）和容器（Container）的关系，就像是面向对象程序设计中的 类 和 实例 一样，镜像是静态的定义，容器是镜像运行时的实体。容器可以被创建、启动、停止、删除、暂停等。

Docker容器命令
```sh
# 列出正在运行的容器
docker container ls
# 列出所有容器
docker container ls --all
# 删除未运行的容器
docker container rm [container_id]
```

## 一个简单的例子

因为我本地的数据库管理工具一般用adminer，所以开始的例子就先搭建一个集成php与mysql的Docker环境吧。

首先使用Docker官方的MySQL镜像运行容器mysql:8.0.12（当本地不存在时将自动从官方服务器获取）
```sh
docker container run \
  -d \
  --rm \
  --name db \
  --env MYSQL_ROOT_PASSWORD=123456 \
  mysql:8.0.12
```
* `-d` 容器启动后在后台运行
* `--rm` 容器停止后自动删除容器文件
* `--name` 容器命名
* `--env` 传入环境变量，将被用来当做mysql的root密码 

接下来运行`docker container ls`就可以看到正在运行中的容器了

基于Docker官方的PHP镜像创建一个运行PHPweb服务的容器

先从[adminer](https://www.adminer.org/)的官方网站上下载主要的php文件和喜欢的css文件，放在项目根目录中命名为index.php
```
wget https://github.com/vrana/adminer/releases/download/v4.6.3/adminer-4.6.3-mysql.php -O index.php
wget https://raw.githubusercontent.com/vrana/adminer/master/designs/pepa-linha/adminer.css
```
接下来在同目录下创建`Dockerfile`文件
```
FROM php:7.2-cli
RUN docker-php-ext-install pdo_mysql
COPY . /usr/src/adminer
WORKDIR /usr/src/adminer
CMD php -S 0.0.0.0:80
```
意思是基于官方php:7.2-cli镜像，构建时运行docker-php-ext-install命令安装mysql依赖，将本地目录下所有文件copy到容器/usr/src/adminer中且设置其为工作目录，镜像容器启动时运行`php -S 0.0.0.0:80`

接下来运行
```sh
docker image build -t adminer .
```
就会根据Dockerfile文件生成一个命名为`adminer`的镜像。运行这个镜像
```sh
docker run \
-d \
--rm \
-p 8080:80 \
--link db:mysql8 \
adminer
```
* `-p 8080:80` 将本地真实环境的8080端口绑定到容器中的80端口
* `--link db:mysql` 本容器将链接db容器运行，且db容器在本容器中的别名为mysql8

通过浏览器访问`127.0.0.1:8080`，数据库填链接的别名`mysql8`输入设定的root密码就可以进入adminer的管理界面了

---
### References
1. [Docker 微服务教程](http://www.ruanyifeng.com/blog/2018/02/docker-wordpress-tutorial.html)
2. [Docker - 从入门到实践](https://yeasy.gitbooks.io/docker_practice/basic_concept/)
