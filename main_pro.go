//go:build !pro

package main

import "net/http"

func initProRoutes(mux *http.ServeMux) {}
func setupProEngine()                  {}
func isProBuild() bool                 { return false }
