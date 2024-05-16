package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"strconv"
	"time"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) GetOrder(orderId int) (model.Order, error) {
	queryOrder := fmt.Sprintf(`SELECT id, status, total_price, COALESCE(deposit, 0) AS deposit,
       								   rental_period_start, rental_period_end, address, COALESCE(latitude, '') AS latitude,
       								   COALESCE(longitude, '') AS longitude, COALESCE(company_name, '') AS company_name,
       								   created_at, user_id FROM %s
       								   WHERE id=$1`, ordersTable)
	queryItems := fmt.Sprintf(`SELECT o.item_id, i.name, i.price, i.price_deposit, o.items_number as count, COALESCE(im.image, '') AS image
									  FROM %s o
									  LEFT JOIN %s i on i.id = o.item_id
									  LEFT JOIN (SELECT DISTINCT ON (item_id) * FROM %s) im ON i.id = im.item_id
									  WHERE o.order_id=$1`, ordersItemsTable, itemsTable, itemsImagesTable)

	tx, err := r.db.Beginx()
	if err != nil {
		return model.Order{}, err
	}
	defer tx.Rollback()

	var order model.Order
	err = tx.Get(&order, queryOrder, orderId)
	if err != nil {
		return model.Order{}, err
	}
	var items []model.OrderSmallItem
	err = tx.Select(&items, queryItems, orderId)
	if err != nil {
		return model.Order{}, err
	}
	order.Items = items

	err = tx.Commit()
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}

func (r *OrderPostgres) GetAllOrders(limit, offset, userId int, statuses []string) (int, []model.Order, error) {
	queryCount := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, ordersTable)

	queryOrders := fmt.Sprintf(`SELECT id, status, total_price, COALESCE(deposit, 0) AS deposit,
       								   rental_period_start, rental_period_end, address, COALESCE(latitude, '') AS latitude,
       								   COALESCE(longitude, '') AS longitude,
       								   COALESCE(company_name, '') AS company_name, created_at, user_id FROM %s`, ordersTable)

	argCounter := 1
	args := make([]interface{}, 0)
	if userId != 0 {
		queryOrders += ` WHERE user_id = $` + strconv.Itoa(argCounter)
		queryCount += ` WHERE user_id = $` + strconv.Itoa(argCounter)
		argCounter += 1
		args = append(args, userId)
	}
	if statuses != nil {
		if argCounter == 1 {
			queryOrders += ` WHERE `
			queryCount += ` WHERE `
		} else {
			queryOrders += ` AND `
			queryCount += ` AND `
		}
		queryOrders += ` status IN (`
		queryCount += ` status IN (`
		for i := range statuses {
			queryOrders += fmt.Sprintf("$%d", argCounter)
			queryCount += fmt.Sprintf("$%d", argCounter)
			if i != len(statuses)-1 {
				queryOrders += ", "
				queryCount += ", "
			}
			args = append(args, statuses[i])
			argCounter += 1
		}
		queryOrders += `)`
		queryCount += `)`
	}
	queryOrders += ` LIMIT $` + strconv.Itoa(argCounter) + ` OFFSET $` + strconv.Itoa(argCounter+1)
	args = append(args, limit, offset)

	queryItems := fmt.Sprintf(`SELECT o.item_id, i.name, i.price, i.price_deposit, o.items_number as count, COALESCE(im.image, '') AS image
									  FROM %s o
									  LEFT JOIN %s i on i.id = o.item_id
									  LEFT JOIN (SELECT DISTINCT ON (item_id) * FROM %s) im ON i.id = im.item_id
									  WHERE o.order_id=$1`, ordersItemsTable, itemsTable, itemsImagesTable)

	tx, err := r.db.Beginx()
	if err != nil {
		return 0, nil, err
	}
	defer tx.Rollback()

	var orders []model.Order
	var count int

	err = tx.Get(&count, queryCount, args[:len(args)-2]...)
	if err != nil {
		return 0, nil, err
	}
	err = tx.Select(&orders, queryOrders, args...)
	if err != nil {
		return 0, nil, err
	}
	for i := range orders {
		var items []model.OrderSmallItem
		err = tx.Select(&items, queryItems, orders[i].Id)
		if err != nil {
			return 0, nil, err
		}
		orders[i].Items = items
	}

	err = tx.Commit()
	if err != nil {
		return 0, nil, err
	}

	return count, orders, nil
}

func (r *OrderPostgres) GetUserById(userId int) (model.User, error) {
	query := fmt.Sprintf(`SELECT id, fullname, COALESCE(email, '') AS email, phone_number,
       							COALESCE(identification_number, '') AS identification_number, is_confirmed, role, created_at
								 FROM %s WHERE id=$1`, usersTable)

	var user model.User
	err := r.db.Get(&user, query, userId)

	return user, err
}

