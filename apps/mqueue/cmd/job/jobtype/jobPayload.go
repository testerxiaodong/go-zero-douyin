package jobtype

type DeferPublishVideoPayload struct {
	Title     string
	OwnerId   int64
	OwnerName string
	SectionID int64
	TagIds    []string
	PlayUrl   string
	CoverUrl  string
}
