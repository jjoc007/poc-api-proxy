package call

const (
	insertCall = `INSERT INTO call_log (id, source_ip, target_path, duration, creation_date) VALUES(?, ?, ?, ?, ?)`
)
