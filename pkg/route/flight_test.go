package route

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_FlightHandler(t *testing.T) {
	tests := []struct {
		name                string
		method              string
		body                string
		status              int
		source, destination []byte
	}{
		{
			name:        `[["SFO", "EWR"]] => ["SFO", "EWR"]`,
			method:      http.MethodPost,
			body:        `[["SFO", "EWR"]]`,
			status:      http.StatusOK,
			source:      []byte("SFO"),
			destination: []byte("EWR"),
		},
		{
			name:        `[["ATL", "EWR"], ["SFO", "ATL"]] => ["SFO", "EWR"]`,
			method:      http.MethodPost,
			body:        `[["ATL", "EWR"], ["SFO", "ATL"]]`,
			status:      http.StatusOK,
			source:      []byte("SFO"),
			destination: []byte("EWR"),
		},
		{
			name:        `[["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]] => ["SFO", "EWR"]`,
			method:      http.MethodPost,
			body:        `[["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]]`,
			status:      http.StatusOK,
			source:      []byte("SFO"),
			destination: []byte("EWR"),
		},
		{
			name:        `missing body`,
			method:      http.MethodPost,
			body:        ``,
			status:      http.StatusBadRequest,
			source:      nil,
			destination: nil,
		},
		{
			name:        `invalid http method`,
			method:      http.MethodGet,
			body:        `[["SFO", "EWR"]]`,
			status:      http.StatusMethodNotAllowed,
			source:      nil,
			destination: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/track", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			FlightHandler(w, req)
			res := w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			if tt.status != res.StatusCode {
				t.Errorf("expected %v got %v", tt.status, res.StatusCode)
			}
			if !bytes.Contains(data, tt.source) {
				t.Errorf("expected %v got %v", string(tt.source), string(data))
			}
			if !bytes.Contains(data, tt.destination) {
				t.Errorf("expected %v got %v", string(tt.destination), string(data))
			}
		})
	}
}

func Test_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		flight  *Flight
		args    args
		wantErr bool
	}{
		{
			name:   `for ["SFO", "EWR"] return ["SFO", "EWR"]`,
			flight: &Flight{},
			args: args{
				data: []byte(`["SFO", "EWR"]`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.flight.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Flight.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_FindBaseFlight(t *testing.T) {
	tests := []struct {
		name    string
		flights Flights
		wantF   *Flight
		wantErr bool
	}{
		{
			name: `for [["SFO", "EWR"]]  return ["SFO", "EWR"]`,
			flights: Flights{
				Flight{
					Source:      "SFO",
					Destination: "EWR",
				},
			},
			wantF: &Flight{
				Source:      "SFO",
				Destination: "EWR",
			},
			wantErr: false,
		},
		{
			name: `for [["ATL", "EWR"], ["SFO", "ATL"]]  return [["ATL", "EWR"], ["SFO", "ATL"]] return ["SFO", "EWR"]`,
			flights: Flights{
				Flight{
					Source:      "ATL",
					Destination: "EWR",
				},
				Flight{
					Source:      "SFO",
					Destination: "ATL",
				},
			},
			wantF: &Flight{
				Source:      "SFO",
				Destination: "EWR",
			},
			wantErr: false,
		},
		{
			name: `for [["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]] return ["SFO", "EWR"]`,
			flights: Flights{
				Flight{
					Source:      "IND",
					Destination: "EWR",
				},
				Flight{
					Source:      "SFO",
					Destination: "ATL",
				},
				Flight{
					Source:      "GSO",
					Destination: "IND",
				},
				Flight{
					Source:      "ATL",
					Destination: "GSO",
				},
			},
			wantF: &Flight{
				Source:      "SFO",
				Destination: "EWR",
			},
			wantErr: false,
		},
		{
			name: `for [["IND", "EWR"], ["BLR", "DEL"] return error`,
			flights: Flights{
				Flight{
					Source:      "IND",
					Destination: "EWR",
				},
				Flight{
					Source:      "BLR",
					Destination: "DEL",
				},
			},
			wantF:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotF, err := tt.flights.FindBaseFlight(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Flights.FindBaseFlight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantF != nil && (tt.wantF.Source != gotF.Source || tt.wantF.Destination != gotF.Destination) {
				t.Errorf("Flights.FindBaseFlight() = %v, want %v", gotF, tt.wantF)
			}
		})
	}
}
