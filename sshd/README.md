# socks5 服务

## 1.1 docker 建立

```bash
# 62002 为ssh端口 62001 为socks5 监听端口
docker run -tid --name docker-ubuntu18.04-rocktan001_socks5 \
-e SSH_PORT=62002 \
-p 62002:62002 \
-p 62001:62001 \
-e ROOT_PWD=root \
-e ROCK_USER_PWD=F96AEB124CXIAOQIANG4423 \
--restart always \
rocktan001/docker-ubuntu18.04-rocktan001:v2.0
```

## 1.2 添加socks5 服务
 > 

```bash
apt-get install  lrzsz

# 编译sshd_v5 ,上传到服务器端
vim /etc/supervisord.conf 
[program:sshd_v5]
command=/usr/bin/sshd_v5
autostart=true
autorestart=true
environment=LANG=C.UTF-8
```