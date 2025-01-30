package models

import "github.com/appointments_api/db"

// get all patients
func GetPatients() ([]Patient, error) {
	sqlStr := `SELECT id, first_name, last_name, date_of_birth, gender, contact_number, address FROM patients
	ORDER BY first_name, last_name`
	rows, err := db.Db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var patients []Patient

	for rows.Next() {
		var patient Patient
		err := rows.Scan(
			&patient.ID,
			&patient.FName,
			&patient.LName,
			&patient.DOB,
			&patient.Gender,
			&patient.Contact,
			&patient.Address,
		)

		if err != nil {
			return nil, err
		}

		patients = append(patients, patient)
	}

	return patients, nil
}

// get single patient
func GetPatient(id int) (Patient, error) {
	var patient Patient

	sqlStr := `SELECT id, first_name, last_name, date_of_birth, gender, contact_number, address FROM patients WHERE id = $1`
	err := db.Db.QueryRow(sqlStr, id).Scan(
		&patient.ID,
		&patient.FName,
		&patient.LName,
		&patient.DOB,
		&patient.Gender,
		&patient.Contact,
		&patient.Address,
	)

	if err != nil {
		return patient, err
	}

	return patient, nil
}

// insert patient record
func InsertPatient(p Patient) error {

	sqlStr := `INSERT INTO patients(first_name, last_name, date_of_birth, gender, contact_number, address) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Db.Exec(sqlStr, p.FName, p.LName, p.DOB, p.Gender, p.Contact, p.Address)

	return err
}

// delete patient
func DeletePatient(id int) error {
	sqlStr := `DELETE FROM patients WHERE id = $1`
	_, err := db.Db.Exec(sqlStr, id)

	return err
}

// update patient
func UpdatePatient(p Patient, id int) error {

	sqlStr := `UPDATE patients SET 
	first_name = $1, 
	last_name = $2, 
	date_of_birth = $3, 
	gender = $4, 
	contact_number = $5, 
	address = $6 
	WHERE id = $7`

	_, err := db.Db.Exec(sqlStr, p.FName, p.LName, p.DOB, p.Gender, p.Contact, p.Address, id)

	return err
}
