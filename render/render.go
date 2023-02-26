package render

import "net/http"

type Render interface {
	Render(w http.ResponseWriter, status int) error
}

type RenderFunc func(w http.ResponseWriter, status int) error

func (f RenderFunc) Render(w http.ResponseWriter, status int) error {
	return f(w, status)
}

func HTML(buf []byte) Render {
	return RenderFunc(func(w http.ResponseWriter, status int) error {
		w.Header().Set("content-type", "text/html; charset=utf-8")
		w.WriteHeader(status)
		_, err := w.Write(buf)
		return err
	})
}

func CompressedHTML(buf []byte) Render {
	return RenderFunc(func(w http.ResponseWriter, status int) error {
		w.Header().Set("content-type", "text/html; charset=utf-8")
		w.Header().Set("content-encoding", "gzip")
		w.WriteHeader(status)
		_, err := w.Write(buf)
		return err
	})
}
