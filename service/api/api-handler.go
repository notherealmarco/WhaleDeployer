package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/api", rt.getHelloWorld)
	rt.router.GET("/api/projects", rt.wrap(rt.getProjects))
	rt.router.GET("/api/projects/:id", rt.wrap(rt.getProject))
	rt.router.POST("/api/projects", rt.wrap(rt.postProject))
	rt.router.DELETE("/api/projects/:id", rt.wrap(rt.deleteProject))

	rt.router.POST("/api/projects/:id", rt.wrap(rt.buildComposeProject))

	rt.router.GET("/api/projects/:id/logs", rt.wrap(rt.getLogs))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
