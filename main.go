package main

import (
	"github.com/BurntSushi/toml"
	"net/http"
	"github.com/urfave/negroni"
	"strconv"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/yanzay/log"
	. "github.com/cool2645/kotori-ng/config"
	"github.com/cool2645/kotori-ng/model"
	"github.com/cool2645/kotori-ng/handler"
)

var mux = httprouter.New()

func main() {

	_, err := toml.DecodeFile("config.toml", &GlobCfg)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open("sqlite3", GlobCfg.DB_FILE)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Database init done")
	defer db.Close()
	db.AutoMigrate()
	model.Db = db

	mux.GET("/api", handler.Pong)
	mux.ServeFiles("/static/*filepath", http.Dir("static"))
	//mux.NotFound = http.HandlerFunc(NotFoundHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   GlobCfg.ALLOW_ORIGIN,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowCredentials: true,
		//AllowedHeaders: []string{""},
	})
	h := c.Handler(mux)

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("app")))
	n.UseHandler(h)

	http.ListenAndServe(":"+strconv.FormatInt(GlobCfg.PORT, 10), n)
	return
}