package http

import (
	"net/http"
)

func NewLayoutTemplateOptions(server AssetServer, r *http.Request, w http.ResponseWriter) LayoutTemplateOptions {
	session, _ := store.Get(r, "flash")
	session.Save(r, w)

	return LayoutTemplateOptions{
		FlashMessages: interfaceSliceToStringSlice(session.Flashes()),
	}
}
