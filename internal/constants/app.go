package constants

type Environment string

const (
	Production  Environment = "production"
	Development Environment = "development"
)

type Header string

const (
	ContentType     Header = "Content-Type"
	ContentEncoding Header = "Content-Encoding"
)

type ContentTypes string

const (
	Json ContentTypes = "application/json"
)

type Encoding string

const (
	Gzip Encoding = "gzip"
)
