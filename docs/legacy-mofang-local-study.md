# 本地老版魔方财务学习记录

更新时间：2026-03-21

## 1. 本地运行入口

- 前台根地址：`http://127.0.0.1:8091/`
- 后台入口：`http://127.0.0.1:8091/admin01/`
- 后台账号：`admin`
- 后台密码：`Admin123!`
- 本地数据库：`MariaDB 10.4`
- 本地数据库端口：`3308`
- 本地数据库名：`zjmf_legacy374`

## 2. 当前可用的本地运行方式

已经确认这套老系统不能稳定跑在单个 `php-cgi` 上。

当前本地运行方式：

- `nginx` 监听 `8091`
- `php-cgi` 作为 FastCGI 池，监听：
  - `9074`
  - `9075`
  - `9076`
  - `9077`

相关文件：

- Nginx 配置：
  [legacy-mofang-nginx.conf](C:/Users/Administrator/Desktop/IDC/scripts/legacy-mofang-nginx.conf)
- 启动脚本：
  [start-legacy-mofang-nginx.ps1](C:/Users/Administrator/Desktop/IDC/scripts/start-legacy-mofang-nginx.ps1)
- PHP 7.4 配置：
  [legacy-mofang-php74.ini](C:/Users/Administrator/Desktop/IDC/scripts/legacy-mofang-php74.ini)

## 3. 真正的登录故障根因

### 3.1 不是程序本身坏了

安装已经完成，数据库也已经正常落表：

- `shd_user`
- `shd_auth_rule`
- 后台登录账号已创建

### 3.2 伪静态是必要条件

这点和官方说明一致。ThinkPHP 路由依赖 rewrite。

### 3.3 单个 php-cgi 会导致登录死锁

之前浏览器点登录时：

- `POST /admin01/login` 返回 `504`
- `POST /admin01/async_curl_multi` 也同时返回 `504`

这不是用户名密码错，而是后台登录链路里触发了异步请求，单个 FastCGI worker 被自己卡死。

修正为 4 个 `php-cgi` worker 后，浏览器后台登录已经成功。

## 4. 后台登录链路

### 4.1 关键接口

- `GET /admin01/login_page`
- `GET /admin01/verify`
- `POST /admin01/login`
- `POST /admin01/async_curl_multi`

### 4.2 登录页前端 bundle

- 入口 bundle：
  [Login~f71cff67.eac1556a.js](C:/Users/Administrator/Desktop/IDC/idccw-src/3.7.4/public/admin01/js/Login~f71cff67.eac1556a.js)
- 公共请求封装：
  [app~5a11b65b.061d499b.js](C:/Users/Administrator/Desktop/IDC/idccw-src/3.7.4/public/admin01/js/app~5a11b65b.061d499b.js)
- 登录接口定义：
  [app~d0ae3f07.dc550a27.js](C:/Users/Administrator/Desktop/IDC/idccw-src/3.7.4/public/admin01/js/app~d0ae3f07.dc550a27.js)

### 4.3 当前确认的登录 payload

本地这套环境当前有效最小 payload：

```text
username=admin&password=Admin123!&code=&captcha=
```

附带 query 参数：

```text
request_time=<timestamp>&languagesys=CN
```

### 4.4 当前本地环境的验证状态

- 图形验证码未开启
- 后台二次验证未开启
- 密码前端未发现额外加密

`GET /admin01/login_page` 当前返回：

```json
{
  "status": 200,
  "msg": "请求成功",
  "data": {
    "second_verify_admin": "0",
    "second_verify_action_admin": [""]
  }
}
```

`GET /admin01/verify?name=allow_login_admin_captcha` 当前返回：

```json
{
  "status": 400,
  "msg": "未开启验证码",
  "is_aff": "0"
}
```

## 5. 后台菜单树

已经通过真实登录返回体确认后台菜单结构，不再只是靠页面观察推断。

一级菜单：

- 客户
- 业务
- 财务
- 工单
- 功能
- 设置
- 资源与商店
- 反馈

高价值二级/三级菜单：

- 客户
  - 客户管理
    - 客户列表
    - 实名认证
    - 客户资源池
  - 我的业绩
  - 运营管理
    - 营销推送
    - 推介计划
