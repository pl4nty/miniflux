// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package core

import (
	"net/http"
	"time"

	"github.com/miniflux/miniflux2/server/template"
)

// Response handles HTTP responses.
type Response struct {
	writer   http.ResponseWriter
	request  *http.Request
	template *template.TemplateEngine
}

// SetCookie send a cookie to the client.
func (r *Response) SetCookie(cookie *http.Cookie) {
	http.SetCookie(r.writer, cookie)
}

// JSON returns a JSONResponse.
func (r *Response) JSON() *JSONResponse {
	r.commonHeaders()
	return &JSONResponse{writer: r.writer, request: r.request}
}

// HTML returns a HTMLResponse.
func (r *Response) HTML() *HTMLResponse {
	r.commonHeaders()
	return &HTMLResponse{writer: r.writer, request: r.request, template: r.template}
}

// XML returns a XMLResponse.
func (r *Response) XML() *XMLResponse {
	r.commonHeaders()
	return &XMLResponse{writer: r.writer, request: r.request}
}

// Redirect redirects the user to another location.
func (r *Response) Redirect(path string) {
	http.Redirect(r.writer, r.request, path, http.StatusFound)
}

// Cache returns a response with caching headers.
func (r *Response) Cache(mimeType, etag string, content []byte, duration time.Duration) {
	r.writer.Header().Set("Content-Type", mimeType)
	r.writer.Header().Set("Etag", etag)
	r.writer.Header().Set("Cache-Control", "public")
	r.writer.Header().Set("Expires", time.Now().Add(duration).Format(time.RFC1123))

	if etag == r.request.Header.Get("If-None-Match") {
		r.writer.WriteHeader(http.StatusNotModified)
	} else {
		r.writer.Write(content)
	}
}

func (r *Response) commonHeaders() {
	r.writer.Header().Set("X-XSS-Protection", "1; mode=block")
	r.writer.Header().Set("X-Content-Type-Options", "nosniff")
	r.writer.Header().Set("X-Frame-Options", "DENY")
}

// NewResponse returns a new Response.
func NewResponse(w http.ResponseWriter, r *http.Request, template *template.TemplateEngine) *Response {
	return &Response{writer: w, request: r, template: template}
}