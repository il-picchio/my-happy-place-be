package validation

import (
	"blenioviva/internal/models"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/biter777/countries"
	"github.com/go-playground/validator/v10"
)

var _swissZipRegex = regexp.MustCompile(`^\d{4}$`)
var _timeRegex = regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

func _RegisterCustomValidators() {
	// Register custom validation functions if needed
	// e.g., Validate.RegisterValidation("custom_tag", CustomValidationFunction)

	// Register translation functions if needed
	// e.g., Validate.RegisterTranslation("custom_tag", translator, func(ut ut.Translator) error { ... })
	Raw.RegisterValidation("ch_zip", func(fl validator.FieldLevel) bool {
		zip := fl.Field().String()
		return _swissZipRegex.MatchString(zip)
	})

	Raw.RegisterValidation("country", func(fl validator.FieldLevel) bool {
		code := strings.ToUpper(fl.Field().String())
		country := countries.ByName(code)
		return country != countries.Unknown
	})

	Raw.RegisterValidation("time", func(fl validator.FieldLevel) bool {
		timeStr := fl.Field().String()
		if !_timeRegex.MatchString(timeStr) {
			return false
		}
		parts := strings.Split(timeStr, ":")
		if len(parts) != 2 {
			return false
		}
		hour, err := strconv.Atoi(parts[0])
		if err != nil || hour < 0 || hour > 23 {
			return false
		}

		minute, err := strconv.Atoi(parts[1])
		if err != nil || minute < 0 || minute > 59 {
			return false
		}
		return true
	})

	Raw.RegisterStructValidation(func(sl validator.StructLevel) {
		tr, ok := sl.Current().Interface().(models.TimeRange)
		if !ok {
			return
		}

		start, err1 := time.Parse("15:04", tr.Start)
		end, err2 := time.Parse("15:04", tr.End)

		if err1 == nil && err2 == nil && !start.Before(end) {
			// Report on the "End" field with tag "timerange"
			sl.ReportError(tr, "", "", "timerange", "")
		}
	}, models.TimeRange{})
}
