package bot

type MessageRecord struct {
	GroupCode  int64 `json:"group_code"`
	EventId    int32 `json:"event_id"`
	InternalId int32 `json:"internal_id"`
	Time       int64
}
