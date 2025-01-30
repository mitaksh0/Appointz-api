package models

// Context key for payload data
type ContextKey string

const PayloadContextKey ContextKey = "jwtPayload"

type Response struct {
	StatusCode int         `json:"status_code,omitempty"`
	Message    interface{} `json:"message,omitempty"`
}

type Credentials struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	RePassword string `json:"re_password,omitempty"`
	Email      string `json:"email"`
	Role       string `json:"role"`
}

type Patient struct {
	ID      int    `json:"id"`
	FName   string `json:"first_name"`
	LName   string `json:"last_name"`
	DOB     string `json:"date_of_birth"`
	Gender  string `json:"gender"`
	Contact string `json:"contact"`
	Address string `json:"address"`
}

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type,omitempty"` // role type: "receptionist" or "doctor"
}

type Appointment struct {
	ID              int     `json:"id"`
	PatientID       int     `json:"patient_id"`       // from patients table
	DoctorID        int     `json:"doctor_id"`        // from users table role doctor
	RecepID         int     `json:"recep_id"`         // from users table role receptionist
	AppointmentDate string  `json:"appointment_date"` // date for appointment
	Notes           string  `json:"notes"`            // about the patient, what the appointment is for
	ReceptionInfo   Role    `json:"reception_info"`
	DocInfo         Role    `json:"doctor_info,omitempty"`
	PatientInfo     Patient `json:"patient_info"`
	Doctors         []Role  `json:"doctors,omitempty"`
}
