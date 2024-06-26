package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"strconv"
	"strings"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetUser(phoneNumber, password string) (model.User, error) {
	query := fmt.Sprintf(`SELECT id, fullname, COALESCE(email, '') AS email, phone_number,
       							COALESCE(identification_number, '') AS identification_number, is_confirmed, role, created_at
								FROM %s WHERE phone_number=$1 AND password_hash=$2`, usersTable)

	var user model.User
	err := r.db.Get(&user, query, phoneNumber, password)

	return user, err
}

func (r *UserPostgres) GetUserById(userId int) (model.User, error) {
	query := fmt.Sprintf(`SELECT id, fullname, COALESCE(email, '') AS email, phone_number,
       							COALESCE(identification_number, '') AS identification_number, is_confirmed, role, created_at
								 FROM %s WHERE id=$1`, usersTable)

	var user model.User
	err := r.db.Get(&user, query, userId)

	return user, err
}

func (r *UserPostgres) GetUserWithPasswordById(userId int) (model.UserPassword, error) {
	query := fmt.Sprintf(`SELECT id, fullname, COALESCE(email, '') AS email, phone_number,
       							COALESCE(identification_number, '') AS identification_number, password_hash,
       							is_confirmed, role, created_at
								FROM %s WHERE id=$1`, usersTable)

	var user model.UserPassword
	err := r.db.Get(&user, query, userId)
	return user, err
}

func (r *UserPostgres) CreateUser(user model.SignUpInput) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (fullname, phone_number, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable)
	queryCart := fmt.Sprintf(`INSERT INTO %s (user_id) VALUES ($1)`, cartsTable)
	querySubs := fmt.Sprintf(`INSERT INTO %s (user_id) VALUES ($1)`, subscriptionsTable)

	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var id int
	err = tx.Get(&id, query, user.FullName, user.PhoneNumber, user.Password)
	if err != nil {
		return 0, err
	}
	_, err = tx.Exec(queryCart, id)
	if err != nil {
		return 0, err
	}
	_, err = tx.Exec(querySubs, id)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserPostgres) SaveSessionData(refreshToken, fingerprint string, userId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback() // TODO: how can I handle it?

	query := fmt.Sprintf(`DELETE FROM %s WHERE fingerprint=$1 AND user_id=$2`, refreshSessionsTable)
	_, err = tx.Exec(query, fingerprint, userId)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`INSERT INTO %s (refresh_token, fingerprint, user_id) VALUES($1, $2, $3)`, refreshSessionsTable)
	_, err = tx.Exec(query, refreshToken, fingerprint, userId)
	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}

func (r *UserPostgres) GetRefreshSessions(userId int) ([]model.RefreshSession, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1", refreshSessionsTable)
	var refreshSessions []model.RefreshSession
	err := r.db.Select(&refreshSessions, query, userId)
	if err != nil {
		return nil, err
	}
	return refreshSessions, nil
}

func (r *UserPostgres) WipeRefreshSessionsWithFingerprint(fingerprint string, userId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id=$1 AND fingerprint!=$2`, refreshSessionsTable)
	_, err := r.db.Exec(query, userId, fingerprint)
	return err
}

func (r *UserPostgres) WipeRefreshSessions(userId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id=$1`, refreshSessionsTable)
	_, err := r.db.Exec(query, userId)
	return err
}

func (r *UserPostgres) GetRefreshSession(refreshToken string, userId int) (model.RefreshSession, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE refresh_token=$1 AND user_id=$2`, refreshSessionsTable)
	var refreshSession model.RefreshSession
	if err := r.db.Get(&refreshSession, query, refreshToken, userId); err != nil {
		return refreshSession, err
	}
	return refreshSession, nil
}

func (r *UserPostgres) DeleteUserRefreshSession(refreshToken string, userId int) (int, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE refresh_token=$1 AND user_id=$2`, refreshSessionsTable)
	res, err := r.db.Exec(query, refreshToken, userId)
	if err != nil {
		return 500, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 500, err
	}
	if rows <= 0 {
		return 400, errors.New("nothing was deleted")
	}
	return 200, nil
}