- 业务
  - 订单
    - 产品订单
    - 续费订单
    - 流量包订单
  - 业务
    - 业务列表
    - 产品暂停请求
- 财务
  - 财务记录
    - 交易流水
    - 账单管理
    - 信用额管理
  - 审核管理
    - 提现审核
  - 发票和合同
    - 发票列表
    - 合同列表
- 工单
  - 工单列表
  - 工单统计
- 功能
  - 插件
  - 统计
    - 年度收入统计
    - 新客户
    - 产品收入
    - 收入排名
  - 系统状态
    - 定时任务状态
    - 任务队列
    - 数据库状态
- 设置
  - 基础设置
    - 客户设置
    - 财务设置
    - 工单设置
  - 系统设置
  - 商品设置
  - 站务设置
- 资源与商店
  - 上下游
    - 上游资源
    - 下游管理
  - 应用商店

## 6. 魔方云模块的真实位置

后台菜单里，“魔方云”并不在“业务”域，而在：

- 商品设置
  - 自动化接口
    - 魔方云

这说明老系统里“魔方云”首先是一个资源/接口提供方，而不是业务主对象。

### 6.1 路由层确认

在后台路由 bundle 里已经确认：

- `path: "zjmfcloud", name: "zjmfcloud"`
- `path: "zjmfcloud-product", name: "zjmfcloudProduct"`

说明至少存在两块：

- 魔方云接口/节点管理
- 魔方云商品/产品映射

相关文件：

- 后台路由 bundle：
  [app~5a11b65b.061d499b.js](C:/Users/Administrator/Desktop/IDC/idccw-src/3.7.4/public/admin01/js/app~5a11b65b.061d499b.js)
- 魔方云页面 chunk：
  [CloudServer~31ecd969.3c96c116.js](C:/Users/Administrator/Desktop/IDC/idccw-src/3.7.4/public/admin01/js/CloudServer~31ecd969.3c96c116.js)
- 魔方云商品 chunk：
  [DCloudProduct~31ecd969.6cafcfe5.js](C:/Users/Administrator/Desktop/IDC/idccw-src/3.7.4/public/admin01/js/DCloudProduct~31ecd969.6cafcfe5.js)

### 6.2 当前从 chunk 中确认的页面职责

`CloudServer~31ecd969` 当前能直接确认是“接口列表/新增接口/连通性状态”类页面，而不是实例控制台。

已确认的界面关键词：

- `rubiks_cube_cloud_list_hint`
- `the_new_interface`
- `search_by_name`
- `api_status`
- `username`
- `password`
- `reinstall_times`

这说明这页更偏：

- 魔方云接口接入
- API 地址和账号配置
- 状态检测
- 默认业务参数配置

而不是“客户云主机实例列表”。

## 7. 当前结论

### 7.1 老系统的对接逻辑分层

从目前确认到的结构看，老系统里魔方云至少拆成三层：

1. 接口接入层  
   配置魔方云节点/API 账号、检测连通性、设置默认能力。

2. 商品映射层  
   将本地商品与魔方云可售卖资源或模板绑定。

3. 业务交付层  
   在订单支付、服务开通、续费、暂停、实例动作时再调用实际的云资源动作。

### 7.2 这也是新系统现在和老系统差距大的根因

当前新系统如果只是“服务详情里直接点 provider 动作”，其实只做到了第 3 层的一小部分。

真正老系统完整的是：

- 接口管理
- 商品映射
- 订单/账单驱动
- 服务生命周期
- 资源同步

## 8. 下一步学习顺序

优先继续抓这几块：

1. 后台首页 `#/home-page`
2. 客户列表与客户详情
3. 产品订单与订单详情
4. 账单管理与账单详情
5. 魔方云接口页 `#/zjmfcloud`
6. 魔方云商品页 `#/zjmfcloud-product`
7. 业务列表 `#/customer-product`

## 9. 当前结论对新系统重建的影响

新系统后面不能只补“实例动作按钮”，必须补齐：

- 魔方云接口管理中心
- 魔方云商品映射中心
- 本地商品与云资源模板绑定
- 订单支付后自动开通链
- 服务与远端云资源双向同步链

