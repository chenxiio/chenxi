package testproj

import (
	"context"

	"github.com/chenxiio/chenxi"
	"github.com/chenxiio/chenxi/api"
	"github.com/chenxiio/chenxi/build/testproj/src/drivers/driver"
	"github.com/chenxiio/chenxi/build/testproj/src/modules"
	"github.com/chenxiio/chenxi/logger"

	"github.com/chenxiio/chenxi/runtime/alarm"
	dlldriver "github.com/chenxiio/chenxi/runtime/driver"
	"github.com/chenxiio/chenxi/runtime/job"
	"github.com/chenxiio/chenxi/runtime/recipe"
)

// func init() {
// 	fmt.Println("testproj init")
// }

func TypesMap(name, ty, path, parm string) (obj any) {
	switch ty {
	case "ioserver":
		ioserver, err := chenxi.NewIOServer(chenxi.CX.Cfg.Basedir, &chenxi.CX.Cfg.IOCfg, chenxi.CX.GetSocketio())
		if err != nil {
			logger.Error(err.Error())
			panic(err)
		}
		return ioserver
	case "recipe":
		rcp := recipe.NewRecipe(chenxi.CX.Cfg.Basedir + "/cfg/recipe/")
		return rcp
	case "job":

		return job.CreateJobInstance()
	case "alarm":

		return alarm.NewAlarms(chenxi.CX.Cfg.Basedir, nil, chenxi.CX.GetSocketio(), logger.GetLog("alarm", "alarm", chenxi.CX.Cfg.Basedir))
	case "classdriver":
		dr := driver.Drivertest{}

		err := dr.Start(context.TODO(), parm)
		if err != nil {
			logger.Error(err.Error())
		}
		return &dr
	// case "driverp":
	// 	dr := driversplugin{}

	// 	err := dr.Start(context.TODO(), parm)
	// 	if err != nil {
	// 		slog.Error(err.Error())
	// 	}
	// 	return dr

	case "dlldriver":
		drv := dlldriver.NewDriverDll(chenxi.CX.Cfg.Basedir+path, parm)
		return drv
	case "plugindriver":
		drv := dlldriver.NewDriverPlugin(chenxi.CX.Cfg.Basedir+path, parm)
		return drv
	case "pmtest":
		pm := &modules.PMTest{Name: name}
		// err := pm.Init(context.TODO(), parm)
		// if err != nil {
		// 	slog.Error(err.Error())
		// 	panic(err)
		// }
		return pm
	case "tmtest":
		pm := &modules.TMTest{Name: name}
		// err := pm.Init(context.TODO(), parm)
		// if err != nil {
		// 	slog.Error(err.Error())
		// }
		return pm
	case "cmtest":
		pm := &modules.CMTest{Name: name}
		// err := pm.Init(context.TODO(), parm)
		// if err != nil {
		// 	slog.Error(err.Error())
		// }
		return pm

	default:
		logger.Error("genmap not fond ", name, ty, path, parm)
		return nil
	}
}

func TypesMapApi(apitp string) (Internal, out any) {
	//apistr := strings.Split(apitp, ".")
	switch apitp {
	case "IOSERVER":
		var w api.IOServerAPIStruct = api.IOServerAPIStruct{}
		return &w.Internal, &w
	case "RECIPE":
		var w api.RecipeApiStruct = api.RecipeApiStruct{}
		return &w.Internal, &w
	case "JOB":
		var out api.JobApiStruct
		return &out.Internal, &out
	case "ALARM":
		var out api.ALMApiStruct
		return &out.Internal, &out
	case "DRIVER":

		return nil, nil
	// case "driverp":
	// 	dr := driversplugin{}

	// 	err := dr.Start(context.TODO(), parm)
	// 	if err != nil {
	// 		slog.Error(err.Error())
	// 	}
	// 	return dr
	case "PM":
		var out api.PMApiStruct
		return &out.Internal, &out
	case "CM":
		var out api.CMApiStruct
		return &out.Internal, &out
	case "TM":
		var out api.TMApiStruct
		return &out.Internal, &out
	case "CUSTOM":
		return nil, nil
	default:
		logger.Error("TypesMapApi not fond ", apitp)
		return nil, nil
	}
}
