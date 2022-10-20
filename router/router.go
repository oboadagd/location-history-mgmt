// Package router defines request urls of microservice api
package router

import (
	"github.com/labstack/echo/v4"
	middleKit "github.com/oboadagd/kit-go/middleware/echo"
	"github.com/oboadagd/location-history-mgmt/controller"
)

// Router represents the router layer.
type Router struct {
	server             *echo.Echo                                // *echo.Echo that has embedded a http server
	locationController controller.LocationControllerInterface    // controller layer
	errorMiddleware    middleKit.ErrorHandlerMiddlewareInterface // error handle middleware
}

// NewRouter initializes router layer
func NewRouter(
	server *echo.Echo,
	locationController controller.LocationControllerInterface,
	errorMiddleware middleKit.ErrorHandlerMiddlewareInterface,
) *Router {
	return &Router{
		server,
		locationController,
		errorMiddleware,
	}
}

// Init implements request urls definition
func (r *Router) Init() {

	basePath := r.server.Group("/location-history-mgmt")

	locations := basePath.Group("/locations", r.errorMiddleware.HandlerError)
	{
		locations.GET("/distance/:userName/:initialDate/:finalDate", r.locationController.GetDistanceTraveled)
		locations.GET("/distance/:userName", r.locationController.GetDistanceTraveled)
	}
}
