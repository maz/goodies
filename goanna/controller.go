// Package goanna is an MVC toolkit with routing, controllers and sessions
package goanna

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type ControllerInterface interface {
	Init(*http.Request)
	Session() Session
}

type Controller struct {
	Request      Request
	sessionStore SessionStore
}

func (c *Controller) Session() Session {
	if c.Request.session == nil {
		c.Request.session = c.sessionStore.GetSession(&c.Request)
	}
	return c.Request.session
}

// NewController instantiates a new Controller using the given request and session store
func NewController(req *http.Request, sessionStore SessionStore) Controller {
	request := Request{Request: req}
	c := Controller{
		Request:      request,
		sessionStore: sessionStore,
	}
	return c
}

// IsGetRequest() returns whether the request is GET
func (c *Controller) IsGetRequest() bool {
	return c.Request.Method == "GET"
}

// IsGetRequest() returns whether the request is POST
func (c *Controller) IsPostRequest() bool {
	return c.Request.Method == "POST"
}

// RenderView renders a template string using the provided template and vars struct
// and returns the rendered tamplate
func (c *Controller) RenderView(templateStr string, vars interface{}) []byte {
	t, err := template.New("RenderView").Parse(templateStr)
	if err != nil {
		log.Panicln(err)
	}
	return c.RenderTemplate(t, vars)
}

// RenderView renders a template using the provided template and vars struct
// and returns the rendered tamplate
func (c *Controller) RenderTemplate(t *template.Template, vars interface{}) []byte {
	out := bytes.NewBuffer(nil)
	err := t.Execute(out, vars)
	if err != nil {
		log.Panicln(err)
	}
	return out.Bytes()
}

// Render renders a template using the provided template and vars struct
// and returns a response with the rendered template
func (c *Controller) RenderTemplateResponse(t *template.Template, vars interface{}) *OkResponse {
	response := NewResponse()
	response.SetContent(c.RenderTemplate(t, vars))
	return response
}

// Redirect returns a new RedirectResponse
// (HTTP 302)
func (c *Controller) Redirect(urlStr string) *RedirectResponse {
	return NewRedirectResponse(urlStr, 302)
}

// PermanentRedirect returns a new RedirectResponse with a permanent redirect
// (HTTP 301)
func (c *Controller) PermanentRedirect(urlStr string) *RedirectResponse {
	return NewRedirectResponse(urlStr, 301)
}

// Render renders a template using the provided template and vars struct
// and returns a response with the rendered template
func (c *Controller) Render(templateStr string, vars interface{}) *OkResponse {
	response := NewResponse()
	response.SetContent(c.RenderView(templateStr, vars))
	return response
}

// CookieHost returns the request's hostname suitable for use in cookies
// by stripping the port and appending a '.'
func (c *Controller) CookieHost() string {
	return "." + strings.Split(c.Request.Host, ":")[0]
}

func (c *Controller) UrlFor(routeName string, args ...string) *url.URL {
	return UrlFor(routeName, args...)
}

func (c *Controller) RedirectRoute(routeName string, args ...string) *RedirectResponse {
	return NewRedirectResponse(c.UrlFor(routeName, args...).String(), 302)
}
