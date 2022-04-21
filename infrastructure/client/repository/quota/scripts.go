package quota

const (
	selectQuota = "SELECT limit_calls FROM quota where source_ip = ? and target_path = ? "
)
