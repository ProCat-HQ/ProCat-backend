package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
)

type CategoryPostgres struct {
	db *sqlx.DB
}

func NewCategoryPostgres(db *sqlx.DB) *CategoryPostgres {
	return &CategoryPostgres{db: db}
}

func (r *CategoryPostgres) CreateCategory(categoryParentId int, name string) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (name, parent_id) VALUES ($1, $2) RETURNING id`, categoriesTable)

	var id int
	err := r.db.Get(&id, query, name, categoryParentId)
	return id, err
}

func (r *CategoryPostgres) ChangeCategory(categoryId int, name string) error {
	query := fmt.Sprintf(`UPDATE %s SET name=$1 WHERE id=$2`, categoriesTable)

	_, err := r.db.Exec(query, name, categoryId)
	return err
}

func (r *CategoryPostgres) GetCategoriesForParent(categoryParentId int) ([]model.Category, error) {
	query := fmt.Sprintf(`SELECT id, COALESCE(name, '') AS name,
       							 COALESCE(parent_id, -1) AS parent_id FROM %s WHERE parent_id=$1`, categoriesTable)

	var categories []model.Category

	err := r.db.Select(&categories, query, categoryParentId)
	return categories, err
}

func (r *CategoryPostgres) DeleteCategory(categoryId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, categoriesTable)
	_, err := r.db.Exec(query, categoryId)
	return err
}

func (r *CategoryPostgres) GetCategoryRoute(categoryId int) ([]model.Category, error) {
	query := fmt.Sprintf(`WITH RECURSIVE cat (id, name, parent_id) AS (
									SELECT id, name, parent_id FROM %s WHERE id=$1
									UNION ALL
									SELECT d.id, d.name, d.parent_id
									FROM %s d
										JOIN cat ON cat.parent_id=d.id
								)
								SELECT id, COALESCE(name, '') AS name, COALESCE(parent_id, -1) AS parent_id FROM cat`,
		categoriesTable, categoriesTable)

	var categories []model.Category
	err := r.db.Select(&categories, query, categoryId)
	return categories, err
}
