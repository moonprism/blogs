---
*奉上自己写的搭建饥荒服务器的shell脚本*  -  [installDSTserver.sh](https://github.com/moonprism/installDSTserver.sh)
---

## 快速搭建指南
> 测试 ubuntu14.04 与 debian8.0 可用

### 安装依赖 ：
64位：
`sudo apt-get install libstdc++6:i386 libgcc1:i386 libcurl4-gnutls-dev:i386`

32位：
`sudo apt-get install libstdc++6 libgcc1 libcurl4-gnutls-dev`

### 下载并安装steamcmd ：

`mkdir ~/steamcmd`

`cd ~/steamcmd`

`wget https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz`

`tar -xvzf steamcmd_linux.tar.gz`

### 创建资源文件 ：

`mkdir -p ~/.klei/DoNotStarveTogether/MyDediServer/Master`

`mkdir -p ~/.klei/DoNotStarveTogether/MyDediServer/Caves`

### 获取token ：

（进入游戏在Acct Info中可以找到Generate Server Token按钮）

`vim ~/.klei/DoNotStarveTogether/MyDediServer/cluster_token.txt`

将获取的token复制进去

### 创建 cluster.ini 文件 ：

`vim ~/.klei/DoNotStarveTogether/MyDediServer/cluster.ini`
```
[GAMEPLAY]
game_mode = survival
max_players = 6
pvp = false
pause_when_empty = true


[NETWORK]
cluster_description = 【描述】
cluster_name = 【你的名字】
cluster_intention = cooperative
cluster_password =【密码】


[MISC]
console_enabled = true


[SHARD]
shard_enabled = true
bind_ip = 127.0.0.1
master_ip = 127.0.0.1
master_port = 10889
cluster_key = supersecretkey
```
### 创建 server.ini 文件 ：

`vim ~/.klei/DoNotStarveTogether/MyDediServer/Master/server.ini`
```
[NETWORK]
server_port = 11000


[SHARD]
is_master = true


[STEAM]
master_server_port = 27018
authentication_port = 8768
```

` vim ~/.klei/DoNotStarveTogether/MyDediServer/Caves/server.ini `
```
[NETWORK]
server_port = 11001


[SHARD]
is_master = false
name = Caves


[STEAM]
master_server_port = 27019
authentication_port = 8769
```
### 创建 worldgenoverride.lua 文件 ：
`vim ~/.klei/DoNotStarveTogether/MyDediServer/Caves/worldgenoverride.lua`
```
return {
    override_enabled = true,
    preset = "DST_CAVE",
}
```
### 创建脚本启动 ：

`vim ~/run_dedicated_servers.sh`
```sh
#!/bin/bash

steamcmd_dir="$HOME/steamcmd"
install_dir="$HOME/dontstarvetogether_dedicated_server"
cluster_name="MyDediServer"
dontstarve_dir="$HOME/.klei/DoNotStarveTogether"

function fail()
{
        echo Error: "$@" >&2
        exit 1
}

function check_for_file()
{
    if [ ! -e "$1" ]; then
            fail "Missing file: $1"
    fi
}

cd "$steamcmd_dir" || fail "Missing $steamcmd_dir directory!"

check_for_file "steamcmd.sh"
check_for_file "$dontstarve_dir/$cluster_name/cluster.ini"
check_for_file "$dontstarve_dir/$cluster_name/cluster_token.txt"
check_for_file "$dontstarve_dir/$cluster_name/Master/server.ini"
check_for_file "$dontstarve_dir/$cluster_name/Caves/server.ini"

./steamcmd.sh +force_install_dir "$install_dir" +login anonymous +app_update 343050 validate +quit

check_for_file "$install_dir/bin"

cd "$install_dir/bin" || fail 

run_shared=(./dontstarve_dedicated_server_nullrenderer)
run_shared+=(-console)
run_shared+=(-cluster "$cluster_name")
run_shared+=(-monitor_parent_process $$)

"${run_shared[@]}" -shard Caves  | sed 's/^/Caves:  /' &
"${run_shared[@]}" -shard Master | sed 's/^/Master: /'
```
赋予执行权限：`chmod u+x ~/run_dedicated_servers.sh`
### 运行
`~/run_dedicated_servers.sh`

### 参照

* [Dedicated Server Quick Setup Guide - Linux](http://forums.kleientertainment.com/topic/64441-dedicated-server-quick-setup-guide-linux/)

### 关于如何添加mod
` vim ~/dontstarvetogether_dedicated_server/mods/dedicated_server_mods_setup.lua`
```
// Mod ID的获取可到stream创意工坊中mod的url后缀id中得到
ServerModSetup("Mod ID")
ServerModSetup("Mod ID")
. . .
```
也可以直接将下载的mod文件夹复制到dontstarvetogether_dedicated_server/mods内

` vim ~/.klei/DoNotStarveTogether/MyDediServer/Master/modoverrides.lua`

` vim ~/.klei/DoNotStarveTogether/MyDediServer/Caves/modoverrides.lua`
```
 return {
 ["mod文件夹名"] = { enabled = true },
 ["（创意工坊内文件夹命名）workshop-Mod ID"] = { enabled = true },
. . .
}
```
但重启上面脚本的时候会发现编辑的相关文件都被还原，分析执行脚本可以发现其由两部分组成，一部分是steamcmd的验证更新，这个操作会将编辑的mod配置还原。所以可以执行后一部分来初始化mod后再开启服务器：
```
install_dir="$HOME/dontstarvetogether_dedicated_server"
cluster_name="MyDediServer"
cd "$install_dir/bin" || fail 

run_shared=(./dontstarve_dedicated_server_nullrenderer)
run_shared+=(-console)
run_shared+=(-cluster "$cluster_name")
run_shared+=(-monitor_parent_process $$)
run_shared+=(-shard)

"${run_shared[@]}" Caves  | sed 's/^/Caves:  /' &
"${run_shared[@]}" Master | sed 's/^/Master: /'
```
退出后再开启服务器就没问题了

`screen ~/run_dedicated_servers.sh` screen可以使进程停留在后台，按`ctrl+a+z`可以离线，想调出时输入`screen -x`

### 关于双服务器（一个地面一个地穴）
按照前面的步奏在主要服务器安装游戏建立`~/.klei/DoNotStarveTogether/MyDediServer/Master/`的内容，在地穴服务器安装游戏建立`~/.klei/DoNotStarveTogether/MyDediServer/Caves/`的内容

* 在Master(地面)服务器`~/.klei/DoNotStarveTogether/MyDediServer/cluster.ini`配置中的

```
bind_ip = 127.0.0.1    //这里改为0.0.0.0
```
* 在Caves(地穴)服务器`~/.klei/DoNotStarveTogether/MyDediServer/cluster.ini`配置中的

```
master_ip = 127.0.0.1  //这里改为Master服务器的ip
```