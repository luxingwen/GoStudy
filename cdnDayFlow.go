package main

import (
	"bmob/daemon/m"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"github.com/astaxie/beego/httplib"
	"log"
	"net/url"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var TimeStamp1 int64
var TimeStamp2 int64
var Status int64

const (
	flowAmount       = 32212254720
	cdn_type   int64 = 0 //1为又拍云
	token            = "Bearer e9915106-eefa-4572-b48a-c8b67ca33294"
)

type Result struct {
	Data  interface{} `json:"data"`
	Peaks interface{} `json:"peaks"`
}

type AppCdnFlowInfo struct {
	AppId     int64  `json:"app_id"`
	UpyunName string `json:"upyun_name"`
	Flow      int64  `json:"flow"`
	Status    int64  `json:"status"`
	Day       int64  `json:"day"`
	Type      int64  `json:"type"`
}

type CdnStatistics struct {
	AppId int64 `json:"app_id"`
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	for {
		log.Println("start...................")
		execute()
		execute2()
		time.Sleep(600 * time.Second)
	}
}

func (a *AppCdnFlowInfo) Flows() {
	db, err := m.GetMySQL()
	if err != nil {
		log.Error("can't connect mysql ,err: ", err)
	}
	defer db.Close()
	rows, err := db.Query("select app_id,upyun_name from t_app_cdn where app_id is not NULL and status!=1")
	if err != nil {
		log.Error("query err:%v", err)
	}
	for rows.Next() {
		appCdn := new(AppCdnFlowInfo)
		if err = rows.Scan(&appCdn.AppId, &appCdn.UpyunName); err == nil {

		} else {
			log.Error("err:%v", err)
		}
	}
}

func (a *AppCdnFlowInfo) QueryCdnFlow(db *sql.DB) (err error) {
	startTime := time.Now().Format("2006-01-02")
	endTime := time.Now().Format("2006-01-02 ")
	flows, err := getFlowSize(a.UpyunName, startTime, endTime)
	if err != nil {
		return err
	}
	a.Flow = flows
	err = a.UpdateCdnDayFlow(db)
	if err != nil {
		return
	}
	if err = a.UpdateCdnMonthFlow(db); err != nil {
		return
	}
	return
}
func (a *AppCdnFlowInfo) UpdateCdnDayFlow(db *sql.DB) (err error) {
	nowDay, err := strconv.FormatInt(time.Now().Format("20060102"), 10)
	if err != nil {
		return
	}
	var appId int64
	err = db.QueryRow("select app_id from  t_cdn_day_statistics where app_id=? and day=?", a.AppId, nowDay)
	if err != nil {
		stmt, err := db.Prepare("insert into t_cdn_day_statistics(flow,day,id)values(?,?,?)")
		if err != nil {
			return
		}
		defer stmt.Close()
		err = stmt.Exec(a.Flow, nowDay, a.AppId)
		if err != nil {
			return
		}
	} else {
		stmt, err := db.Prepare("update t_cdn_day_statistics set flow=? where id=?)")
		if err != nil {
			return
		}
		defer stmt.Close()
		err = stmt.Exec(a.AppId)
		if err != nil {
			return
		}
	}
}
func (a *AppCdnFlowInfo) UpdateCdnMonthFlow(db *sql.DB) (err error) {
	monthStart, monthEdn := getMon()
	nowMonth, err := strconv.FormatInt(time.Now().Format("200601"), 10)
	if err != nil {
		return
	}
	var monthFlows float64
	err = db.QueryRow("select sum(flow) t_cdn_day_statistics where id=? and month>=? and month <=?", a.AppId, monthStart, monthEdn).Scan(&monthFlows)
	if err != nil {
		stmt, err := db.Prepare("insert t_cdn_month_statistics(id,month,flow)values(?,?,?)")
		if err != nil {
			return
		}
		defer stmt.Close()
		err = stmt.Exec(a.AppId, nowMonth, a.Flow)
		if err != nil {
			return
		}
	} else {
		stmt, err := db.Prepare("update t_cdn_month_statistics set flow=? where id=? and month=?")
		if err != nil {
			return
		}
		defer stmt.Close()
		err = stmt.Exec(monthFlows, a.AppId, nowMonth)
		if err != nil {
			return
		}
	}

}
func getFlowSize(bucket_name, startTime, endTime string) (resultSize int64, err error) {
	log.Println("getSize.data", bucket_name, startTime, endTime)
	var result = Result{}
	var data = []Size{}
	u, _ := url.Parse("https://api.upyun.com/v2/nstats")
	q := u.Query()
	q.Set("bucket_name", bucket_name)
	q.Set("start_day", startTime)
	q.Set("end_day", endTime)
	u.RawQuery = q.Encode()

	req := httplib.Get(u.String())
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	req.Header("Content-Type", "application/json")
	req.Header("Authorization", token)
	reqByte, _ := req.Bytes()
	err = json.Unmarshal(reqByte, &result)
	if err != nil {
		return
	}
	var sumFlow float64
	if v, ok := result.Data.([]interface{}); ok {
		for _, vs := range v {
			if size, ok := vs["size"].(float64); ok {
				sumFlow += size
			}
		}
	}

	return int64(sumFlow * 1.1), nil
}

func insertData(t_cdn CdnStatistics) error {
	log.Println("insertData.data", t_cdn)
	db, err := m.GetMySQL()
	if err != nil {
		return err
	}
	defer db.Close()

	sql := "SELECT app_id FROM t_cdn_day_statistics where app_id=? and day=? and type=?"
	row, err := db.Query(sql, t_cdn.AppId, t_cdn.Day, t_cdn.Type)
	if nil != err {
		return err
	}

	if row.Next() {
		sql = "update t_cdn_day_statistics set flow=?,status=? where app_id=? and day=? and type=?"
		stmt, err := db.Prepare(sql)
		if nil != err {
			return err
		}
		res, err := stmt.Exec(t_cdn.Flow, t_cdn.Status, t_cdn.AppId, t_cdn.Day, t_cdn.Type)
		if nil != err {
			return err
		}
		defer stmt.Close()
		_, err = res.RowsAffected()
		return err

	} else {
		sql = "insert into t_cdn_day_statistics (app_id,flow,day,type,status) VALUES(?, ?, ?, ?,?)"
		stmt, err := db.Prepare(sql)
		if nil != err {
			return err
		}
		res, err := stmt.Exec(t_cdn.AppId, t_cdn.Flow, t_cdn.Day, cdn_type, t_cdn.Status)
		if nil != err {
			return err
		}
		defer stmt.Close()
		_, err = res.RowsAffected()
		return err
	}

}

func getMonthFlow() ([]CdnStatistics, error) {
	var t_cdn = CdnStatistics{}
	l_cdn := []CdnStatistics{}

	t1, t2 := getMon()
	db, err := m.GetMySQL()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sql := "select app_id,sum(flow)as flow from  t_cdn_day_statistics where  day>=? and day <=?  group by app_id"

	rows, err := db.Query(sql, t1, t2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&t_cdn.AppId, &t_cdn.Flow); err != nil {
			return nil, err
		}
		l_cdn = append(l_cdn, t_cdn)
	}
	return l_cdn, nil
}

