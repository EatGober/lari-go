package domain

type Appointment struct {
	AppointmentID          int    `json:"appointmentid`
	PatientID              int    `json:"patientid`
	DepartmentID           int    `json:"departmentid`
	ProviderID             int    `json:"providerid`
	PatientPhone           int    `json:"patientPhone`
	ProviderName           string `json:"providerName`
	ScheduleDateTimeString string `json:"scheduledDateTimeString`
	Duration               int    `json:"duration`
}
