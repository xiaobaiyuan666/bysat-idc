# 魔方业务系统 / 魔方云接口摘录

来源：

- 文档入口：<https://my.idcsmart.com/doc/>
- 文档数据接口：<https://my.idcsmart.com/v1/doc>

我已确认文档树中存在 `server\\mf_cloud` 与 `addon\\idcsmart_cloud` 相关接口分组。

## 已确认的前台 `mf_cloud` 路由

- `GET /console/v1/product/:id/mf_cloud/order_page`
- `GET /console/v1/product/:id/mf_cloud/image`
- `POST /console/v1/product/:id/mf_cloud/duration`
- `GET /console/v1/mf_cloud`
- `GET /console/v1/mf_cloud/:id`
- `GET /console/v1/mf_cloud/:id/status`
- `POST /console/v1/mf_cloud/:id/on`
- `POST /console/v1/mf_cloud/:id/off`
- `POST /console/v1/mf_cloud/:id/reboot`
- `POST /console/v1/mf_cloud/:id/reset_password`
- `POST /console/v1/mf_cloud/:id/reinstall`
- `GET /console/v1/mf_cloud/:id/disk`
- `POST /console/v1/mf_cloud/:id/disk/order`
- `POST /console/v1/mf_cloud/:id/disk/resize/order`
- `GET /console/v1/mf_cloud/:id/snapshot`
- `POST /console/v1/mf_cloud/:id/snapshot`
- `GET /console/v1/mf_cloud/:id/backup`
- `POST /console/v1/mf_cloud/:id/backup`
- `GET /console/v1/mf_cloud/:id/vpc_network`
- `POST /console/v1/mf_cloud/:id/vpc_network`
- `PUT /console/v1/mf_cloud/:id/vpc_network`
- `GET /console/v1/mf_cloud/:id/common_config`
- `POST /console/v1/mf_cloud/:id/common_config/order`
- `GET /console/v1/mf_cloud/:id/whether_renew`

## 已确认的后台 `mf_cloud` 路由

- `GET /admin/v1/mf_cloud/:id`
- `POST /admin/v1/mf_cloud/:id/on`
- `POST /admin/v1/mf_cloud/:id/off`
- `POST /admin/v1/mf_cloud/:id/reboot`
- `POST /admin/v1/mf_cloud/:id/hard_off`
- `POST /admin/v1/mf_cloud/:id/hard_reboot`
- `GET /admin/v1/mf_cloud/:id/status`
- `POST /admin/v1/mf_cloud/:id/reset_password`
- `POST /admin/v1/mf_cloud/:id/reinstall`
- `GET /admin/v1/mf_cloud/:id/remote_info`
- `GET /admin/v1/mf_cloud/:id/ip`
- `POST /admin/v1/mf_cloud/:id/ip`
- `PUT /admin/v1/mf_cloud/:id/ip`
- `DELETE /admin/v1/mf_cloud/:id/ip`
- `GET /admin/v1/mf_cloud/config`
- `PUT /admin/v1/mf_cloud/config`

## 与当前系统的映射

当前代码里 `src/lib/providers/mofang-cloud.ts` 已做了 provider 抽象，默认用 mock：

- `provisionService`
- `activateService`
- `suspendService`
- `renewService`
- `terminateService`
- `syncService`

你后续要做的真实对接建议顺序：

1. 先用文档中的 `status / on / off / reboot / detail` 把实例生命周期打通。
2. 再补 `disk / snapshot / backup / vpc_network` 这些附加资源。
3. 最后再把 `common_config / recommend_config / upgrade_defence / whether_renew` 接入订单和续费规则。
