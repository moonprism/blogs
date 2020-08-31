## cerbot

[cerbot](https://certbot.eff.org/) 是一个从 [Let's Encrypt](https://letsencrypt.org/) 获取证书的客户端工具。

### 获取

```
wget https://dl.eff.org/certbot-auto -O certbot
chmod a+x ./certbot
```

### 生成证书

```
./certbot-auto certonly --standalone -d www.kicoe.com
```

* 如果出现`Problem binding to port 80: Could not bind to IPv4 or IPv6.`，把nginx服务先关掉试试

### 配置nginx

```
server {
    listen 80;
    server_name www.kicoe.com;
    return 301 https://$host$request_uri;
}
server {
    listen 443;
    ssl on;
    server_name www.kicoe.com;
    ssl_certificate /etc/letsencrypt/live/www.kicoe.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/www.kicoe.com/privkey.pem;
    location / {
            proxy_pass http://localhost:8080;
    }
}
```

### 更新证书

因为免费的Let's Encrypt证书有效期只有三个月，可以使用如下命令更新证书

```shell
certbot renew
```