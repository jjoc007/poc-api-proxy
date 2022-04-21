package calls_stats_summary

const (
	selectCallsStatsSummary = "select total_calls  from calls_stats_summary where source_ip = ? and target_path = ? and `year` = ? and `month` = ?"
	insertCallsStatsSummary = "INSERT INTO calls_stats_summary (source_ip, target_path, `year`, `month`) VALUES(?, ?, ?, ?) ON DUPLICATE KEY UPDATE total_calls = total_calls +1;"
)
