package models

import (
	"strconv"
	"time"

	"github.com/appointments_api/db"
)

// get appointments (all) GET
func GetAppointments(date string) ([]Appointment, error) {
	var appointments []Appointment

	sqlStr := `SELECT a.id, p.first_name, p.last_name, p.date_of_birth, p.gender, p.contact_number, p.address , COALESCE(doc.id, 0), COALESCE(doc.name, ''), recep.id, recep.name, a.appointment_date, a.notes FROM appointments a
	LEFT JOIN patients p ON p.id = a.patient_id
	LEFT JOIN users doc ON doc.id = a.doctor_id 
	LEFT JOIN users recep ON recep.id = a.recep_id 
	ORDER BY a.appointment_date
	`

	// if date is passed, return results upcoming and todays date appointments
	if date != "" {
		sqlStr += ` WHERE a.appointment_date >= ` + date
	}

	rows, err := db.Db.Query(sqlStr)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var appointment Appointment
		var doc, recep Role
		var p Patient
		var appTimestamp time.Time
		err := rows.Scan(
			&appointment.ID,
			&p.FName,
			&p.LName,
			&p.DOB,
			&p.Gender,
			&p.Contact,
			&p.Address,

			&doc.ID,
			&doc.Name,

			&recep.ID,
			&recep.Name,

			&appTimestamp,
			&appointment.Notes,
		)
		utcTime := appTimestamp.UTC()
		appointment.AppointmentDate = utcTime.Format("2006-01-02")
		appointment.AppointmentTime = utcTime.Format("03:04 PM")

		appointment.DocInfo = doc
		appointment.ReceptionInfo = recep
		appointment.PatientInfo = p

		appointments = append(appointments, appointment)

		if err != nil {
			return nil, err
		}
	}

	return appointments, nil
}

// get single appointment GET
func GetAppointment(id int) (Appointment, error) {
	var appointment Appointment
	var doc, recep Role
	var p Patient

	sqlStr := `SELECT a.id, p.id, p.first_name, p.last_name, p.date_of_birth, p.gender, p.contact_number, p.address , doc.id, doc.name, recep.id, recep.name, a.appointment_date, a.notes FROM appointments a
	LEFT JOIN patients p ON p.id = a.patient_id
	LEFT JOIN users doc ON doc.id = a.doctor_id 
	LEFT JOIN users recep ON recep.id = a.recep_id 
	WHERE a.id = $1
	`
	var appointmentTimestamp time.Time
	err := db.Db.QueryRow(sqlStr, id).Scan(
		&appointment.ID,
		&p.ID,
		&p.FName,
		&p.LName,
		&p.DOB,
		&p.Gender,
		&p.Contact,
		&p.Address,

		&doc.ID,
		&doc.Name,

		&recep.ID,
		&recep.Name,

		&appointmentTimestamp,
		&appointment.Notes,
	)

	if err != nil {
		return appointment, err
	}
	utcTime := appointmentTimestamp.UTC()
	appointment.AppointmentDate = utcTime.Format("2006-01-02")
	appointment.AppointmentTime = utcTime.Format("03:04 PM")

	appointment.DocInfo = doc
	appointment.ReceptionInfo = recep
	appointment.PatientInfo = p

	return appointment, nil

}

// add appointment POST
func AddAppointment(a Appointment, recepID int) error {

	combined := a.AppointmentDate + " " + a.AppointmentTime
	// Assuming the input is in "2006-01-02 15:04:05" format
	layout := "2006-01-02 15:04"
	// parsedTime, err := time.Parse(layout, combined)
	parsedTime, err := time.ParseInLocation(layout, combined, time.UTC)
	if err != nil {
		return err
	}
	sqlStr := `INSERT INTO appointments(patient_id, doctor_id, appointment_date, notes, recep_id) VALUES ($1, $2, $3, $4, $5)`

	var docID *string

	if a.DoctorID != "" {
		docID = &a.DoctorID
	}

	_, err = db.Db.Exec(sqlStr, a.PatientID, docID, parsedTime, a.Notes, recepID)

	return err
}

// edit appointment PATCH
func UpdateAppointment(a Appointment, id int) error {

	sqlStr := `UPDATE appointments SET 
	patient_id = $1, 
	doctor_id = $2, 
	appointment_date = $3, 
	notes = $4 
	WHERE id = $5`

	combined := a.AppointmentDate + " " + a.AppointmentTime
	// Assuming the input is in "2006-01-02 15:04:05" format
	layout := "2006-01-02 15:04"
	// parsedTime, err := time.Parse(layout, combined)
	parsedTime, err := time.ParseInLocation(layout, combined, time.UTC)
	if err != nil {
		return err
	}

	patientID, _ := strconv.Atoi(a.PatientID)
	docID, _ := strconv.Atoi(a.DoctorID)

	_, err = db.Db.Exec(sqlStr, patientID, docID, parsedTime, a.Notes, id)

	return err
}

// delete appointment DELETE
func DeleteAppointment(id int) error {
	sqlStr := `DELETE FROM appointments WHERE id = $1`

	_, err := db.Db.Exec(sqlStr, id)

	return err
}
