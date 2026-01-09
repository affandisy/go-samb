// src/pages/StockReport.jsx
import { useEffect, useState } from 'react';
import { Container, Table } from 'react-bootstrap';
import api from '../utils/api';

const StockReport = () => {
  const [stocks, setStocks] = useState([]);

  useEffect(() => {
    api.get('/stock-report')
      .then(res => setStocks(res.data || []))
      .catch(err => console.error(err));
  }, []);

  return (
    <Container>
      <h2 className="mb-4">Laporan Stok Terkini</h2>
      <Table bordered hover responsive>
        <thead className="table-dark">
          <tr>
            <th>Gudang</th>
            <th>Produk</th>
            <th className="text-end">Qty Dus</th>
            <th className="text-end">Qty Pcs</th>
          </tr>
        </thead>
        <tbody>
          {stocks.length > 0 ? stocks.map((s, idx) => (
            <tr key={idx}>
              <td>{s.warehouse}</td>
              <td>{s.product}</td>
              <td className="text-end">{s.qty_dus}</td>
              <td className="text-end">{s.qty_pcs}</td>
            </tr>
          )) : <tr><td colSpan="4" className="text-center">Data kosong</td></tr>}
        </tbody>
      </Table>
    </Container>
  );
};

export default StockReport;