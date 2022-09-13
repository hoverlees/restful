package restful

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type RestfulServerHandler func(w http.ResponseWriter, r *http.Request, uriParams map[string]string)

type restfulServiceItem struct {
	httpMethod  string
	urlPattern  []string
	handlerFunc RestfulServerHandler
}

type restfulMuxHandler struct {
	http.Handler
	apiRoot  string
	services []restfulServiceItem
}

func urlMatch(pattern []string, uri []string) (bool, map[string]string) {
	if len(pattern) != len(uri) {
		return false, nil
	}
	hasUriParam := false
	for i, patternStr := range pattern {
		uriStr := uri[i]
		if len(patternStr) > 0 && patternStr[0] == '{' {
			hasUriParam = true
			continue
		} else if patternStr != uriStr {
			return false, nil
		}
	}
	if !hasUriParam {
		return true, nil
	}
	uriParams := make(map[string]string)
	for i, patternStr := range pattern {
		uriStr := uri[i]
		patternStrLen := len(patternStr)
		if patternStrLen > 0 && patternStr[0] == '{' {
			key := patternStr[1 : patternStrLen-1]
			uriParams[key] = uriStr
		}
	}
	return true, uriParams
}

func (h *restfulMuxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI[len(h.apiRoot):]
	uriPieces := strings.Split(uri, "/")
	for _, serviceItem := range h.services {
		if r.Method != serviceItem.httpMethod {
			continue
		}
		if ok, uriParams := urlMatch(serviceItem.urlPattern, uriPieces); ok {
			serviceItem.handlerFunc(w, r, uriParams)
			return
		}
	}
	w.WriteHeader(404)

}

type Server struct {
	muxHandlerMap map[string]*restfulMuxHandler
	mux           *http.ServeMux
	server        *http.Server
	addr          string
	readTimeout   time.Duration
	writeTimeout  time.Duration
}

func (s *Server) Start() {
	if s.server != nil {
		fmt.Println("can't start, server has been started.")
		return
	}
	s.server = &http.Server{
		Addr:         s.addr,
		Handler:      s.mux,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}
	s.server.ListenAndServe()
}

func (s *Server) Stop() {
	if s.server == nil {
		fmt.Println("can't stop, server not started.")
		return
	}
	s.server.Close()
	s.server = nil
}

func (s *Server) AddRestfulHandler(apiRoot string, httpMethod string, uri string, handlerFunc RestfulServerHandler) {
	if s.server != nil {
		fmt.Println("can't config timeouts, server has been started.")
		return
	}
	muxHandler, ok := s.muxHandlerMap[apiRoot]
	if !ok {
		muxHandler = &restfulMuxHandler{
			apiRoot: apiRoot,
		}
		s.muxHandlerMap[apiRoot] = muxHandler
		s.mux.Handle(apiRoot+"/", muxHandler)
	}
	item := restfulServiceItem{
		urlPattern:  strings.Split(uri, "/"),
		httpMethod:  httpMethod,
		handlerFunc: handlerFunc,
	}
	muxHandler.services = append(muxHandler.services, item)
}

func (s *Server) SetTimeouts(readTimeout time.Duration, writeTimeout time.Duration) {
	if s.server != nil {
		fmt.Println("can't config timeouts, server has been started.")
		return
	}
	s.readTimeout = readTimeout
	s.writeTimeout = writeTimeout
}

func NewServer(bindAddress string) *Server {
	server := &Server{
		addr:          bindAddress,
		muxHandlerMap: make(map[string]*restfulMuxHandler),
		mux:           &http.ServeMux{},
	}
	return server
}
