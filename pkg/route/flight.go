package route

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Flight struct {
	Source      string
	Destination string
}

type Flights []Flight

var FlightHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if r.ContentLength == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	flights := Flights{}
	err = json.Unmarshal(data, &flights)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()
	flight, err := flights.FindBaseFlight(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := []string{flight.Source, flight.Destination}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (flights Flights) FindBaseFlight(ctx context.Context) (f *Flight, err error) {
	m := make(map[string]int)
	for _, flight := range flights {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		if v, ok := m[flight.Source]; ok {
			m[flight.Source] = v + 1
		} else {
			m[flight.Source] = 1
		}
		if v, ok := m[flight.Destination]; ok {
			m[flight.Destination] = v + 1
		} else {
			m[flight.Destination] = 1
		}
	}

	keyWithVal1 := []string{}
	for k, v := range m {
		if v == 1 {
			keyWithVal1 = append(keyWithVal1, k)
		}
	}

	if len(keyWithVal1) != 2 {
		return nil, fmt.Errorf("invalid flights")
	}

	f = &Flight{
		Source:      keyWithVal1[1],
		Destination: keyWithVal1[0],
	}

	for _, flight := range flights {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		if flight.Source == keyWithVal1[0] {
			f.Source = keyWithVal1[0]
			f.Destination = keyWithVal1[1]
			return
		}
	}
	return
}

func (flight *Flight) UnmarshalJSON(data []byte) (err error) {
	var values []json.RawMessage

	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}

	if len(values) != 2 {
		return fmt.Errorf("invalid flight - %v", values)
	}

	if err := json.Unmarshal(values[0], &flight.Source); err != nil {
		return fmt.Errorf("invalid source - %v, err=%s", values[0], err)
	}

	if err := json.Unmarshal(values[1], &flight.Destination); err != nil {
		return fmt.Errorf("invalid destination - %v, err=%s", values[0], err)
	}

	return nil
}