func (r *OrderPostgres) CreateOrder(status string, deposit bool, rpStart, rpEnd time.Time,
	address string, lat, lon float64, companyName string, userId int,
	deliveryMethod string, tStart, tEnd time.Time, rentPeriodDays int) (model.OrderCheque, error) {

	cartId, err := r.GetUsersCartId(userId)
	if err != nil {
		return model.OrderCheque{}, err
	}
	totalPrice, totalDeposit, err := r.GetTotalCartPrices(cartId)
	if err != nil {
		return model.OrderCheque{}, err
	}
	totalPrice *= rentPeriodDays

	itemsCheque, err := r.GetItemCheque(cartId)
	if err != nil {
		return model.OrderCheque{}, err
	}

	orderCreationQuery := fmt.Sprintf(`INSERT INTO %s (status, total_price, deposit, rental_period_start, rental_period_end,
                    							address, latitude, longitude, company_name, user_id)
												VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`, ordersTable)
	depositPrice := func() int {
		if deposit {
			return totalDeposit
		}
		return 0
	}()

	tx, err := r.db.Beginx()
	if err != nil {
		return model.OrderCheque{}, err
	}
	defer tx.Rollback()

	var orderId int
	err = tx.Get(&orderId, orderCreationQuery, status, totalPrice, depositPrice, rpStart, rpEnd, address, lat, lon, companyName, userId)
	if err != nil {
		return model.OrderCheque{}, err
	}

	// if it won't work change $1 to order_id using sprintf
	queryMoveFromCartToOrder := fmt.Sprintf(`INSERT INTO %s (items_number, order_id, item_id)
													SELECT c.items_number, $1 as order_id, i.id
													FROM %s c
													LEFT JOIN %s i ON c.item_id = i.id
													WHERE c.cart_id=$2`, ordersItemsTable, cartsItemsTable, itemsTable)

	_, err = tx.Exec(queryMoveFromCartToOrder, orderId, cartId)
	if err != nil {
		return model.OrderCheque{}, err
	}

	queryDeleteFromCart := fmt.Sprintf(`DELETE FROM %s WHERE cart_id=$1`, cartsItemsTable)
	_, err = tx.Exec(queryDeleteFromCart, cartId)
	if err != nil {
		return model.OrderCheque{}, err
	}

	queryAddDelivery := fmt.Sprintf(`INSERT INTO %s (time_start, time_end, method, order_id) 
											VALUES ($1, $2, $3, $4)`, deliveriesTable)
	_, err = tx.Exec(queryAddDelivery, tStart, tEnd, deliveryMethod, orderId)
	if err != nil {
		return model.OrderCheque{}, err
	}

	queryAddPayment := fmt.Sprintf(`INSERT INTO %s (price, order_id) VALUES ($1, $2)`, paymentsTable)
	_, err = tx.Exec(queryAddPayment, totalPrice+depositPrice, orderId)

	err = tx.Commit()
	if err != nil {
		return model.OrderCheque{}, err
	}

	orderCheque := model.OrderCheque{
		OrderId:    orderId,
		TotalPrice: totalPrice,
		Deposit:    depositPrice,
		Items:      itemsCheque,
	}

	return orderCheque, nil
}

func (r *OrderPostgres) GetTotalCartPrices(cartId int) (int, int, error) {
	query := fmt.Sprintf(`SELECT SUM(i.price * c.items_number) as price, SUM(i.price_deposit * c.items_number) as price_deposit FROM %s c
								LEFT JOIN %s i ON c.item_id = i.id WHERE c.cart_id=$1`, cartsItemsTable, itemsTable)
	var res struct {
		Price        int `db:"price"`
		PriceDeposit int `db:"price_deposit"`
	}
	err := r.db.Get(&res, query, cartId)
	return res.Price, res.PriceDeposit, err
}

func (r *OrderPostgres) GetUsersCartId(userId int) (int, error) {
	query := fmt.Sprintf(`SELECT id FROM %s WHERE user_id = $1`, cartsTable)

	var id int
	err := r.db.Get(&id, query, userId)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *OrderPostgres) GetItemCheque(cartId int) ([]model.ItemCheque, error) {
	query := fmt.Sprintf(`SELECT i.name, c.items_number as count, i.price, i.price_deposit FROM %s c
								LEFT JOIN %s i ON c.item_id = i.id WHERE c.cart_id=$1`, cartsItemsTable, itemsTable)
	var itemCheque []model.ItemCheque
	err := r.db.Select(&itemCheque, query, cartId)
	return itemCheque, err
}

func (r *OrderPostgres) ChangeOrderStatus(orderId int, status string) error {
	query := fmt.Sprintf(`UPDATE %s SET status=$1 WHERE id=$2`, ordersTable)
	_, err := r.db.Exec(query, status, orderId)
	return err
}

func (r *OrderPostgres) GetPaymentsForOrder(orderId int) ([]model.Payment, error) {
	query := fmt.Sprintf(`SELECT id, paid, COALESCE(method, '') AS method, price, created_at
								 FROM %s WHERE order_id=$1`, paymentsTable)
	var payments []model.Payment
	err := r.db.Select(&payments, query, orderId)
	return payments, err
}

func (r *OrderPostgres) ChangePaymentStatus(paymentId, paid int, method string) error {
	query := fmt.Sprintf(`UPDATE %s SET paid = paid + $1,
                    			 method = COALESCE(method, '') || $2 || ';'
                				 WHERE id=$3`, paymentsTable)

	_, err := r.db.Exec(query, paid, method, paymentId)
	return err
}
