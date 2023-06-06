package data

import (
	"buptcs/vtime"
	"fmt"
	"log"
	"strconv"
)

const (
	StepSub    int = 0
	StepInline int = 1
	StepCall   int = 2
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
	Id                    int     // only visible to database
	Num                   int64   `json:"num"` // NewReport()
	Charge_id             int     `json:"charge_id"`
	Charge_mode           int     `json:"charge_mode"`
	Username              string  `json:"username"` // NewReport()
	User_id               int     `json:"user_id"`  // NewReport()
	Request_charge_amount float64 `json:"request_charge_amount"`
	Real_charge_amount    float64 `json:"real_charge_amount"` // ticker
	Charge_time           int64   `json:"charge_time"`        // ticker
	Charge_fee            float64 `json:"charge_fee"`         // ticker
	Service_fee           float64 `json:"service_fee"`        // ticker
	Tot_fee               float64 `json:"tot_fee"`            // ticker
	Step                  int     `json:"step"`               // NewReport()
	Queue_number          string  `json:"queue_number"`
	Subtime               int64   `json:"subtime"`    // NewReport()
	Inlinetime            int64   `json:"inlinetime"` // i.e. StageQueueing
	Calltime              int64   `json:"calltime"`
	Charge_start_time     int64   `json:"charge_start_time"`
	Charge_end_time       int64   `json:"charge_end_time"`
	Terminate_flag        bool    `json:"terminate_flag"`
	Terminate_time        int64   `json:"terminate_time"`
	Failed_flag           bool    `json:"failed_flag"` // chargeSubmit
	Failed_msg            string  `json:"failed_msg"`  // chargeSubmit
}

func (r *Report) Archive() {
	statement :=
		"insert into reports (num,charge_id,charge_mode,username,user_id,request_charge_amount,real_charge_amount,charge_time,charge_fee,service_fee,tot_fee,step,queue_number,subtime,inlinetime,calltime,charge_start_time,charge_end_time,terminate_flag,terminate_time,failed_flag,failed_msg)" +
			"values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)" +
			"returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(r.Num, r.Charge_id, r.Charge_mode, r.Username, r.User_id, r.Request_charge_amount, r.Real_charge_amount, r.Charge_time, r.Charge_fee, r.Service_fee, r.Tot_fee, r.Step, r.Queue_number, r.Subtime, r.Inlinetime, r.Calltime, r.Charge_start_time, r.Charge_end_time, r.Terminate_flag, r.Terminate_time, r.Failed_flag, r.Failed_msg).Scan(&r.Id)
	if err != nil {
		log.Println(err)
	}
}

func NewReport(u User) Report {
	rp := Report{}
	rp.Num = generateReportNum(u.Id)
	rp.Username = u.Name
	rp.User_id = u.Id
	rp.Step = StepSub
	rp.Subtime = vtime.Now().Unix()
	return rp
}

func ArchivedReportsByUser(u User) (rps []Report) {
	rp := Report{}
	rows, err := Db.Query("SELECT * FROM reports WHERE user_id = $1", u.Id)
	if err != nil {
		log.Println(err, "Db.Query()")
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&rp.Id, &rp.Num, &rp.Charge_id, &rp.Charge_mode, &rp.Username, &rp.User_id, &rp.Request_charge_amount, &rp.Real_charge_amount, &rp.Charge_time, &rp.Charge_fee, &rp.Service_fee, &rp.Tot_fee, &rp.Step, &rp.Queue_number, &rp.Subtime, &rp.Inlinetime, &rp.Calltime, &rp.Charge_start_time, &rp.Charge_end_time, &rp.Terminate_flag, &rp.Terminate_time, &rp.Failed_flag, &rp.Failed_msg)
		if err != nil {
			log.Println(err, "rows.Scan()")
		}
		rps = append(rps, rp)
	}
	return
}

// using unix user_id + timestamp
func generateReportNum(user_id int) int64 {
	ts := vtime.Now().Unix()
	str := strconv.FormatInt(ts, 10)
	str = strconv.Itoa(user_id) + str
	ret, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Println(err)
	}
	return ret
}

func (r Report) String() string {
	ret := ""
	ret += fmt.Sprintln("Id                   ", r.Id)
	ret += fmt.Sprintln("Num                  ", r.Num)
	ret += fmt.Sprintln("Charge_id            ", r.Charge_id)
	ret += fmt.Sprintln("Charge_mode          ", r.Charge_mode)
	ret += fmt.Sprintln("Username             ", r.Username)
	ret += fmt.Sprintln("User_id              ", r.User_id)
	ret += fmt.Sprintln("Request_charge_amount", r.Request_charge_amount)
	ret += fmt.Sprintln("Real_charge_amount   ", r.Real_charge_amount)
	ret += fmt.Sprintln("Charge_time          ", r.Charge_time)
	ret += fmt.Sprintln("Charge_fee           ", r.Charge_fee)
	ret += fmt.Sprintln("Service_fee          ", r.Service_fee)
	ret += fmt.Sprintln("Tot_fee              ", r.Tot_fee)
	ret += fmt.Sprintln("Step                 ", r.Step)
	ret += fmt.Sprintln("Queue_number         ", r.Queue_number)
	ret += fmt.Sprintln("Subtime              ", r.Subtime)
	ret += fmt.Sprintln("Inlinetime           ", r.Inlinetime)
	ret += fmt.Sprintln("Calltime             ", r.Calltime)
	ret += fmt.Sprintln("Charge_start_time    ", r.Charge_start_time)
	ret += fmt.Sprintln("Charge_end_time      ", r.Charge_end_time)
	ret += fmt.Sprintln("Terminate_flag       ", r.Terminate_flag)
	ret += fmt.Sprintln("Terminate_time       ", r.Terminate_time)
	ret += fmt.Sprintln("Failed_flag          ", r.Failed_flag)
	ret += fmt.Sprintln("Failed_msg           ", r.Failed_msg)
	return ret
}
