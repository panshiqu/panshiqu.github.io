---
layout: default
---

# expect spawn scp * shell路径名展开
_2018-10-11 10:20:45_

* * *

一直通过 scp 同步多台服务器上配置文件，虽然已经写了 shell 脚本，但是密码还需要手动输入，终于忍不了啦，经查可以使用 expect 改变这种现状（当然不只这一种解决方案）

```
#!/usr/bin/expect
#filename scp.exp
spawn scp -r files/* panshiqu@172.16.10.182:/home/panshiqu
expect "password:"
send "123456\n"
expect eof
```

```
panshiqudeiMac:exp panshiqu$ ./scp.exp 
spawn scp -r files/* panshiqu@172.16.10.182:/home/panshiqu
panshiqu@172.16.10.182's password: 
files/*: No such file or directory
```

这个错误困扰我很长时间，根本原因是 路径名展开 是 shell 的特性，expect 没有

虽然网上有提供 一次仅拷贝一个文件 的替代方案

```
#!/usr/bin/expect
#filename scp2.exp
spawn scp -r [lindex $argv 0] panshiqu@172.16.10.182:/home/panshiqu
expect "password:"
send "123456\n"
expect eof
```

```
#!/bin/bash
#filename scp.sh
for f in `echo files/*`
do
    ./scp2.exp $f
done
```

即便这只是一个辅助工具，而且我也不关心效率，但我也不想如此，一个文件一次调用，建立连接等等，会慢的

反复尝试，终不得解，最后灵机一动，分享给大家

```
#!/bin/bash
#filename scp2.sh
scp -r $1 $2
```

```
#!/usr/bin/expect
#filename scp3.exp
spawn ./scp2.sh files/* panshiqu@172.16.10.182:/home/panshiqu
expect "password:"
send "123456\n"
expect eof
```
