package pluginmanager

import (
	"plugin"
	"github.com/pkg/errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"path"
	. "github.com/cool2645/kotori-ng/kotori_plugin"
)

const (
	defaultPluginNum = 10
)

type PluginManager struct {
	pluginDir   string
	pluginCount int
	router      *mux.Router
	db          *gorm.DB
	Plugins     []PluginDescriptor
}

type PluginDescriptor struct {
	Name    string
	Version string
	Path    string
}

func NewPluginManager(path string, r *mux.Router, d *gorm.DB) (*PluginManager) {
	pm := PluginManager{
		pluginDir:   path,
		pluginCount: 0,
		router:      r,
		db:          d,
		Plugins:     make([]PluginDescriptor, defaultPluginNum),
	}
	return &pm
}

func (pm *PluginManager) GetCount() int {
	return pm.pluginCount
}

func (pm *PluginManager) GetPath() string {
	return pm.pluginDir
}

func (pm *PluginManager) loadPlugin(path string) (err error) {
	// Load .so
	p, err := plugin.Open(path)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Failed to load plugin %v: Open file error", path))
	}
	// Get PluginInstance
	pi, err := p.Lookup("PluginInstance")
	if err != nil {
		err = errors.Wrap(err,
			fmt.Sprintf("Failed to load plugin %v: PluginInstance not found", path))
	}
	// Register Plugin
	pm.regPlugin(*pi.(*Plugin))
	return
}

func (pm *PluginManager) regPlugin(p Plugin) () {
	p.RegRouter(pm.router)
	p.InitDB(pm.db)
	pm.Plugins = append(pm.Plugins, PluginDescriptor{
		Name:    p.GetName(),
		Version: p.GetVersion(),
	})
	return
}

func (pm *PluginManager) LoadPlugins() (err error) {
	ps, err := ioutil.ReadDir(pm.pluginDir)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Failed to open plugin dir %v", pm.pluginDir))
	}
	for _, p := range ps {
		if path.Ext(p.Name()) == ".so" {
			pm.loadPlugin(pm.pluginDir + p.Name())
		}
	}
	return
}
