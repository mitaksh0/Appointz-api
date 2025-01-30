package models

import "github.com/appointments_api/db"

// get appointments (all) GET
func GetAppointments(date string) ([]Appointment, error) {
	var appointments []Appointment

	sqlStr := `SELECT p.first_name, p.last_name, p.date_of_birth, p.gender, p.contact_number, p.address , doc.id, doc.name, recep.id, recep.name, a.appointment_date, a.notes FROM appointments a
	LEFT JOIN patients p ON p.id = a.patient_id
	LEFT JOIN users doc ON doc.id = a.doctor_id 
	LEFT JOIN users recep ON recep.id = a.recep_id 
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
		err := rows.Scan(
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

			&appointment.AppointmentDate,
			&appointment.Notes,
		)

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

	sqlStr := `SELECT p.first_name, p.last_name, p.date_of_birth, p.gender, p.contact_number, p.address , doc.id, doc.name, recep.id, recep.name, a.appointment_date, a.notes FROM appointments a
	LEFT JOIN patients p ON p.id = a.patient_id
	LEFT JOIN users doc ON doc.id = a.doctor_id 
	LEFT JOIN users recep ON recep.id = a.recep_id 
	WHERE a.id = $1
	`
	err := db.Db.QueryRow(sqlStr, id).Scan(
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

		&appointment.AppointmentDate,
		&appointment.Notes,
	)

	if err != nil {
		return appointment, err
	}

	appointment.DocInfo = doc
	appointment.ReceptionInfo = recep
	appointment.PatientInfo = p

	return appointment, nil

}

// add appointment POST
func AddAppointment(a Appointment) error {

	sqlStr := `INSERT INTO appointments(patient_id, doctor_id, appointment_date, notes, recep_id) VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Db.Exec(sqlStr, a.PatientID, a.DoctorID, a.AppointmentDate, a.Notes, a.RecepID)

	return err
}

// edit appointment PATCH
func UpdateAppointment(a Appointment, id int) error {

	sqlStr := `UPDATE appointments SET 
	patient_id = $1, 
	doctor_id = $2, 
	appointment_date = $3, 
	notes = $4, 
	recep_id = $5
	WHERE id = $6`

	_, err := db.Db.Exec(sqlStr, a.PatientID, a.DoctorID, a.AppointmentDate, a.Notes, a.RecepID, id)

	return err
}

// delete appointment DELETE
func DeleteAppointment(id int) error {
	sqlStr := `DELETE FROM appointments WHERE id = $1`

	_, err := db.Db.Exec(sqlStr, id)

	return err
}
