---
layout: default
---

# 依据官方文档源码编译Nginx实现普通TCP服务反向代理负载均衡支持SSL/TLS的配置
_2018-08-23 11:45:28_

* * *

一直没有怎么用过Nginx，用的最多的就是搭建静态网站，确实屈才啊。最近要用Nginx部署反向代理服务，网上讲解的文章太多，不能拿来直接用就觉得写的不好，随着研究的深入现整理本文，希望读到的朋友会觉得有用。

看了几篇文章无果后，便去看了Nginx的官方文档，搜了关键字proxy，发现仅有四篇文档
[WebSocket proxying](https://nginx.org/en/docs/http/websocket.html)
[ngx_http_proxy_module](https://nginx.org/en/docs/http/ngx_http_proxy_module.html)
[ngx_mail_proxy_module](https://nginx.org/en/docs/mail/ngx_mail_proxy_module.html)
[ngx_stream_proxy_module](https://nginx.org/en/docs/stream/ngx_stream_proxy_module.html)
从名字可以排除掉前三个（我要代理的是普通的TCP服务），现在看第四篇文档
```
The ngx_stream_proxy_module module (1.9.0) allows proxying data streams over TCP, UDP (1.9.13), and UNIX-domain sockets.
```
大致翻译如下：这个模块允许代理数据流通过TCP、UDP、UNIX域套接字

官方文档提供的示例配置
```
server {
    listen 127.0.0.1:12345;
    proxy_pass 127.0.0.1:8080;
}
```
感觉正合我意，安装尝试一番吧

Ubuntu 16.04.4 LTS 安装 nginx/1.10.3 (Ubuntu)
```
apt-get install nginx
```

nginx.conf新增如下配置
```
stream {
    server {
        listen 6666;
        proxy_pass 172.16.10.177:6666;
    }
}
```
通过以下命令测试配置正确，重启使配置生效
```
nginx -t
service nginx restart
```
现在你已经可以通过部署Nginx的服务器IP访问到177服务器上的真正服务啦

现实服务很可能不只一组，想要借此实现负载均衡，就要好好理解proxy_pass这段说明啦
```
If a domain name resolves to several addresses, all of them will be used in a round-robin fashion. In addition, an address can be specified as a server group.

The address can also be specified using variables (1.11.3):

proxy_pass $upstream;

In this case, the server name is searched among the described server groups, and, if not found, is determined using a resolver.
```
接下来我们要好好看看所谓的[server groups](https://nginx.org/en/docs/stream/ngx_stream_upstream_module.html)

研究一番后我们可以写出如下配置
```
stream {
    upstream backend {
        server 172.16.10.79:6666;
        server 172.16.10.177:6666;
    }

    server {
        listen 6666;
        proxy_pass backend;
    }
}
```
生效配置，循环访问79和177服务

有时线上的环境只能通过源码编译，这里给出简单过程
```
Ubuntu
apt-get install gcc
apt-get install libpcre3-dev
apt-get install zlib1g-dev
apt-get install libssl-dev

CentOS
yum install gcc
yum install pcre-devel
yum install zlib-devel
yum install openssl-devel
```

下载解压后执行以下命令
```
./configure
make
make install
```

你会发现```./nginx -t```报以下错误
```
nginx: [emerg] unknown directive "stream"
```

我们就在官方[Building nginx from Sources](https://nginx.org/en/docs/configure.html)文档中寻求帮助
```
--with-stream
--with-stream=dynamic
enables building the stream module for generic TCP/UDP proxying and load balancing. This module is not built by default.
```

所以再来一遍喽
```
./configure --with-stream
make
```

但是并不需要重新安装，你只需要覆盖可执行程序就好
```
mv /usr/local/nginx/sbin/nginx /usr/local/nginx/sbin/nginx.old
mv objs/nginx /usr/local/nginx/sbin/nginx
```

啰哩啰唆了这么多，大家可能发现我并不是想说某一种特定需求该如何配置，而是在分享我是如何借助官方文档写出配置的历程

再设想一种场景，目前有支持WebSocket的服务，需要扩展其支持SSL/TLS协议，我们就可以摸索出如下配置
```
./configure --with-stream --with-stream_ssl_module
```
```
stream {
    server {
        listen              7667 ssl;
        proxy_pass          172.16.10.177:7667;
        ssl_protocols       TLSv1 TLSv1.1 TLSv1.2;
        ssl_ciphers         AES128-SHA:AES256-SHA:RC4-SHA:DES-CBC3-SHA:RC4-MD5;
        ssl_certificate     /usr/local/nginx/conf/cert.pem;
        ssl_certificate_key /usr/local/nginx/conf/cert.key;
        ssl_session_cache   shared:SSL:10m;
        ssl_session_timeout 10m;
    }
}
```

从官方文档中可以看出，Nginx甚至可以代理WebSocket，如此以来，自己可以只实现普通的TCP监听，其它的各种类型ws/wss都可以由Nginx代理
譬如有一套游戏服务，客户端现在新增H5版本用WebSocket实现，之后客户端再出微信小游戏版本，微信要求支持SSL/TLS，这些都可以通过Nginx实现

最后说一下我对反向代理和正向代理的理解
正向代理：客户知晓他正在使用代理，譬如翻墙
反向代理：客户并不知道他正在被代理到合适的服务
