package vote

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	type args struct {
		httpClient *http.Client
		baseURL    *string
	}

	sampleClient := &http.Client{}
	sampleURL := "http://azert.y"

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "all nil",
			args: args{
				httpClient: nil,
				baseURL:    nil,
			},
			wantErr: false,
		},
		{
			name: "all set",
			args: args{
				httpClient: sampleClient,
				baseURL:    &sampleURL,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.httpClient, tt.args.baseURL)
			assert.Equal(t, err != nil, tt.wantErr)
			assert.NotNil(t, got)
			assert.NotNil(t, got.client)
			assert.NotNil(t, got.client.Jar)
			assert.NotNil(t, got.urlFormal)

			if tt.args.baseURL == nil {
				assert.Equal(t, defaultURL, got.url)
				assert.Equal(t, defaultURL, got.urlFormal.String())
			} else {
				assert.Equal(t, *tt.args.baseURL, got.url)
			}

			if tt.args.httpClient != nil {
				assert.Equal(t, tt.args.httpClient, got.client)
			}
		})
	}
}
