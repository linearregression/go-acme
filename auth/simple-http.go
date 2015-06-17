package auth

// SimpleHTTPChallenge
type SimpleHTTPChallenge struct {
	// Token is the expected response when the resource at Path is queried
	Token string
	// Path is the path suffix for a simpleHttp verification
	Path string
}
