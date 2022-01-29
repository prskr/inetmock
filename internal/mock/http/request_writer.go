//go:generate mockgen --build_flags=--mod=mod -destination=./response_writer.mock.go -package=httpmock net/http ResponseWriter

package httpmock
