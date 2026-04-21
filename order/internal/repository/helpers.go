package repository

var RepoFields []string = []string{
	"order_id",
	"user_id",
	"part_ids",
	"total_price",
	"transaction_id",
	"payment_method",
	"status",
}

const TABLE_NAME = "orders"
