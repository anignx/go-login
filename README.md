# go-login
登录注册

# 注册流程
1.先调用getToken接口，获取手机验证码
2.再调用register接口进行注册
3.注册成功后

# 登录
调用login接口进行登录


# 1.服务注册
  1.在main的init阶段，对conf文件进行解析，并将自己的ip：port + api注册到etcd中
  2.对router中所有的接口都进行注册

# 2.服务发现
  1.获取conf中的所有client，读取服务发现数据，并建立client长链接，写入map中，通过clientInit在业务中获取

# 注意：
  1.服务注册取消时，可以通过超时时间失效（30s），也可以在k8s上线平台处理服务销毁