-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE SCHEMA IF NOT EXISTS product_management;

CREATE TABLE IF NOT EXISTS product_management.products( 
  id UUID PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	sku  VARCHAR(255) NOT NULL UNIQUE,
  barcode VARCHAR(255) NOT NULL UNIQUE,
	description VARCHAR(255),
	category  VARCHAR(255),
	brand  VARCHAR(255),
	unit_type VARCHAR(255),
	is_active BOOLEAN,
	created_at TIMESTAMPTZ,
	updated_at TIMESTAMPTZ
);

-- CREATE INDEX idx_products_sku on product(sku);
CREATE INDEX idx_products_barcode ON product_management.products(barcode);
-- CREATE INDEX idx_products_is_active ON products(is_active);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP SCHEMA IF EXISTS product_management CASCADE;
-- +goose StatementEnd
