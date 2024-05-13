package model

type Environment struct {
	ID                string
	UserID            string
	HardwareInput     string
	SnapshotURL       string
	Images            []string
	ServiceNames      []string
	NumberOfInstances int32
}
