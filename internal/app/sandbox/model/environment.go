package model

type Requirements struct {
	ID                string
	UserID            string
	HardwareInput     string
	SnapshotURL       string
	Images            []string
	ServiceNames      []string
	NumberOfInstances int32
}
