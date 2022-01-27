## 2021-12-15 
	sshd_v2 运行在服务器端，监听62001 端口。
	目前遇到问题： 手机端断开转发后，服务端监听端口62001 对应socket 没有及时关闭。测试等待一分钟，变成IN_WAIT1 状态。

## 2021-12-16
	sshd_v2 修改main-timeout 参数，socket 远端断开后 ，server 30内关闭。

	tcpConn.SetDeadline(time.Now().Add(*maintimeout)) 30s 内有连接情况，也被断开了。应该是作者设计失误。在确认下。	

	增加tcpkeepalive 解决了ssh_forward_v2 断开连接后，tcpkeepalive 超时后，断开socket 问题。

	新增： ssh_forward_v2 监听的端口，还在继续监听。正常应该close 掉。



	新增： keepalive 不知道为什么么会使read 判断出错。

## 2022-01-10 去掉timeout ，目前测试满足keepalive 时间到，自动断开需求。
	1:测试sshd_forward_v2 转发，在网络断开情况下，所有监听的端口(62001)和 已建立连接的端口(62002) ，全部断开。	
		结果：X .  测试流程： 手机端打开sshd_forward_v2 .电脑端打开vncserver, 同时打开vncviwer。断开手机WiFi。62002 端口会不断发送数据到62001 。keepalive 没有断开。
		10:43:06.788998 IP 172.17.0.2.62001 > 112.95.160.227.2261: Flags [P.], seq 0:36, ack 1, win 1399, options [nop,nop,TS val 252377434 ecr 10468460], length 36

		keepalive 底层判断有错。

## 2022-01-11  
	ssh -i /media/disk2/socks5/.ssh/id_rsa -N -D 0.0.0.0:1080 root@www.rocktan001.com -p 1196


	结果：
		一：打开google .显示有连接。网页断开后，连接消失。符合要求。后续按这个标准做。
		二：curl --proxy socks5h://192.168.199.173:1080 www.google.com
			执行完后，连接www.google.com socket 消失。

##2022-01-12
	#关于tcp 死连接的讨论
	https://gmd20.github.io/blog/tcp%E7%9A%84%E6%AD%BB%E8%BF%9E%E6%8E%A5%E6%A3%80%E6%B5%8B-dead-peer-%E9%87%8D%E4%BC%A0%E8%AE%A1%E6%97%B6%E5%99%A8-%E5%90%84%E7%A7%8Dtimeout%E8%B6%85%E6%97%B6%E9%80%89%E9%A1%B9-TCP_USER_TIMEOUT-keepalive-heartbeart%E5%BF%83%E8%B7%B3%E6%9C%BA%E5%88%B6%E7%AD%89/
	netstat -natupo 可以查看到timer 信息
	# 遇到新问题，不停往socket 写数据，keepalive 超时后会立刻启动自动重写机制。
	# 通过往 echo 0 > /proc/sys/net/ipv4/tcp_retries2 可以重写机制失效。
	# 注意是IPV4 .l, err = net.Listen("tcp4", *host+":"+*port) 	
