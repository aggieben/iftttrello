package crashtrello

import (
	"net/http"

	"appengine"
	"appengine/memcache"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

func init() {
	m := martini.Classic()

	m.Use(func(c martini.Context, r *http.Request) {
		c.MapTo(appengine.NewContext(r), (*appengine.Context)(nil))
	})
	m.Use(render.Renderer())

	m.Get("/", index_get)
	m.Post("/configure", binding.MultipartForm(Configuration{}), configure_post)

	http.Handle("/", m)
}

func index_get(r render.Render) {
	r.HTML(200, "index", nil)
}

type Configuration struct {
	AppKey    string `form:"appKey"`
	AppSecret string `form:"appSecret"`
}

func configure_post(c martini.Context, gae appengine.Context, r render.Render, config Configuration) {
	var err error = nil
	var cachedItem *memcache.Item

	cachedItem, err = memcache.Get(gae, "config")
	switch err {
	case memcache.ErrCacheMiss:
		var newItem = memcache.Item{Key: "config", Object: config}
		if err2 := memcache.JSON.Add(gae, &newItem); err2 != nil {
			panic(err2)
		}
		gae.Infof("got item: %v", cachedItem)
	default:
		panic(err)
	}

	r.Redirect("/", 303)
}
