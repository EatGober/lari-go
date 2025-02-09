package scheduler

import (
	"lari-go/internal/domain"
	"regexp"
	"strconv"
	"time"
)

func CreateSlot(appoints []domain.Appointment) string {

	currentTime := time.Now().String()
	re := regexp.MustCompile(`\D`)
	apptID := re.ReplaceAllString(currentTime, "")

	var pids []string

	for _, appt := range appoints {
		pids = append(pids, strconv.Itoa(appt.PatientID))
	}

	saveList(apptID, pids)

	return apptID
}
