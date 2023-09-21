
### 一、topic 命名规则

```bash
地区_企业_项目名_pro(dev)_业务名
# 举例：
hb_xiaomi_tg_dev_subscribe
hb_xiaomi_tg_pro_subscribe
```
### 二、broker 地址
内网地址：
```
101.227.212.17:9092
101.227.212.16:9092
101.227.212.15:9092
```

公网地址：
```
10.49.32.158:9092
10.49.32.157:9092
10.49.32.156:9092
```

### 三、测试信息

生产者用户与权限：

账号：writer

密码：writer-123

topic：test_topic

权限：写

---

消费者用户与权限：

账号：reader

密码：reader-123

topic：test_topic

权限：读

---
管理员用户与权限

账号：admin

密码：admin-123

topic：*

权限：读、写

---


