//go:generate mockgen --build_flags=--mod=mod -destination=./conn.mock.go -package=audit_mock net Conn

package audit_mock
