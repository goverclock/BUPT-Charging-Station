# BUPT-Charging-Station

TODO:

- [x] JWT or user_id
- [ ] JWT token in header or response.data.token
- [x] turn all user_id into ctx.Get("user_name")
- [ ] some operation needs admin power, check in handler function
- [ ]前端根据user_id or user_type确定是否为管理员
- [ ] /register/user现在只返回user_id
- [ ] /charge/getChargingMsg是否保留

- [ ] remove all log.Fatal

- [ ] 浮点数显示问题(后端处理)

- [ ] 余额不足Request_charge_amount时直接拒绝充电请求

- [ ] StationReport缺少tot_charge_fee...(/chargeports/getreport(s))

- [ ] endCharge之后没有结束时间

FRONTEND:

- [ ] 格林威治时间转换为北京时间

- [ ] 余额充值时负数卡死,以及刷新界面后不立刻显示余额

- [ ] chargeports/getreport

- [ ] 充电桩状态切换应该发到switchBroken而不是switch路径

- [ ] chargeports/waitingCars
