package domain

import (
	"encoding/json"
	"log"
)

type CancelledAppointment struct {
	ScheduleDateTimeString string `json:"startTimeISO"`
	Duration               int    `json:"duration"`
}

type Appointment struct {
	AppointmentID          int    `json:"appointmentid"`
	PatientID              int    `json:"patientid"`
	DepartmentID           int    `json:"departmentid"`
	ProviderID             int    `json:"providerid"`
	PatientPhone           int    `json:"patientPhone"`
	ProviderName           string `json:"providerName"`
	ScheduleDateTimeString string `json:"scheduledDateTimeString"`
	Duration               int    `json:"duration"`
}

func ConstructCancelledAppointment(jsonInput string) CancelledAppointment {
	var appt CancelledAppointment

	if err := json.Unmarshal([]byte(jsonInput), &appt); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	return appt
}

func ConstructApppointment(jsonInput string) Appointment {
	var appt Appointment

	if err := json.Unmarshal([]byte(jsonInput), &appt); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	return appt
}

func ConstructApptLists(jsonInput string) []Appointment {
	var appts []Appointment

	if err := json.Unmarshal([]byte(jsonInput), &appts); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return appts
}
