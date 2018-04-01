package pluginmanager

import (
	"plugin"
	"github.com/pkg/errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"path"
	. "github.com/cool2645/kotori-ng/kotoriplugin"
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
	pm.regPlugin(*pi.(*Plugin), path)
	return
}

func (pm *PluginManager) regPlugin(p Plugin, ppath string) () {
	filename := path.Base(ppath)
	filename = path.Base(ppath)[0 : len(filename)-len(path.Ext(ppath))]
	sr := pm.router.PathPrefix("/" + filename).Subrouter()
	//sr.Methods("GET").Path("/").HandlerFunc(handler.Pong)
	p.RegRouter(sr)
	p.InitDB(pm.db)
	pm.Plugins = append(pm.Plugins, PluginDescriptor{
		Name:    p.GetName(),
		Version: p.GetVersion(),
		Path:    ppath,
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
