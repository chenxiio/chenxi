package cfg

import (
	"fmt"
)

type Cfg struct {
	Project Project
	Modules Modules
	IOCfg   IOCfg
	CTCCfg  CTCCfg
	Basedir string
}

func NewCfg(_basedir string) *Cfg {
	return &Cfg{Basedir: _basedir}
}

func (c *Cfg) LoadAll() error {

	c.Project = Project{path: c.Basedir + "cfg/project_cfg.xml"}
	err := c.Project.ReadConfigFile()
	if err != nil {
		return err
	}
	c.Modules = Modules{path: fmt.Sprintf("%scfg/module_cfg.xml", c.Basedir)}
	err = c.Modules.ReadConfigFile()
	if err != nil {
		return err
	}

	c.IOCfg = IOCfg{path: fmt.Sprintf("%scfg/io_cfg.xml", c.Basedir)}
	err = c.IOCfg.ReadConfigFile()
	if err != nil {
		return err
	}

	// // 复制并递增IO配置
	// var newConfigs IODefines = make(IODefines)
	// for i := 0; i < 50; i++ {
	// 	for _, io := range c.IOCfg.Items {
	// 		newIO := io
	// 		newIO.Name = fmt.Sprintf("%s%d", io.Name, i+1)
	// 		newIO.Pw = newIO.Name
	// 		newIO.Pr = newIO.Name
	// 		newConfigs[newIO.Name] = newIO
	// 	}
	// }
	// c.IOCfg.Items = newConfigs
	// err = c.IOCfg.SaveConfigFile()

	c.CTCCfg = CTCCfg{path: fmt.Sprintf("%scfg/ctc.json", c.Basedir)}
	err = c.CTCCfg.ReadFile()
	if err != nil {
		return err
	}
	for _, v := range c.CTCCfg.Interlocking {
		v.Pall = make(map[string]*TAction)
	}
	return nil
}
func (c *Cfg) SaveAll() error {

	err := c.Project.SaveConfigFile()
	if err != nil {
		return err
	}
	err = c.Modules.SaveConfigFile()
	if err != nil {
		return err
	}
	err = c.IOCfg.SaveConfigFile()
	if err != nil {
		return err
	}
	return nil
}