func (r *UserPostgres) GetAllUsers(limit, offset int, role, isConfirmed string) (int, []model.User, error) {
	queryCount := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, usersTable)
	var count int
	err := r.db.Get(&count, queryCount)
	if err != nil {
		return 0, nil, err
	}

	query := fmt.Sprintf(`SELECT id, fullname, COALESCE(email, '') AS email, phone_number, COALESCE(identification_number, '') AS identification_number,
       							is_confirmed, role, created_at FROM %s`, usersTable)

	argCounter := 1
	var params []string
	if role != "" || isConfirmed != "" {
		if role != "" {
			params = append(params, `role=$`+strconv.Itoa(argCounter))
			argCounter += 1
		}
		if isConfirmed != "" {
			params = append(params, `is_confirmed=$`+strconv.Itoa(argCounter))
			argCounter += 1
		}
	}
	if len(params) > 0 {
		query = query + ` WHERE ` + strings.Join(params, ` AND `)
	}
	query += ` OFFSET $` + strconv.Itoa(argCounter)
	argCounter += 1
	query += ` LIMIT $` + strconv.Itoa(argCounter)

	var users []model.User
	switch len(params) {
	case 0:
		err = r.db.Select(&users, query, offset, limit)
	case 1:
		if role != "" {
			err = r.db.Select(&users, query, role, offset, limit)
		}
		if isConfirmed != "" {
			err = r.db.Select(&users, query, isConfirmed, offset, limit)
		}
	case 2:
		err = r.db.Select(&users, query, role, isConfirmed, offset, limit)
	}

	if err != nil {
		return 0, nil, err
	}
	return count, users, nil
}

func (r *UserPostgres) DeleteUserById(userId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, usersTable)
	res, err := r.db.Exec(query, userId)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserPostgres) ChangeFullName(userId int, fullName string) error {
	queryChangeFullName := fmt.Sprintf(`UPDATE %s SET fullname=$1 WHERE id=$2`, usersTable)
	queryClearIinBin := fmt.Sprintf(`UPDATE %s SET identification_number=NULL WHERE id=$1`, usersTable)

	tx, err := r.db.Beginx()

	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(queryChangeFullName, fullName, userId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(queryClearIinBin, userId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (r *UserPostgres) ChangeIdentificationNumber(userId int, identificationNumber string) error {
	queryChangeIin := fmt.Sprintf(`UPDATE %s SET identification_number=$1 WHERE id=$2`, usersTable)

	_, err := r.db.Exec(queryChangeIin, identificationNumber, userId)
	return err
}

func (r *UserPostgres) ChangePassword(userId int, passwordHash string) error {
	query := fmt.Sprintf(`UPDATE %s SET password_hash=$1 WHERE id=$2`, usersTable)
	_, err := r.db.Exec(query, passwordHash, userId)
	return err
}

func (r *UserPostgres) ChangePhoneNumber(userId int, phoneNumber, passwordHash string) error {
	query := fmt.Sprintf(`UPDATE %s SET phone_number=$1, password_hash=$2 WHERE id=$3`, usersTable)
	_, err := r.db.Exec(query, phoneNumber, passwordHash, userId)
	return err
}

func (r *UserPostgres) ChangeEmail(userId int, email string) error {
	query := fmt.Sprintf(`UPDATE %s SET email=$1 WHERE id=$2`, usersTable)
	_, err := r.db.Exec(query, email, userId)
	return err
}

func (r *UserPostgres) ChangeUserRole(userId int, role string) error {
	query := fmt.Sprintf(`UPDATE %s SET role=$1 WHERE id=$2`, usersTable)
	_, err := r.db.Exec(query, role, userId)
	return err
}
