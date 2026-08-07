package model

// ShortUrlStats 短链接统计
type ShortUrlStats struct {
	ShortUrl               string `db:"short_url" json:"short_url"`
	TodayCount             int    `db:"today_count" json:"today_count"`
	YesterdayCount         int    `db:"yesterday_count" json:"yesterday_count"`
	Last7DaysCount         int    `db:"last_7_days_count" json:"last_7_days_count"`
	MonthlyCount           int    `db:"monthly_count" json:"monthly_count"`
	TotalCount             int    `db:"total_count" json:"total_count"`
	DistinctTodayCount     int    `db:"d_today_count" json:"d_today_count"`
	DistinctYesterdayCount int    `db:"d_yesterday_count" json:"d_yesterday_count"`
	DistinctLast7DaysCount int    `db:"d_last_7_days_count" json:"d_last_7_days_count"`
	DistinctMonthlyCount   int    `db:"d_monthly_count" json:"d_monthly_count"`
	DistinctTotalCount     int    `db:"d_total_count" json:"d_total_count"`
}

// StatsSum 短链接统计
// key:
// 		1.today_count
//		2.yesterday_count
//		3.last_7_days_count
//		4.monthly_count
//		5.d_today_count
//		6.d_yesterday_count
//		7.d_last_7_days_count
//		8.d_monthly_count
//      9.total_pv
//      10.total_uv
type StatsSum struct {
	Key   string `db:"stats_key"`
	Value int    `db:"stats_value"`
}

// ShortUrlDailyStats 短链接日统计
type ShortUrlDailyStats struct {
	ID       int    `db:"id" json:"id"`
	ShortUrl string `db:"short_url" json:"short_url"`
	Date     string `db:"date" json:"date"`
	Pv       int    `db:"pv" json:"pv"`
	Uv       int    `db:"uv" json:"uv"`
}
