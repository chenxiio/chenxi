package models

type Alarm struct {
	Sn         int64
	Module     string
	Level      int
	Aid        int
	ClearType  int
	Text       string
	CreateTime int64
	OffTime    int64
}
