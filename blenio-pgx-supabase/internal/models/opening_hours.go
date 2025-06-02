package models

const (
	Monday = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

type TimeRange struct {
	Start string `json:"start" validate:"required,time"` // custom "time" validator (HH:mm)
	End   string `json:"end" validate:"required,time"`   // same here
}

type DaySchedule struct {
	Intervals []TimeRange `json:"intervals" validate:"required,dive"` // dive to validate each TimeRange
}

type DateRange struct {
	StartDate string `json:"start_date" validate:"required,datetime=2006-01-02"`                  // Go time layout for date only
	EndDate   string `json:"end_date" validate:"required,datetime=2006-01-02,gtefield=StartDate"` // EndDate >= StartDate
}

type Exception struct {
	DateRange DateRange   `json:"date_range,omitempty" validate:"omitempty"`                     // optional, dive into DateRange
	Dates     []string    `json:"dates,omitempty" validate:"omitempty,dive,datetime=2006-01-02"` // optional list of dates, validate each
	Closed    bool        `json:"closed"`                                                        // no validation needed
	Hours     []TimeRange `json:"hours,omitempty" validate:"omitempty,dive"`                     // optional list of TimeRanges
	Reason    string      `json:"reason,omitempty"`                                              // optional free text
}

type OpeningHours struct {
	Weekly     map[int8]DaySchedule `json:"weekly" validate:"required,dive"` // required weekly schedule
	Exceptions []Exception          `json:"exceptions,omitempty" validate:"omitempty,dive"`
}
