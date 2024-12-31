package api

import (
	"context"

	"github.com/chenxiio/chenxi/models"
)

type ALMApi interface {
	Insert(ctx context.Context, alarm *models.Alarm) error                                //perm:none
	ClearAlarms(ctx context.Context, aid int64, cleartype int) error                      //perm:none
	GetAlarms(ctx context.Context) ([]models.Alarm, error)                                //perm:none
	GetAlarmsHistory(ctx context.Context, start int64, end int64) ([]models.Alarm, error) //perm:none
}
