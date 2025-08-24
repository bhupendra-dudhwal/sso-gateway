package egress

import "io"

type HttpClientPorts interface {
	Execute(url, method string, reqPayload io.Reader, resPayload any) error
}
