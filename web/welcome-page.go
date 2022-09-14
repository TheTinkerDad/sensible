package web

// WelcomePage returns a Page object displaying a welcome message
func WelcomePage() *Page {

	return &Page{Title: "Hello", Body: "Welcome to \"Sensible\"(tm)!"}
}
