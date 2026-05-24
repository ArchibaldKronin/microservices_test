package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	auth_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/auth/v1"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	statusv3 "google.golang.org/genproto/googleapis/rpc/status"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

const (
	SessionCookieName = "X-Session-Uuid"

	HeaderUserUUID    = "X-User-Uuid"
	HeaderUserLogin   = "X-User-Login"
	HeaderContentType = "content-type"
	HeaderAuthStatus  = "X-Auth-Status"

	HeaderCookie        = "cookie"
	HeaderAuthorization = "authorization"

	ContentTypeJSON = "application/json"

	AuthStatusDenied = "denied"

	SessionUUIDHeader = "X-Session-Uuid"
)

func (a *apiAuth) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	logger.Info(ctx, "External Authorization Check called")

	sessionUUID, err := a.extractSessionUUID(req)
	if err != nil {
		logger.Error(ctx, "Session extraction failed", zap.Error(err))
		return a.denyRequest("Missing or invalid session", 403), nil
	}

	logger.Info(ctx, "Extracted session_uuid", zap.String("uuid", sessionUUID))

	whoamiResp, err := a.Whoami(ctx, &auth_v1.WhoamiRequest{
		SessionUuid: sessionUUID,
	})
	if err != nil {
		logger.Error(ctx, "Whoami failed", zap.Error(err))
		return a.denyRequest("Invalid session", 403), nil
	}

	return a.allowRequest(whoamiResp, sessionUUID), nil
}

func (a *apiAuth) extractSessionUUID(req *authv3.CheckRequest) (string, error) {
	if req.Attributes == nil || req.Attributes.Request == nil {
		return "", fmt.Errorf("no HTTP request found")
	}

	headers := req.Attributes.Request.Http.Headers

	if cookieHeader, ok := headers[HeaderCookie]; ok && cookieHeader != "" {
		sessionUUID := a.extractSessionFromCookies(cookieHeader)
		if sessionUUID != "" {
			return sessionUUID, nil
		}
	}

	return "", fmt.Errorf("session uuid not found in cookies")
}

func (a *apiAuth) extractSessionFromCookies(cookieHeader string) string {
	req := &http.Request{Header: make(http.Header)}
	req.Header.Add(HeaderCookie, cookieHeader)

	if cookie, err := req.Cookie(SessionCookieName); err == nil {
		var sessionUUID string
		sessionUUID, err = url.QueryUnescape(cookie.Value)
		if err != nil {
			return cookie.Value
		}

		return sessionUUID
	}

	return ""
}

func (a *apiAuth) denyRequest(message string, statusCode int32) *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status: &statusv3.Status{Code: int32(codes.Unauthenticated)},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{
					Code: typev3.StatusCode(statusCode),
				},
				Body: fmt.Sprintf(`{"error": "%s", "timestamp": "%s"}`,
					message, time.Now().Format(time.RFC3339)),
				Headers: []*corev3.HeaderValueOption{
					{
						Header: &corev3.HeaderValue{
							Key:   HeaderContentType,
							Value: ContentTypeJSON,
						},
					},
					{
						Header: &corev3.HeaderValue{
							Key:   HeaderAuthStatus,
							Value: AuthStatusDenied,
						},
					},
				},
			},
		},
	}
}

func (a *apiAuth) allowRequest(whoamiResp *auth_v1.WhoamiResponse, sessionUUID string) *authv3.CheckResponse {
	headers := []*corev3.HeaderValueOption{
		{
			Header: &corev3.HeaderValue{
				Key:   HeaderUserUUID,
				Value: whoamiResp.User.UserUuid,
			},
		},
		{
			Header: &corev3.HeaderValue{
				Key:   HeaderUserLogin,
				Value: whoamiResp.User.Login,
			},
		},
		{
			Header: &corev3.HeaderValue{
				Key:   SessionUUIDHeader,
				Value: sessionUUID,
			},
		},
	}

	return &authv3.CheckResponse{
		Status: &statusv3.Status{Code: 0},
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{
				Headers:         headers,
				HeadersToRemove: []string{HeaderCookie, HeaderAuthorization},
			},
		},
	}
}
