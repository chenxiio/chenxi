package driver

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

type Drivertest struct {
	mymap map[string]interface{}
}

func (dp *Drivertest) Start(ctx context.Context, parm string) error {
	if dp.mymap == nil {
		dp.mymap = make(map[string]interface{})
	}
	rand.Seed(time.Now().UnixNano())
	//dp.mymap[parm] = nil
	return nil
}
func (dp *Drivertest) Stop(ctx context.Context, parm string) error {
	// if dp.mymap == nil {
	// 	return nil
	// }
	// delete(dp.mymap, parm)
	return nil
}
func (dp *Drivertest) ReadInt(ctx context.Context, parm string) (int32, error) {
	if dp.mymap == nil {
		return 0, errors.New("map is not initialized")
	}
	val0, ok := dp.mymap[parm]
	if !ok {
		return rand.Int31(), nil
	}
	val, ok := val0.(int32)
	if !ok {
		return 0, errors.New("value is not an integer")
	}
	return val, nil
}
func (dp *Drivertest) ReadString(ctx context.Context, parm string) (string, error) {
	if dp.mymap == nil {
		return "", errors.New("map is not initialized")
	}
	val, ok := dp.mymap[parm].(string)
	if !ok {
		return "value is not a string", nil
	}
	return val, nil
}
func (dp *Drivertest) ReadDouble(ctx context.Context, parm string) (float64, error) {
	if dp.mymap == nil {
		return 0, errors.New("map is not initialized")
	}
	val0, ok := dp.mymap[parm]
	if !ok {
		return rand.Float64(), nil
	}
	val, ok := val0.(float64)
	if !ok {
		return 0, errors.New("value is not a float")
	}
	return val, nil
}
func (dp *Drivertest) WriteInt(ctx context.Context, parm string, value int32) error {
	if dp.mymap == nil {
		return errors.New("map is not initialized")
	}
	dp.mymap[parm] = value
	return nil
}
func (dp *Drivertest) WriteString(ctx context.Context, parm string, value string) error {
	if dp.mymap == nil {
		return errors.New("map is not initialized")
	}
	dp.mymap[parm] = value
	return nil
}
func (dp *Drivertest) WriteDouble(ctx context.Context, parm string, value float64) error {
	if dp.mymap == nil {
		return errors.New("map is not initialized")
	}
	dp.mymap[parm] = value
	return nil
}
