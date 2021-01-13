package mock_test

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"gitlab.com/inetmock/inetmock/internal/endpoint/handler/http/mock"
	api_mock "gitlab.com/inetmock/inetmock/internal/mock/api"
	audit_mock "gitlab.com/inetmock/inetmock/internal/mock/audit"
	"gitlab.com/inetmock/inetmock/pkg/api"
	"gitlab.com/inetmock/inetmock/pkg/config"
	"gitlab.com/inetmock/inetmock/pkg/logging"
	"go.uber.org/zap"
)

var (
	testLogger          = logging.NewLogger(zap.NewNop())
	availableExtensions = []string{"gif", "html", "ico", "jpg", "png", "txt"}
	charSet             = "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func Benchmark_httpHandler(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	listenPort := randomHighPort()
	_, handler := setupHandler(b, ctrl, listenPort)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		extension := availableExtensions[rand.Intn(len(availableExtensions))]
		if resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s.%s", listenPort, randomString(15), extension)); err != nil {
			b.Error(err)
		} else if resp.StatusCode != 200 {
			b.Error("")
		}
	}

	defer handler.Shutdown(context.Background())
}

func randomString(length int) (result string) {
	buffer := strings.Builder{}
	for i := 0; i < length; i++ {
		buffer.WriteByte(charSet[rand.Intn(len(charSet))])
	}
	return buffer.String()
}

func setupHandler(b *testing.B, ctrl *gomock.Controller, listenPort uint16) (api.HandlerRegistry, api.ProtocolHandler) {
	b.Helper()

	registry := api.NewHandlerRegistry()
	if err := mock.AddHTTPMock(registry); err != nil {
		b.Errorf("AddHTTPMock() error = %v", err)
	}
	handler, ok := registry.HandlerForName("http_mock")
	if !ok {
		b.Error("handler not registered")
	}

	emitter := audit_mock.NewMockEmitter(ctrl)
	emitter.EXPECT().Emit(gomock.Any()).AnyTimes()

	mockApp := api_mock.NewMockPluginContext(ctrl)
	mockApp.EXPECT().
		Logger().
		Return(testLogger)

	mockApp.EXPECT().
		Audit().
		Return(emitter)

	v := viper.New()
	v.Set("rules", []map[string]string{
		{
			"pattern":  ".*\\.(?i)gif",
			"response": "./../../assets/fakeFiles/default.gif",
		},
		{
			"pattern":  ".*\\.(?i)html",
			"response": "./../../assets/fakeFiles/default.html",
		},
		{
			"pattern":  ".*\\.(?i)ico",
			"response": "./../../assets/fakeFiles/default.ico",
		},
		{
			"pattern":  ".*\\.(?i)jpg",
			"response": "./../../assets/fakeFiles/default.jpg",
		},
		{
			"pattern":  ".*\\.(?i)png",
			"response": "./../../assets/fakeFiles/default.png",
		},
		{
			"pattern":  ".*\\.(?i)txt",
			"response": "./../../assets/fakeFiles/default.txt",
		},
	})

	handlerConfig := config.HandlerConfig{
		HandlerName:   "http_test",
		Port:          listenPort,
		ListenAddress: "localhost",
		Options:       v,
	}

	if err := handler.Start(mockApp, handlerConfig); err != nil {
		b.Error(err)
		b.FailNow()
	}

	return registry, handler
}

func randomHighPort() uint16 {
	var err error
	var listener net.Listener
	defer func() {
		if listener != nil {
			_ = listener.Close()
		}
	}()

	for {
		if listener, err = net.Listen("tcp", ":0"); err == nil {
			parts := strings.Split(listener.Addr().String(), ":")
			port, _ := strconv.Atoi(parts[len(parts)-1])
			return uint16(port)
		}
	}
}
