DROP VIEW IF EXISTS vw_stock_report;

CREATE MATERIALIZED VIEW vw_stock_report AS
WITH stock_in AS (
    SELECT 
        tih.whs_idf,
        tid.trx_in_d_product_idf as product_id,
        SUM(tid.trx_in_d_qty_dus) as total_dus,
        SUM(tid.trx_in_d_qty_pcs) as total_pcs
    FROM trx_in_detail tid
    JOIN trx_in_header tih ON tid.trx_in_idf = tih.trx_in_pk
    GROUP BY tih.whs_idf, tid.trx_in_d_product_idf
),
stock_out AS (
    SELECT 
        toh.whs_idf,
        tod.trx_out_d_product_idf as product_id,
        SUM(tod.trx_out_d_qty_dus) as total_dus,
        SUM(tod.trx_out_d_qty_pcs) as total_pcs
    FROM trx_out_detail tod
    JOIN trx_out_header toh ON tod.trx_out_idf = toh.trx_out_pk
    GROUP BY toh.whs_idf, tod.trx_out_d_product_idf
)
SELECT 
    w.whs_name AS warehouse,
    p.product_name AS product,
    COALESCE(si.total_dus, 0) - COALESCE(so.total_dus, 0) AS qty_dus,
    COALESCE(si.total_pcs, 0) - COALESCE(so.total_pcs, 0) AS qty_pcs
FROM master_warehouse w
CROSS JOIN master_product p
LEFT JOIN stock_in si ON si.whs_idf = w.whs_pk AND si.product_id = p.product_pk
LEFT JOIN stock_out so ON so.whs_idf = w.whs_pk AND so.product_id = p.product_pk
WHERE (COALESCE(si.total_dus, 0) - COALESCE(so.total_dus, 0)) > 0 
   OR (COALESCE(si.total_pcs, 0) - COALESCE(so.total_pcs, 0)) > 0
ORDER BY w.whs_name, p.product_name;

CREATE INDEX idx_mv_stock_warehouse ON vw_stock_report(warehouse);
CREATE INDEX idx_mv_stock_product ON vw_stock_report(product);

CREATE OR REPLACE FUNCTION refresh_stock_report()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY vw_stock_report;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION trigger_refresh_stock()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM refresh_stock_report();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_refresh_stock_on_trx_in ON trx_in_header;
CREATE TRIGGER trg_refresh_stock_on_trx_in
AFTER INSERT OR UPDATE OR DELETE ON trx_in_header
FOR EACH STATEMENT
EXECUTE FUNCTION trigger_refresh_stock();

DROP TRIGGER IF EXISTS trg_refresh_stock_on_trx_out ON trx_out_header;
CREATE TRIGGER trg_refresh_stock_on_trx_out
AFTER INSERT OR UPDATE OR DELETE ON trx_out_header
FOR EACH STATEMENT
EXECUTE FUNCTION trigger_refresh_stock();