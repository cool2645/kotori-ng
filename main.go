package main

import (
	"github.com/BurntSushi/toml"
	. "github.com/cool2645/kotori-ng/config"
	"github.com/cool2645/kotori-ng/handler"
	"github.com/cool2645/kotori-ng/model"
	"github.com/cool2645/kotori-ng/pluginmanager"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"github.com/yanzay/log"
	"net/http"
	"strconv"
)

var (
	BaseApi    = "/api"
	BaseApiVer = "/v1"
	Base       = BaseApi + BaseApiVer
	r          = mux.NewRouter().StrictSlash(true)
	api        = r.PathPrefix(BaseApi).Subrouter()
	v1Api      = api.PathPrefix(BaseApiVer).Subrouter()
)

func InitRouter() {
	// Static files
	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// 404
	r.NotFoundHandler = http.HandlerFunc(handler.NotFoundHandler)
	// Ping
	api.Methods("GET").Path("/").HandlerFunc(handler.Pong)
	v1Api.Methods("GET").Path("/").HandlerFunc(handler.Pong)
}

func main() {
	// Load global config
	_, err := toml.DecodeFile("config.toml", &GlobCfg)
	if err != nil {
		panic(err)
	}
	// Init DB connection
	db, err := gorm.Open("sqlite3", GlobCfg.DB_FILE)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Infof("Database init done")
	db.AutoMigrate()
	model.Db = db

	// Init global router
	InitRouter()
	log.Infof("Router init done")

	// Load plugins
	pm := pluginmanager.NewPluginManager(GlobCfg.PLUGIN_DIR, api, model.Db)
	err = pm.LoadPlugins()
	if err != nil {
		log.Fatal(err)
		return
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   GlobCfg.ALLOW_ORIGIN,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowCredentials: true,
		//AllowedHeaders: []string{""},
	})
	h := c.Handler(r)

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("app")))
	n.UseHandler(h)

	http.ListenAndServe(":"+strconv.FormatInt(GlobCfg.PORT, 10), n)
	return
}
