package data

import (
	"log"
	"strconv"
	"time"
)

const (
	StepSub	int = 0
	StepInline int = 1
	StepCall int = 2
	StepCharge int = 3
	StepFinish int = 4
)

// 2 types of report:archived and ongoing
// DB should only store archived reports,
// ongoing reports should be managered by scheduler
// when ongoing report reach a state(e.g. end charge)
// that is surely
// not going to change in the future, archive it to DB

type Report struct {
	Id                    int // only visible to database
	Num                   int64	// NewReport()
	Charge_id             int
	Charge_mode           int
	Username              string	// NewReport()
	User_id               int// NewReport()
	Request_charge_amount float64
	Real_charge_amount    float64
	Charge_time           int64
	Charge_fee            float64
	Service_fee           float64
	Tot_fee               float64
	Step                  int// NewReport()
	Queue_number          string
	Subtime               int64// NewReport()
	Inlinetime            int64
	Calltime              int64
	Charge_start_time     int64
	Charge_end_time       int64
	Terminate_flag        bool
	Terminate_time        int64
	Failed_flag           bool
	Failed_msg            string
}

func (r *Report) Archive() {
	statement :=
		"insert into reports (num,charge_id,charge_mode,username,user_id,request_charge_amount,real_charge_amount,charge_time,charge_fee,service_fee,tot_fee,step,queue_number,subtime,inlinetime,calltime,charge_start_time,charge_end_time,terminate_flag,terminate_time,failed_flag,failed_msg)" +
			"values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)" +
			"returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(r.Num, r.Charge_id, r.Charge_mode, r.Username, r.User_id, r.Request_charge_amount, r.Real_charge_amount, r.Charge_time, r.Charge_fee, r.Service_fee, r.Tot_fee, r.Step, r.Queue_number, r.Subtime, r.Inlinetime, r.Calltime, r.Charge_start_time, r.Charge_end_time, r.Terminate_flag, r.Terminate_time, r.Failed_flag, r.Failed_msg).Scan(&r.Id)
	if err != nil {
		log.Fatal(err)
	}
}

func NewReport(u User) Report{
	rp := Report{}
	rp.Num = generateReportNum(u.Id)
	rp.Username = u.Name
	rp.User_id = u.Id
	rp.Step = StepSub
	rp.Subtime = time.Now().Unix()
	return rp
}

// using unix user_id + timestamp
func generateReportNum(user_id int) int64 {
	ts := time.Now().Unix()
	str := strconv.FormatInt(ts, 10)
	str = strconv.Itoa(user_id) + str
	ret, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return ret
}
