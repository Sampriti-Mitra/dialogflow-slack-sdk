package externals

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

type errReader string

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestNewSlackRequest(t *testing.T) {
	var (
		credentialsPath = "i am credentials path"
		payload         = `{
			"type": "dialog_submission",
			"token": "M1AqUUw3FqayAbqNtsGMch72",
			"callback_id": "employee_offsite_1138b",
			"response_url": "https://hooks.slack.com/app/T012AB0A1/123456789/JpmK0yzoZDeRiqfeduTBYXWQ"
		}`
		slackCallback = slack.InteractionCallback{
			Type:        "dialog_submission",
			Token:       "M1AqUUw3FqayAbqNtsGMch72",
			CallbackID:  "employee_offsite_1138b",
			ResponseURL: "https://hooks.slack.com/app/T012AB0A1/123456789/JpmK0yzoZDeRiqfeduTBYXWQ",
		}
	)

	type args struct {
		req             *http.Request
		credentialsPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *SlackRequest
		wantErr bool
	}{
		{
			name: "request is nil",
			args: args{
				credentialsPath: credentialsPath,
			},
			want: &SlackRequest{
				credentials: credentialsPath,
			},
		},
		{
			name: "request body is nil",
			args: args{
				req: &http.Request{
					Form: url.Values{
						"payload": []string{payload},
					},
				},
				credentialsPath: credentialsPath,
			},
			want: &SlackRequest{
				credentials:         credentialsPath,
				InteractionCallback: &slackCallback,
			},
		},
		{
			name: "request is ok",
			args: args{
				req: &http.Request{
					Form: url.Values{
						"payload": []string{payload},
					},
					Body: io.NopCloser(strings.NewReader("Request body")),
				},
				credentialsPath: credentialsPath,
			},
			want: &SlackRequest{
				Body:                []byte("Request body"),
				credentials:         credentialsPath,
				InteractionCallback: &slackCallback,
			},
		},
		{
			name: "payload unmarshal error",
			args: args{
				req: &http.Request{
					Form: url.Values{
						"payload": []string{
							`invalid payload`,
						},
					},
				},
				credentialsPath: credentialsPath,
			},
			want: &SlackRequest{
				credentials: credentialsPath,
			},
		},
		{
			name: "read body error",
			args: args{
				req: &http.Request{
					Body: io.NopCloser(errReader("Request body")),
				},
				credentialsPath: credentialsPath,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSlackRequest(tt.args.req, tt.args.credentialsPath)

			switch tt.wantErr {
			case true:
				assert.Error(t, err)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestSlackRequest_VerifyIncomingSlackRequests(t *testing.T) {
	var (
		validSigningSecret = "e6b19c573432dcc6b075501d51b51bb8"
		validBody          = []byte(`{"token":"aF5ynEYQH0dFN9imlgcADxDB"}`)
		ts                 = fmt.Sprintf("%d", time.Now().Unix())
	)

	secret := hmac.New(sha256.New, []byte(validSigningSecret))
	secret.Write([]byte(fmt.Sprintf("v0:%s:", ts)))
	secret.Write(validBody)

	validHeader := http.Header{
		"X-Slack-Signature":         []string{fmt.Sprintf("v0=%s", hex.EncodeToString(secret.Sum(nil)))},
		"X-Slack-Request-Timestamp": []string{ts},
	}

	type args struct {
		headers       http.Header
		body          []byte
		signingSecret string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "failed secret verification",
			args: args{
				headers: http.Header{},
			},
			want:    http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "failed to write body",
			args: args{
				headers:       http.Header{},
				signingSecret: "abcdefg12345",
			},
			want:    http.StatusInternalServerError,
			wantErr: true,
		},
		{
			name: "failed to ensure secret",
			args: args{
				headers:       http.Header{},
				body:          validBody,
				signingSecret: "abcdefg12345",
			},
			want:    http.StatusUnauthorized,
			wantErr: true,
		},
		{
			name: "everything is ok",
			args: args{
				headers:       validHeader,
				body:          validBody,
				signingSecret: validSigningSecret,
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slackReq := new(SlackRequest)

			got, err := slackReq.VerifyIncomingSlackRequests(tt.args.headers, tt.args.body, tt.args.signingSecret)

			switch tt.wantErr {
			case true:
				assert.Error(t, err)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
