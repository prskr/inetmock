package multiplexing

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

const (
	http2ClientPrefaceLength = len(http2.ClientPreface)
	maxHTTPRead              = 4096

	HTTP2PseudoHeaderMethod    HTTP2PseudoHeader = ":method"
	HTTP2PseudoHeaderScheme    HTTP2PseudoHeader = ":scheme"
	HTTP2PseudoHeaderAuthority HTTP2PseudoHeader = ":authority"
	HTTP2PseudoHeaderPath      HTTP2PseudoHeader = ":path"
)

type HTTP2PseudoHeader string

func parseHTTPVersion(reader io.Reader) (HTTPVersion, *bufio.Reader) {
	bufferedReader := bufio.NewReaderSize(reader, maxHTTPRead)
	b, err := bufferedReader.Peek(http2ClientPrefaceLength)

	if err != nil {
		return HTTPVersionUnknown, nil
	} else if bytes.HasPrefix(b, []byte(http2.ClientPreface)) {
		_, _ = io.CopyN(io.Discard, bufferedReader, int64(http2ClientPrefaceLength))
		return HTTPVersion2, bufferedReader
	}

	for i := 4; i < maxHTTPRead; i++ {
		if b, err = bufferedReader.Peek(i); err != nil {
			return HTTPVersionUnknown, nil
		} else if bytes.HasSuffix(b, []byte("\r\n")) {
			b = b[:len(b)-2]
			break
		}
	}

	_, _, proto, ok := parseRequestLine(string(b))
	if !ok {
		return HTTPVersionUnknown, nil
	}

	major, minor, ok := http.ParseHTTPVersion(proto)
	if !ok {
		return HTTPVersionUnknown, nil
	}

	switch {
	case major == 1 && minor == 0:
		return HTTPVersion10, bufferedReader
	case major == 1 && minor == 1:
		return HTTPVersion11, bufferedReader
	default:
		return HTTPVersionUnknown, nil
	}
}

// parseRequestLine parses "GET /foo HTTP/1.1" into its three parts.
func parseRequestLine(line string) (method, requestURI, proto string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	s2 += s1 + 1
	return line[:s1], line[s1+1 : s2], line[s2+1:], true
}

func parseHTTP1Request(bufReader *bufio.Reader) (*RequestPreface, error) {
	if req, err := http.ReadRequest(bufReader); err != nil {
		return nil, err
	} else {
		return &RequestPreface{
			Path:   req.RequestURI,
			Header: req.Header,
			Method: req.Method,
			Host:   req.Host,
			Scheme: req.URL.Scheme,
		}, nil
	}
}

func parseHTTP2Request(reader *bufio.Reader) (*RequestPreface, error) {
	var (
		req    = new(RequestPreface)
		framer = http2.NewFramer(io.Discard, reader)
		done   bool
	)

	hdec := emittingDecoder(req)

	for !done {
		f, err := framer.ReadFrame()
		if err != nil {
			return nil, err
		}

		switch f := f.(type) {
		case *http2.SettingsFrame:
			// Sender acknowledged the SETTINGS frame. No need to write
			// SETTINGS again.
			if f.IsAck() {
				break
			}
			if err := framer.WriteSettings(); err != nil {
				return nil, err
			}
		case *http2.ContinuationFrame:
			if _, err := hdec.Write(f.HeaderBlockFragment()); err != nil {
				return nil, err
			}
			done = done || f.FrameHeader.Flags&http2.FlagHeadersEndHeaders != 0
		case *http2.HeadersFrame:
			if _, err := hdec.Write(f.HeaderBlockFragment()); err != nil {
				return nil, err
			}
			done = done || f.FrameHeader.Flags&http2.FlagHeadersEndHeaders != 0
		}
	}
	return req, nil
}

func emittingDecoder(req *RequestPreface) *hpack.Decoder {
	const maxDynamicTableSize = uint32(4 << 10)
	if req.Header == nil {
		req.Header = http.Header{}
	}
	return hpack.NewDecoder(maxDynamicTableSize, func(f hpack.HeaderField) {
		switch HTTP2PseudoHeader(f.Name) {
		case HTTP2PseudoHeaderMethod:
			req.Method = f.Value
		case HTTP2PseudoHeaderScheme:
			req.Scheme = f.Value
		case HTTP2PseudoHeaderPath:
			req.Path = f.Value
		case HTTP2PseudoHeaderAuthority:
			req.Host = f.Value
		default:
			req.Header.Set(f.Name, f.Value)
		}
	})
}
