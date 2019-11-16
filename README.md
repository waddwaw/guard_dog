# guard_dog
服务状态监测系统

> 这是一个通过配置文件就可以实时监测服务器的运行状态，代码写的比较low，但是基本功能没有问题，也在实际项目中用了一段时间也比较稳定。

#### 使用方法修改配置文件然后运行可执行文件即可

```html
#服务主配置
[app_config]
#服务运行端口号
port = 18080
#钉钉机器人Webhook 地址需要配置自己的钉钉机器人
#请先电脑端安装钉钉 再参考官方服务api文档 https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq
notify_url = https://oapi.dingtalk.com/robot/send?access_token=
#程序间隔扫时间单位秒 同样的监测消息60分钟内不会重复推送
interval_time = 10
#安全设置->自定义关键词 目前支持自定义关键词一种模式
msg_key = 服务器消息通知
#配置需要@相关人的手机号 多个用 , 号隔开
at = 

#配置需要监控的服务
[monitorings]
redis_server = 127.0.0.1:8888
```
