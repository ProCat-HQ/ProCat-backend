package repository

import (
	"errors"
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

	var isOkCartItems int
	queryCheckCartStock := fmt.Sprintf(`SELECT COALESCE(MIN(CAST(in_stock_number - items_number >= 0 AS INTEGER)), 0)
										       FROM %s i
												JOIN %s c ON i.item_id = c.item_id WHERE c.cart_id=$1`, itemsStoresTable, cartsItemsTable)

	err = tx.Get(&isOkCartItems, queryCheckCartStock, cartId)
	if err != nil {
		return model.OrderCheque{}, err
	}

	if isOkCartItems == 0 {
		return model.OrderCheque{}, errors.New("some items from the cart are out of stock")
	}

	var orderId int
	err = tx.Get(&orderId, orderCreationQuery, status, totalPrice, depositPrice, rpStart, rpEnd, address, lat, lon, companyName, userId)
	if err != nil {
		return model.OrderCheque{}, err
	}

	queryMoveFromCartToOrder := fmt.Sprintf(`INSERT INTO %s (items_number, order_id, item_id)
													SELECT c.items_number, $1 as order_id, i.id
													FROM %s c
													LEFT JOIN %s i ON c.item_id = i.id
													WHERE c.cart_id=$2`, ordersItemsTable, cartsItemsTable, itemsTable)

	_, err = tx.Exec(queryMoveFromCartToOrder, orderId, cartId)
	if err != nil {
		return model.OrderCheque{}, err
	}

	defaultStoreId := 1

	queryUpdateStockItems := fmt.Sprintf(`UPDATE %s i SET in_stock_number = in_stock_number - items_number
												FROM %s c WHERE store_id=$1 AND c.cart_id=$2 AND c.item_id = i.item_id`, itemsStoresTable, cartsItemsTable)

	_, err = tx.Exec(queryUpdateStockItems, defaultStoreId, cartId)
	if err != nil {
		return model.OrderCheque{}, err
	}

	queryUpdateItemBoolStock := fmt.Sprintf(`UPDATE %s it SET is_in_stock = i.in_stock_number > 0
													FROM %s i WHERE it.id = i.item_id AND store_id=$1`, itemsTable, itemsStoresTable)

	_, err = tx.Exec(queryUpdateItemBoolStock, defaultStoreId)
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
	query := fmt.Sprintf(`SELECT COALESCE(SUM(i.price * c.items_number), 0) as price,
       							COALESCE(SUM(i.price_deposit * c.items_number), 0) as price_deposit FROM %s c
								LEFT JOIN %s i ON c.item_id = i.id WHERE c.cart_id=$1`, cartsItemsTable, itemsTable)
	var res struct {
		Price        int `db:"price"`
		PriceDeposit int `db:"price_deposit"`
	}
	err := r.db.Get(&res, query, cartId)
	if err != nil {
		return 0, 0, err
	}
	if res.Price == 0 && res.PriceDeposit == 0 {
		return 0, 0, errors.New("empty cart")
	}
	return res.Price, res.PriceDeposit, nil
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
	if err != nil {
		return err
	}
	if status == model.Rejected {
		// delete delivery corresponding to the order
		queryDelete := fmt.Sprintf(`DELETE FROM %s WHERE order_id=$1`, deliveriesTable)
		_, err = r.db.Exec(queryDelete, orderId)
		if err != nil {
			return err
		}
	}
	if status == model.Returned {
		defaultStoreId := 1
		queryUpdateStock := fmt.Sprintf(`UPDATE %s i SET in_stock_number = i.in_stock_number + o.items_number
												FROM %s o WHERE o.order_id=$1 AND i.store_id=$2 AND o.item_id=i.item_id`,
			itemsStoresTable, ordersItemsTable)

		queryUpdateInStockBoolItem := fmt.Sprintf(`UPDATE %s it SET is_in_stock = i.in_stock_number > 0
														  FROM %s i WHERE it.id = i.item_id AND store_id=$1`,
			itemsTable, itemsStoresTable)

		tx, err := r.db.Beginx()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		_, err = tx.Exec(queryUpdateStock, orderId, defaultStoreId)
		if err != nil {
			return err
		}
		_, err = tx.Exec(queryUpdateInStockBoolItem, defaultStoreId)
		if err != nil {
			return err
		}
		return tx.Commit()
	}

	return nil
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
                				 WHERE id=$3 RETURNING paid, price, order_id`, paymentsTable)

	queryUpdateOrderStatus := fmt.Sprintf(`UPDATE %s SET status = $1 WHERE id = $2`, ordersTable)

	queryPaymentsCount := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE order_id=$1`, paymentsTable)

	var money struct {
		Paid    int `db:"paid"`
		Price   int `db:"price"`
		OrderId int `db:"order_id"`
	}
	err := r.db.Get(&money, query, paid, method, paymentId)
	if err != nil {
		return err
	}
	if money.Paid >= money.Price {
		var count int
		err = r.db.Get(&count, queryPaymentsCount, money.OrderId)
		if err != nil {
			return err
		}
		if count > 1 {
			_, err = r.db.Exec(queryUpdateOrderStatus, model.Extended, money.OrderId)
		} else if count == 1 {
			_, err = r.db.Exec(queryUpdateOrderStatus, model.Accepted, money.OrderId)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *OrderPostgres) ExtendOrder(orderId int, rentalPeriodEnd time.Time) error {
	query := fmt.Sprintf(`INSERT INTO %s (rental_period_end, order_id) VALUES ($1, $2)`, ordersExtensionTable)
	queryStatus := fmt.Sprintf(`UPDATE %s SET status=$1 WHERE id=$2`, ordersTable)

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, rentalPeriodEnd, orderId)
	if err != nil {
		return err
	}
	_, err = tx.Exec(queryStatus, model.ExtensionRequest, orderId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *OrderPostgres) GetRentalPeriodEndFromExtension(orderId int) (time.Time, error) {
	query := fmt.Sprintf(`SELECT rental_period_end FROM %s WHERE order_id=$1`, ordersExtensionTable)
	var rentalPeriodEnd time.Time
	err := r.db.Get(&rentalPeriodEnd, query, orderId)
	return rentalPeriodEnd, err
}

func (r *OrderPostgres) ConfirmOrderExtension(orderId int, rentalPeriodEnd time.Time, rentalPeriodDays int,
	status string, deposit bool) error {

	queryUpdateOrder := fmt.Sprintf(`UPDATE %s SET status=$1, rental_period_end=$2,
              							    total_price = total_price + $3, deposit=COALESCE(deposit, 0) + $4
          									WHERE id=$5`, ordersTable)

	queryDeleteOrdersExtension := fmt.Sprintf(`DELETE FROM %s WHERE order_id=$1`, ordersExtensionTable)

	queryGetOrderItemPriceAndDeposit := fmt.Sprintf(`SELECT SUM(i.price * o.items_number) as price,
       														SUM(i.price_deposit * o.items_number) as price_deposit
															FROM %s o
															LEFT JOIN %s i ON o.item_id = i.id WHERE o.order_id=$1`,
		ordersItemsTable, itemsTable)

	queryAddPayment := fmt.Sprintf(`INSERT INTO %s (price, order_id) VALUES ($1, $2)`, paymentsTable)

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var prices struct {
		Price        int `db:"price"`
		PriceDeposit int `db:"price_deposit"`
	}
	err = tx.Get(&prices, queryGetOrderItemPriceAndDeposit, orderId)
	if err != nil {
		return err
	}

	totalPrice := prices.Price * rentalPeriodDays
	depositPrice := func() int {
		if deposit {
			return prices.PriceDeposit
		}
		return 0
	}()

	_, err = tx.Exec(queryUpdateOrder, status, rentalPeriodEnd, totalPrice, depositPrice, orderId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(queryDeleteOrdersExtension, orderId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(queryAddPayment, totalPrice+depositPrice, orderId)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *OrderPostgres) ReturnOrder(orderId int, timeStart, timeEnd time.Time, newStatus, deliveryMethod string) error {
	queryUpdateStatus := fmt.Sprintf(`UPDATE %s SET status=$1 WHERE id=$2`, ordersTable)
	queryCreateDelivery := fmt.Sprintf(`INSERT INTO %s (time_start, time_end, method, order_id)
												VALUES ($1, $2, $3, $4)`, deliveriesTable)

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(queryUpdateStatus, newStatus, orderId)
	if err != nil {
		return err
	}
	_, err = tx.Exec(queryCreateDelivery, timeStart, timeEnd, deliveryMethod, orderId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *OrderPostgres) NeedRepairForOrder(orderId, price int) error {
	queryUpdateStatus := fmt.Sprintf(`UPDATE %s SET status=$1 WHERE id=$2`, ordersTable)
	queryCreatePayment := fmt.Sprintf(`INSERT INTO %s (price, order_id) VALUES ($1, $2)`, paymentsTable)

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(queryUpdateStatus, model.AwaitingRepairPayment, orderId)
	if err != nil {
		return err
	}
	_, err = tx.Exec(queryCreatePayment, price, orderId)
	if err != nil {
		return err
	}

	return tx.Commit()
}
