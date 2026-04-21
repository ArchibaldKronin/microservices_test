-- +goose Up
CREATE table orders (
  order_id uuid PRIMARY KEY,
  user_id uuid NOT NULL,
  part_ids uuid[] NOT NULL,
  total_price decimal NOT NULL,
  transaction_id uuid,
  payment_method text CHECK (payment_method IN ('UNKNOWN', 'CARD', 'SBP', 'CREDIT_CARD', 'INVESTOR_MONEY')),
  status text NOT NULL CHECK (status IN ('PENDING_PAYMENT', 'PAID', 'CANCELLED'))
);

-- +goose Down
drop table orders