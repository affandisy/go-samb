CREATE INDEX idx_trx_in_header_whs ON trx_in_header(whs_idf);
CREATE INDEX idx_trx_in_header_date ON trx_in_header(trx_in_date);
CREATE INDEX idx_trx_in_detail_header ON trx_in_detail(trx_in_idf);
CREATE INDEX idx_trx_in_detail_product ON trx_in_detail(trx_in_d_product_idf);
CREATE INDEX idx_trx_in_detail_whs_product ON trx_in_detail(trx_in_d_product_idf, trx_in_idf);

CREATE INDEX idx_trx_out_header_whs ON trx_out_header(whs_idf);
CREATE INDEX idx_trx_out_header_date ON trx_out_header(trx_out_date);
CREATE INDEX idx_trx_out_detail_header ON trx_out_detail(trx_out_idf);
CREATE INDEX idx_trx_out_detail_product ON trx_out_detail(trx_out_d_product_idf);
CREATE INDEX idx_trx_out_detail_whs_product ON trx_out_detail(trx_out_d_product_idf, trx_out_idf);

CREATE INDEX idx_trx_in_whs_product ON trx_in_detail(trx_in_d_product_idf) 
    INCLUDE (trx_in_d_qty_dus, trx_in_d_qty_pcs);
    
CREATE INDEX idx_trx_out_whs_product ON trx_out_detail(trx_out_d_product_idf) 
    INCLUDE (trx_out_d_qty_dus, trx_out_d_qty_pcs);
