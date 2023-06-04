# BUPT-Charging-Station

TODO:

- [x] 充电桩数量作为系统启动参数

- [x] JWT or user_id

// - [ ] JWT token in header or response.data.token

- [x] turn all user_id into ctx.Get("user_name")

// - [ ] some operation needs admin power, check in handler function

- [x] remove all log.Fatal

- [x] 余额不足Request_charge_amount时直接拒绝充电请求

- [x] StationReport缺少tot_charge_fee...(/chargeports/getreport(s))

- [x] system/xxxsettings接口

FRONTEND:

- [x] 根据user_id or user_type确定是否为管理员

- [x] 格林威治时间转换为北京时间

- [x] 余额充值时负数卡死,以及刷新界面后不立刻显示余额

- [x] chargeports/getreport

- [x] 充电桩状态切换应该发到switchBroken而不是switch路径

- [x] endCharge之后没有结束时间

- [x] 浮点数显示问题

---

// - [ ] 后端也处理一下浮点数问题

- [x] 加入车辆总容量,前后端逻辑判断
