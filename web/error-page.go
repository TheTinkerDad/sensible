package web

// ErrorPage returns a Page object displaying an error message
func ErrorPage() *Page {

	return &Page{Title: "Uh oh...", Body: "The infamous 404 strikes again!"}
}
