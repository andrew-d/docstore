package controllers

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/lann/squirrel"
	"github.com/zenazn/goji/web"
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

	// Is the application running in debug mode?
	Debug bool
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

			// TODO: logging here

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
			c.JSON(w, http.StatusInternalServerError, map[string]interface{}{
				"error":   v.Error(),
				"message": "internal server error",
			})
		}
	})
}
