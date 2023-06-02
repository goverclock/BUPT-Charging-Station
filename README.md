# BUPT-Charging-Station

TODO:

- [x] JWT or user_id

// - [ ] JWT token in header or response.data.token

- [x] turn all user_id into ctx.Get("user_name")

// - [ ] some operation needs admin power, check in handler function

- [ ] 前端根据user_id or user_type确定是否为管理员

- [ ] /charge/getChargingMsg是否保留

- [ ] remove all log.Fatal

- [ ] 浮点数显示问题

- [x] 余额不足Request_charge_amount时直接拒绝充电请求

- [x] StationReport缺少tot_charge_fee...(/chargeports/getreport(s))

FRONTEND:

- [ ] 格林威治时间转换为北京时间

- [ ] 余额充值时负数卡死,以及刷新界面后不立刻显示余额

- [x] chargeports/getreport

- [ ] 充电桩状态切换应该发到switchBroken而不是switch路径

- [ ] chargeports/waitingCars

- [ ] endCharge之后没有结束时间
