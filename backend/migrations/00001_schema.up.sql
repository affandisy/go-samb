CREATE TABLE master_supplier (
    supplier_pk SERIAL PRIMARY KEY,
    supplier_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE master_customer (
    customer_pk SERIAL PRIMARY KEY,
    customer_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE master_product (
    product_pk SERIAL PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE master_warehouse (
    whs_pk SERIAL PRIMARY KEY,
    whs_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE trx_in_header (
    trx_in_pk SERIAL PRIMARY KEY,
    trx_in_no VARCHAR(50) NOT NULL UNIQUE,
    whs_idf INTEGER REFERENCES master_warehouse(whs_pk),
    trx_in_date DATE NOT NULL,
    trx_in_supp_idf INTEGER REFERENCES master_supplier(supplier_pk),
    trx_in_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE trx_in_detail (
    trx_in_d_pk SERIAL PRIMARY KEY,
    trx_in_idf INTEGER REFERENCES trx_in_header(trx_in_pk) ON DELETE CASCADE,
    trx_in_d_product_idf INTEGER REFERENCES master_product(product_pk),
    trx_in_d_qty_dus INTEGER NOT NULL DEFAULT 0,
    trx_in_d_qty_pcs INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE trx_out_header (
    trx_out_pk SERIAL PRIMARY KEY,
    trx_out_no VARCHAR(50) NOT NULL UNIQUE,
    whs_idf INTEGER REFERENCES master_warehouse(whs_pk),
    trx_out_date DATE NOT NULL,
    trx_out_cust_idf INTEGER REFERENCES master_customer(customer_pk),
    trx_out_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE trx_out_detail (
    trx_out_d_pk SERIAL PRIMARY KEY,
    trx_out_idf INTEGER REFERENCES trx_out_header(trx_out_pk) ON DELETE CASCADE,
    trx_out_d_product_idf INTEGER REFERENCES master_product(product_pk),
    trx_out_d_qty_dus INTEGER NOT NULL DEFAULT 0,
    trx_out_d_qty_pcs INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE VIEW vw_stock_report AS
SELECT 
    w.whs_name AS warehouse,
    p.product_name AS product,
    COALESCE(SUM(tid.trx_in_d_qty_dus), 0) - COALESCE(SUM(tod.trx_out_d_qty_dus), 0) AS qty_dus,
    COALESCE(SUM(tid.trx_in_d_qty_pcs), 0) - COALESCE(SUM(tod.trx_out_d_qty_pcs), 0) AS qty_pcs
FROM master_warehouse w
CROSS JOIN master_product p
LEFT JOIN trx_in_header tih ON tih.whs_idf = w.whs_pk
LEFT JOIN trx_in_detail tid ON tid.trx_in_idf = tih.trx_in_pk AND tid.trx_in_d_product_idf = p.product_pk
LEFT JOIN trx_out_header toh ON toh.whs_idf = w.whs_pk
LEFT JOIN trx_out_detail tod ON tod.trx_out_idf = toh.trx_out_pk AND tod.trx_out_d_product_idf = p.product_pk
GROUP BY w.whs_name, p.product_name
ORDER BY w.whs_name, p.product_name;

INSERT INTO master_warehouse (whs_name) VALUES ('Gudang A');

INSERT INTO master_product (product_name) VALUES ('Produk A'), ('Produk B'), ('Produk C');

INSERT INTO master_supplier (supplier_name) VALUES ('Supplier 1'), ('Supplier 2'), ('Supplier 3');
    
INSERT INTO master_customer (customer_name) VALUES ('Customer 1'), ('Customer 2'), ('Customer 3');
