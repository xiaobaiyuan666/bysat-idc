package audit

type Entry struct {
	ID          int64          `json:"id"`
	ActorType   string         `json:"actorType"`
	ActorID     int64          `json:"actorId"`
	Actor       string         `json:"actor"`
	Action      string         `json:"action"`
	TargetType  string         `json:"targetType"`
	TargetID    int64          `json:"targetId"`
	Target      string         `json:"target"`
	RequestID   string         `json:"requestId,omitempty"`
	Description string         `json:"description"`
	Payload     map[string]any `json:"payload,omitempty"`
	CreatedAt   string         `json:"createdAt"`
}

type Repository interface {
	Record(entry Entry) Entry
	List() []Entry
	ListByTarget(targetType string, targetID int64) []Entry
}
