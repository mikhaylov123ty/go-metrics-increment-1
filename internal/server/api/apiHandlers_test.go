package api

import (
	"github.com/stretchr/testify/assert"
	"metrics/internal/storage"
	"net/http/httptest"
	"testing"
)

func TestHandler_Update(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}

	type args struct {
		url         string
		method      string
		contentType string
		metricType  string
		metricName  string
		metricValue string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no url values, error code 400",
			args: args{
				url:    "http://localhost:8080/update",
				method: "POST",
			},
			want: want{
				code: 400,
			},
		},
		{
			name: "wrong content type, error code 400",
			args: args{
				url:         "http://localhost:8080/update",
				method:      "POST",
				contentType: "application/json",
				metricType:  "gauge",
				metricName:  "alloc",
				metricValue: "32.4123",
			},
			want: want{
				code: 400,
			},
		},
		{
			name: "valid gauge",
			args: args{
				url:         "http://localhost:8080/update",
				method:      "POST",
				contentType: "text/plain",
				metricType:  "gauge",
				metricName:  "alloc",
				metricValue: "32.4123",
			},
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "valid counter",
			args: args{
				url:         "http://localhost:8080/update",
				method:      "POST",
				contentType: "text/plain",
				metricType:  "counter",
				metricName:  "alloc",
				metricValue: "3251325234",
			},
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			handler := NewHandler(storage.NewMemoryStorage())
			request := httptest.NewRequest(tt.args.method, tt.args.url, nil)
			request.Header.Add("Content-Type", tt.args.contentType)
			request.SetPathValue("type", tt.args.metricType)
			request.SetPathValue("name", tt.args.metricName)
			request.SetPathValue("value", tt.args.metricValue)

			w := httptest.NewRecorder()

			handler.UpdatePost(w, request)

			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode, "Codes are not equal")
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"), "Content types are note equal")
		})
	}
}