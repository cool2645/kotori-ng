package kotori_plugin

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

//type Plugin struct {
//	Name      string
//	Version   string
//	RegRouter func(*mux.Router)
//	InitDB    func(*gorm.DB)
//}

type Plugin interface {
	GetName() string
	GetVersion() string
	RegRouter(*mux.Router)
	InitDB(*gorm.DB)
}
