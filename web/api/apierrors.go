package api

import "net/http"

// Returns an error explaining that a security token is missing
func ErrMissingToken() SimpleApiResult {

	return SimpleApiResult{Result: Body{Result: "Security token missing!"}, Status: http.StatusUnauthorized}
}

// Returns an error explaining that a security token is wrong
func ErrWrongToken() SimpleApiResult {

	return SimpleApiResult{Result: Body{Result: "Security token is wrong, please check either the URL or your settings!"}, Status: http.StatusUnauthorized}
}
