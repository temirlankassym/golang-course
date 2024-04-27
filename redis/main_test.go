package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestProductHandler(t *testing.T) {
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
		id      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Success",
			args: args{
				id: "1",
			},
			want: "{\"Id\":1,\"Name\":\"Product A\",\"Description\":\"Description of Product A\",\"Price\":10.99}",
		},
		{
			name: "Success",
			args: args{
				id: "2",
			},
			want: "{\"Id\":2,\"Name\":\"Product B\",\"Description\":\"Description of Product B\",\"Price\":15.49}",
		},
		{
			name: "Success",
			args: args{
				id: "3",
			},
			want: "{\"Id\":3,\"Name\":\"Product C\",\"Description\":\"Description of Product C\",\"Price\":5.99}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("http://localhost:8080/product?id=%s", tt.args.id))
			if err != nil {
				log.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("Failed to read response body: %v", err)
			}

			if got := string(body); got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}
