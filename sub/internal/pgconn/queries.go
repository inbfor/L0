package pg

import (
	"context"
	"fmt"
	"test/internal/model"
)

const (
	queryOrder = `
	INSERT INTO "order" 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	queryDelivery = `
	INSERT INTO "delivery"
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	queryPayment = `
	INSERT INTO "payment"
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	queryItem = `
	INSERT INTO "item"
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) `

	queryGetSingleOrder = `
	SELECT 
		o.order_uid,
		o.track_number,
		o.entry,
		d.name,
		d.phone,
		d.zip,
		d.city,
		d.address,
		d.region,
		d.email,
		p.transaction,
		p.request_id,
		p.currency,
		p.provider,
		p.amount,
		p.payment_dt,
		p.bank,
		p.delivery_cost,
		p.goods_total,
		p.custom_fee,
		o.locale,
		o.internal_signature,
		o.customer_id,
		o.delivery_service, 
		o.shardkey, 
		o.sm_id, 
		o.date_created, 
		o.oof_shard
	FROM 
		"order" AS o 
		Inner Join "delivery" AS d ON o.order_uid = d.order_id
		Inner Join "payment" AS p ON d.order_id = p.order_id
		Right Join "item" AS i ON p.order_id = i.order_id
	WHERE o.order_uid = $1
	`

	queryGetItem = `
	SELECT 
			chrt_id,
			track_number,
			price,
			rid,
			name,
			sale,
			size,
			total_price,
			nm_id,
			brand,
			status
		FROM item AS i
		WHERE order_id = $1
	`

	queryGetOrders = `
	SELECT 
	o.order_uid,
	o.track_number,
	o.entry,
	d.name,
	d.phone,
	d.zip,
	d.city,
	d.address,
	d.region,
	d.email,
	p.transaction,
	p.request_id,
	p.currency,
	p.provider,
	p.amount,
	p.payment_dt,
	p.bank,
	p.delivery_cost,
	p.goods_total,
	p.custom_fee,
	o.locale,
	o.internal_signature,
	o.customer_id,
	o.delivery_service, 
	o.shardkey, 
	o.sm_id, 
	o.date_created, 
	o.oof_shard
FROM 
	order AS o 
	Inner Join delivery AS d ON o.order_uid = d.order_id
	Inner Join payment AS p ON d.order_id = p.order_id
	`
)

func (db *DB) InsertMessage(order model.OrderMessage) error {
	_, err := db.db.Exec(context.Background(),
		queryOrder,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = db.db.Exec(context.Background(),
		queryDelivery,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = db.db.Exec(context.Background(),
		queryPayment,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, item := range order.Items {
		_, err = db.db.Exec(context.Background(),
			queryItem,
			order.OrderUID,
			item.ChrtID,
			order.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)

		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (db *DB) GetSingleMessage(uuid string) (model.OrderMessage, error) {
	var order model.OrderMessage

	err := db.db.QueryRow(context.Background(), queryGetSingleOrder, uuid).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDt,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)

	if err != nil {
		fmt.Println(err)
		return model.OrderMessage{}, err
	}

	itemRows, err := db.db.Query(context.Background(), queryGetItem, order.OrderUID)

	if err != nil {
		fmt.Println(err)
		return model.OrderMessage{}, err
	}

	defer itemRows.Close()

	items := make([]model.Item, 0)

	for itemRows.Next() {
		var item model.Item

		err = itemRows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)

		if err != nil {
			fmt.Println(err)
			return model.OrderMessage{}, err
		}
		items = append(items, item)
	}

	order.Items = items

	return order, nil
}

func (db *DB) GetAllMessages() ([]model.OrderMessage, error) {
	rows, err := db.db.Query(context.Background(), queryGetOrders)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := make([]model.OrderMessage, 0)

	for rows.Next() {
		var order model.OrderMessage

		err = rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Delivery.Name,
			&order.Delivery.Phone,
			&order.Delivery.Zip,
			&order.Delivery.City,
			&order.Delivery.Address,
			&order.Delivery.Region,
			&order.Delivery.Email,
			&order.Payment.Transaction,
			&order.Payment.RequestID,
			&order.Payment.Currency,
			&order.Payment.Provider,
			&order.Payment.Amount,
			&order.Payment.PaymentDt,
			&order.Payment.Bank,
			&order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal,
			&order.Payment.CustomFee,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
		)
		if err != nil {
			return nil, err
		}

		itemRows, err := db.db.Query(context.Background(), queryGetItem, order.OrderUID)

		if err != nil {
			return nil, err
		}

		defer itemRows.Close()

		items := make([]model.Item, 0)

		for itemRows.Next() {
			var item model.Item

			err = itemRows.Scan(
				&item.ChrtID,
				&item.TrackNumber,
				&item.Price,
				&item.Rid,
				&item.Name,
				&item.Sale,
				&item.Size,
				&item.TotalPrice,
				&item.NmID,
				&item.Brand,
				&item.Status,
			)

			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
		order.Items = items

		orders = append(orders, order)
	}
	return orders, nil
}
