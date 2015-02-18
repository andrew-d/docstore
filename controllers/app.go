package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
	"github.com/lann/squirrel"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

// Action defines a standard function signature to use when creating controller
// actions.  A controller action is essentially a method attached to a
// controller.
type Action func(c web.C, w http.ResponseWriter, r *http.Request) error

// AppController is our base controller type
type AppController struct {
	// The database connection for this application
	DB *sqlx.DB

	// SQL builder
	Builder squirrel.StatementBuilderType

	// Decoder for form bodies
	Decoder *schema.Decoder

	// Is the application running in debug mode?
	Debug bool

	// Logger
	Logger *logrus.Entry
}

func (c *AppController) Action(a Action) web.Handler {
	return web.HandlerFunc(func(ctx web.C, w http.ResponseWriter, r *http.Request) {
		err := a(ctx, w, r)
		if err == nil {
			return
		}

		switch v := err.(type) {
		case VError:
			status := v.Status
			if status == 0 {
				status = 500
			}

			c.Logger.WithFields(logrus.Fields{
				"err":        v.Base,
				"status":     status,
				"method":     r.Method,
				"url":        r.URL.String(),
				"request_id": middleware.GetReqID(ctx),
			}).Error(v.Message)

			// (Possibly) convert the error to something readable
			var errVal interface{}

			switch ev := v.Base.(type) {
			case error:
				errVal = ev.Error()
			case []error:
				ret := make([]string, 0, len(ev))
				for _, err := range ev {
					ret = append(ret, err.Error())
				}
				errVal = ret
			default:
				errVal = ev
			}

			c.JSON(w, status, map[string]interface{}{
				"error":   errVal,
				"message": v.Message,
			})

		case RedirectError:
			http.Redirect(w, r, v.Location, v.Code)

		default:
			c.Logger.WithFields(logrus.Fields{
				"err":        v,
				"method":     r.Method,
				"url":        r.URL.String(),
				"request_id": middleware.GetReqID(ctx),
			}).Error("error while processing route")

			c.JSON(w, http.StatusInternalServerError, map[string]interface{}{
				"error":   v.Error(),
				"message": "internal server error",
			})
		}
	})
}
