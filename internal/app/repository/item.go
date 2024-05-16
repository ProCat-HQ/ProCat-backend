package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"strconv"
	"strings"
)

type ItemPostgres struct {
	db *sqlx.DB
}

func NewItemPostgres(db *sqlx.DB) *ItemPostgres {
	return &ItemPostgres{db: db}
}

func (r *ItemPostgres) GetCategoryChildren(categoryId int) ([]int, error) {
	query := fmt.Sprintf(`WITH RECURSIVE cat (id, id_parent) AS (
				SELECT id, parent_id FROM %s WHERE id=$1
				UNION ALL
				SELECT d.id, d.parent_id
				FROM %s d
				JOIN cat ON cat.id=d.parent_id
			)
			SELECT id FROM cat`, categoriesTable, categoriesTable)
	var ids []int
	err := r.db.Select(&ids, query, categoryId)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// GetAllItems TODO someday: it returns overall items count, but have to return count of filtered items
func (r *ItemPostgres) GetAllItems(limit, offset, categoryId int, stock bool, search string) (int, []model.PieceOfItem, error) {
	queryCount := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, itemsTable)
	var count int
	err := r.db.Get(&count, queryCount)
	if err != nil {
		return 0, nil, err
	}

	query := fmt.Sprintf(`SELECT i.id, i.name, LEFT(COALESCE(i.description, ''), 30) AS description, i.price, i.is_in_stock,
	  							COALESCE(i.category_id, -1) AS category_id, COALESCE(c.name, '') AS category_name,
	 							COALESCE(im.image, '') AS image
								FROM %s i
								LEFT JOIN %s c ON i.category_id = c.id
								LEFT JOIN (SELECT DISTINCT ON (item_id) * FROM %s) im ON i.id = im.item_id
								`, itemsTable, categoriesTable, itemsImagesTable)

	argCounter := 0
	if categoryId != 0 {
		categoryIdsInt, err := r.GetCategoryChildren(categoryId)
		if err != nil {
			return 0, nil, err
		}
		categoryIdsStrings := make([]string, 0)
		for _, v := range categoryIdsInt {
			categoryIdsStrings = append(categoryIdsStrings, strconv.Itoa(v))
		}

		query += " WHERE category_id IN " + "(" + strings.Join(categoryIdsStrings, ", ") + ")"
		argCounter += 1
	}
	if stock {
		if argCounter == 0 {
			query += " WHERE "
			argCounter += 1
		} else {
			query += " AND "
		}
		query += `i.is_in_stock = true`
	}
	if search != "" {
		if argCounter == 0 {
			query += " WHERE "
			argCounter += 1
		} else {
			query += " AND "
		}
		query += "i.name LIKE '%' || " + "$" + strconv.Itoa(argCounter) + " || '%'"
		argCounter += 1
		query += " OR i.description LIKE '%' || " + "$" + strconv.Itoa(argCounter) + " || '%'"
		argCounter += 1
	}

	if argCounter == 0 {
		argCounter += 1
	}

	query += ` OFFSET $` + strconv.Itoa(argCounter)
	argCounter += 1
	query += ` LIMIT $` + strconv.Itoa(argCounter)

	var items []model.PieceOfItem

	if argCounter == 2 {
		err = r.db.Select(&items, query, offset, limit)
	} else {
		err = r.db.Select(&items, query, search, search, offset, limit)
	}
	if err != nil {
		return 0, nil, err
	}

	return count, items, nil
}

func (r *ItemPostgres) GetItem(itemId int) (model.Item, error) {
	queryItem := fmt.Sprintf(`SELECT i.id, i.name, COALESCE(i.description, '') AS description, i.price, i.is_in_stock,
       								  COALESCE(i.category_id, -1) AS category_id, COALESCE(c.name, '') AS category_name
									  FROM %s i LEFT JOIN %s c ON i.category_id = c.id
									  WHERE i.id = $1`, itemsTable, categoriesTable)

	queryInfo := fmt.Sprintf(`SELECT id, name, description FROM %s WHERE item_id = $1`, infosTable)

	queryImage := fmt.Sprintf(`SELECT id, image FROM %s WHERE item_id=$1`, itemsImagesTable)

	queryStores := fmt.Sprintf(`SELECT s.id, i.in_stock_number, s.name, s.address, s.working_hours_start, s.working_hours_end
									  FROM %s s JOIN %s i ON s.id = i.store_id WHERE i.item_id=$1`, storesTable, itemsStoresTable)

	var item model.Item
	var infos []model.Info
	var images []model.ItemImage
	var stores []model.ItemStores

	tx, err := r.db.Beginx()
	if err != nil {
		return item, err
	}
	defer tx.Rollback()

	err = tx.Get(&item, queryItem, itemId)
	if err != nil {
		return item, err
	}

	err = tx.Select(&infos, queryInfo, itemId)
	if err != nil {
		return item, err
	}

	err = tx.Select(&images, queryImage, itemId)
	if err != nil {
		return item, err
	}

	err = tx.Select(&stores, queryStores, itemId)
	if err != nil {
		return item, err
	}

	err = tx.Commit()
	if err != nil {
		return item, nil
	}

	item.Info = infos
	item.Images = images
	item.ItemStores = stores
	return item, nil
}

func (r *ItemPostgres) CreateItem(name, description string, price, priceDeposit, categoryId int) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (name, description, price, price_deposit, category_id)
								VALUES ($1, CASE WHEN LENGTH($2)=0 THEN NULL ELSE $3 END, $4, $5, $6) RETURNING id`, itemsTable)

	var id int
	err := r.db.Get(&id, query, name, description, description, price, priceDeposit, categoryId)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ItemPostgres) SaveFilenames(itemId int, filenames []string) error {
	query := fmt.Sprintf(`INSERT INTO %s (image, item_id) VALUES `, itemsImagesTable)

	for i := 0; i < len(filenames)-1; i++ {
		query += "('" + filenames[i] + "', '" + strconv.Itoa(itemId) + "'), "
	}
	query += "('" + filenames[len(filenames)-1] + "', '" + strconv.Itoa(itemId) + "')"

	_, err := r.db.Exec(query)
	return err
}

func (r *ItemPostgres) DeleteItem(itemId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, itemsTable)
	_, err := r.db.Exec(query, itemId)
	return err
}

// ChangeItem TODO: SQL Injections sensitive
func (r *ItemPostgres) ChangeItem(itemId int, name, description, price, priceDeposit, categoryId *string) error {
	query := fmt.Sprintf(`UPDATE %s SET `, itemsTable)
	argCounter := 0
	if name != nil {
		query += fmt.Sprintf(` name = '%s'`, *name)
		argCounter += 1
	}
	if description != nil {
		if argCounter != 0 {
			query += `, `
		}
		query += fmt.Sprintf(` description = '%s'`, *description)
		argCounter += 1
	}
	if price != nil {
		if argCounter != 0 {
			query += `, `
		}
		query += fmt.Sprintf(` price = %s`, *price)
		argCounter += 1
	}
	if priceDeposit != nil {
		if argCounter != 0 {
			query += `, `
		}
		query += fmt.Sprintf(` price_deposit = %s`, *priceDeposit)
		argCounter += 1
	}
	if categoryId != nil {
		if argCounter != 0 {
			query += `, `
		}
		query += fmt.Sprintf(` category_id = %s`, *categoryId)
		argCounter += 1
	}
	query += ` WHERE id = $1 `

	_, err := r.db.Exec(query, itemId)
	return err
}
