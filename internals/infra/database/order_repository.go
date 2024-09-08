package database

import (
	"database/sql"
	"github.com/guirialli/go-pos/clean_arch/internals/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	smt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = smt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) FindAll() ([]entity.Order, error) {
	smt, err := r.Db.Prepare("SELECT * FROM orders")
	if err != nil {
		return make([]entity.Order, 0), err
	}
	defer func(smt *sql.Stmt) {
		_ = smt.Close()
	}(smt)

	rows, err := smt.Query()
	if err != nil {
		return make([]entity.Order, 0), err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		err := rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice)
		if err != nil {
			return make([]entity.Order, 0), err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.Db.QueryRow("Select count(*) from orders").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
