package api

type Body struct {
	Result string
}

// SimpleApiResult Holds data about the result of a simple operation
type SimpleApiResult struct {
	Result Body
	Status int
}