func insertMonthData(t_cdn CdnStatistics) error {
	log.Println("insertMonthData.data", t_cdn)
	t1, t2 := getMon()
	t_now := getTimeStamp(0)
	db, err := m.GetMySQL()
	if err != nil {
		return err
	}
	defer db.Close()

	sql := "SELECT app_id FROM t_cdn_month_statistics where app_id=? and month>=? and month<=?"
	row, err := db.Query(sql, t_cdn.AppId, t1, t2)
	if nil != err {
		return err
	}

	if row.Next() {
		sql = "update t_cdn_month_statistics set flow=?,month=? where app_id=? and month>=? and month<=?"
		stmt, err := db.Prepare(sql)
		if nil != err {
			return err
		}
		res, err := stmt.Exec(t_cdn.Flow, t_now, t_cdn.AppId, t1, t2)
		if nil != err {
			return err
		}
		defer stmt.Close()
		_, err = res.RowsAffected()
		return err

	} else {
		sql = "insert into t_cdn_month_statistics (app_id,flow,type,month) VALUES(?, ?, ?, ?)"
		stmt, err := db.Prepare(sql)
		if nil != err {
			return err
		}
		res, err := stmt.Exec(t_cdn.AppId, t_cdn.Flow, cdn_type, t_now)
		if nil != err {
			return err
		}
		defer stmt.Close()
		_, err = res.RowsAffected()
		return err
	}

}

//获取一个时间区间
func getDay(y, m, d int) (string, string) {
	nTime := time.Now()
	time := nTime.AddDate(y, m, d)
	time2 := nTime.AddDate(y, m, d+1)
	Status = 0
	TimeStamp1, _ = strconv.ParseInt(time.Format("20060102"), 10, 0)
	if TimeStamp2 == 0 {
		TimeStamp2 = TimeStamp1
	}
	return time.Format("2006-01-02") + " 00:00:00", time2.Format("2006-01-02") + " 01:00:00"
}

//获取一个时间区间
func getMon() (int64, int64) {
	now := time.Now()
	t_day := now.Day()
	time1 := now.AddDate(0, 0, -t_day+1)
	time2 := now.AddDate(0, 1, -t_day)
	t_time1, _ := strconv.ParseInt(time1.Format("20060102"), 10, 0)
	t_time2, _ := strconv.ParseInt(time2.Format("20060102"), 10, 0)
	return t_time1, t_time2
}

//获取时间,格式20060102
func getTimeStamp(day int) int64 {
	nTime := time.Now()
	time := nTime.AddDate(0, 0, day)
	t_result, _ := strconv.ParseInt(time.Format("20060102"), 10, 0)
	return t_result
}
