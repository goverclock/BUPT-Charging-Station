package scheduler

import (
	"buptcs/data"
	"log"
	"time"
)

// assume sched.mu is locked
// return nil if no report found
func ongoingReportByUser(u data.User) *data.Report {
	if sched.mu.TryLock() {
		log.Fatal("should have locked sched.mu in ongoingReportByUser")
	}
	for _, r := range sched.ongoing_reports {
		if r.User_id == u.Id {
			return r
		}
	}
	return nil
}

// returns all reports, no matter archived or ongoing
func ReportsByUser(u data.User) []data.Report {
	sched.mu.Lock()
	defer sched.mu.Unlock()

	rps := []data.Report{}
	// get archived reports from DB
	rps = append(rps, data.ArchivedReportsByUser(u)...)
	// get ongoing report for the user
	for _, r := range sched.ongoing_reports {
		if r.User_id == u.Id {
			rps = append(rps, *r)
		}
	}
	return rps
}

// assume sched.mu is locked
// archive and remove from sched's ongoing_reports
// also remove the car in the station
// also updates user's balance
func archiveOngoingReport(rp *data.Report) {
	if sched.mu.TryLock() {
		log.Fatal("should have locked sched.mu in archiveOngoingReport")
	}
	for ri, r := range sched.ongoing_reports {
		if r.Num == rp.Num {
			sched.ongoing_reports = append(sched.ongoing_reports[:ri], sched.ongoing_reports[ri+1:]...)
			r.Archive()
			// remove car from station
			st := stationById(r.Charge_id)
			st.Leave(r.Queue_number)
			// update user's balance
			user, err := data.UserByName(r.Username)
			if err != nil {
				log.Fatal("no such user ", user)
			}
			user.Balance -= r.Tot_fee
			user.Update()
			break
		}
	}
}

// TODO: check if user has no ongoing report before creating new
// assume sched.mu is locked
func newOngoingReport(u data.User) *data.Report {
	if sched.mu.TryLock() {
		log.Fatal("should have locked sched.mu in newOngoingReport")
	}
	rp := data.NewReport(u)
	sched.ongoing_reports = append(sched.ongoing_reports, &rp)
	return &rp
}

// assume sched.mu is locked
// only update real_charge_amount, charge_time, charge_fee, service_fee, tot_fee
func updateOngoingReports() {
	cur := time.Now().Unix()
	for _, r := range sched.ongoing_reports {
		// if user isn't charging, nothing should update here
		if r.Step != data.StepCharge {
			continue
		}
		
		// actually charge the car here
		st := stationById(r.Charge_id)
		elec_fee, service_fee := getFee()
		select {
		case elec := <-st.ChargeChan:
			finished := false
			if r.Real_charge_amount + elec >= r.Request_charge_amount {
				elec = r.Request_charge_amount - r.Real_charge_amount
				finished = true
			}
			r.Real_charge_amount += elec // update real_charge_amount
			r.Charge_fee += elec * elec_fee	// update charge_fee
			r.Service_fee += elec * service_fee	// update service_fee
			r.Tot_fee = r.Charge_fee + r.Service_fee	// update tot_fee
			if finished {
				r.Charge_end_time = cur
				r.Step = data.StepFinish
				archiveOngoingReport(r)
			}
		default:
		}

		r.Charge_time = (r.Charge_start_time - cur) / 60 // update charge_time
	}
}
