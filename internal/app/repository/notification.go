package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
)

type NotificationPostgres struct {
	db *sqlx.DB
}

func NewNotificationPostgres(db *sqlx.DB) *NotificationPostgres {
	return &NotificationPostgres{db: db}
}

func (r *NotificationPostgres) GetUsersNotification(userId int) ([]model.Notification, error) {
	query := fmt.Sprintf(`SELECT id, title, LEFT(description, 30) AS description, is_viewed, created_at FROM %s WHERE user_id=$1`, notificationsTable)

	var notifications []model.Notification
	err := r.db.Select(&notifications, query, userId)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationPostgres) CreateNotification(userId int, title, description string) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (title, description, user_id) VALUES ($1, $2, $3) RETURNING id`, notificationsTable)

	var id int
	err := r.db.Get(&id, query, title, description, userId)
	return id, err
}

func (r *NotificationPostgres) ReadAndGetNotification(notificationId int) (model.Notification, error) {
	query := fmt.Sprintf(`UPDATE %s SET is_viewed=TRUE WHERE id=$1
                                RETURNING id, title, description, is_viewed, created_at`, notificationsTable)

	var notification model.Notification
	err := r.db.Get(&notification, query, notificationId)
	return notification, err
}

func (r *NotificationPostgres) GetNotificationUserId(notificationId int) (int, error) {
	query := fmt.Sprintf(`SELECT user_id FROM %s WHERE id=$1`, notificationsTable)
	var userId int
	err := r.db.Get(&userId, query, notificationId)
	return userId, err
}

func (r *NotificationPostgres) DeleteNotification(notificationId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, notificationsTable)
	_, err := r.db.Exec(query, notificationId)
	return err
}
