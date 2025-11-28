package domain

import "time"

type TotalSnapshot struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type SegmentSnapshot struct {
	FlightNumber string    `json:"flightNumber"`
	Departure    string    `json:"from"`
	Arrival      string    `json:"to"`
	DepartTime   time.Time `json:"depart"`
	ArriveTime   time.Time `json:"arrive"`
}

type FlightSnapshot struct {
	ID            string            `json:"id"`
	Status        string            `json:"status"`
	PassengerName string            `json:"passengerName"`
	Segments      []SegmentSnapshot `json:"segments"`
	Total         TotalSnapshot     `json:"total"`
	Source        string            `json:"source"`
}

type FlightsSnapshot []FlightSnapshot

func (t Total) Snapshot() TotalSnapshot {
	return TotalSnapshot{
		Amount:   t.amount,
		Currency: t.currency,
	}
}

func (s Segment) Snapshot() SegmentSnapshot {
	return SegmentSnapshot{
		FlightNumber: s.flightNumber,
		Departure:    s.departure,
		Arrival:      s.arrival,
		DepartTime:   s.departTime,
		ArriveTime:   s.arriveTime,
	}
}

func (f Flight) Snapshot() FlightSnapshot {
	segs := make([]SegmentSnapshot, len(f.segments))
	for i, s := range f.segments {
		segs[i] = s.Snapshot()
	}
	return FlightSnapshot{
		ID:            f.id,
		Status:        f.status,
		PassengerName: f.passengerName,
		Segments:      segs,
		Total:         f.total.Snapshot(),
		Source:        f.source,
	}
}

func (fs Flights) ToSnapshot() FlightsSnapshot {
	out := make(FlightsSnapshot, len(fs))
	for i, f := range fs {
		out[i] = f.Snapshot()
	}
	return out
}
