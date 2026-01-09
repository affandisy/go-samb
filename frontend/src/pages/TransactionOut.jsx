// src/pages/TransactionOut.jsx
import { useEffect, useState } from 'react';
import { Container, Table, Button, Form, Row, Col, Card } from 'react-bootstrap';
import api from '../utils/api';

const TransactionOut = () => {
  const [transactions, setTransactions] = useState([]);
  const [showForm, setShowForm] = useState(false);

  const [errorMsg, setErrorMsg] = useState('');
  const [successMsg, setSuccessMsg] = useState('');
  
  const [warehouses, setWarehouses] = useState([]);
  const [customers, setCustomers] = useState([]);
  const [products, setProducts] = useState([]);

  const [formData, setFormData] = useState({
    trx_out_no: '', whs_idf: '', trx_out_date: new Date().toISOString().split('T')[0],
    trx_out_cust_idf: '', trx_out_notes: '', details: []
  });

  useEffect(() => { fetchTransactions(); fetchMasterData(); }, []);

  const fetchTransactions = async () => {
    try { const res = await api.get('/trx-out'); setTransactions(res.data || []); } catch (e) {}
  };

  const fetchMasterData = async () => {
    const [w, c, p] = await Promise.all([api.get('/warehouses'), api.get('/customers'), api.get('/products')]);
    setWarehouses(w.data || []); setCustomers(c.data || []); setProducts(p.data || []);
  };

  const handleAddDetail = () => {
    setFormData({ ...formData, details: [...formData.details, { trx_out_d_product_idf: '', trx_out_d_qty_dus: 0, trx_out_d_qty_pcs: 0 }] });
  };

  const handleDetailChange = (index, field, value) => {
    const newDetails = [...formData.details];
    
    if (field === 'trx_out_d_product_idf') {
      newDetails[index][field] = parseInt(value) || '';
    } else {
      let val = parseInt(value) || 0;
      if (val < 0) val = 0;
      newDetails[index][field] = val;
    }
    
    setFormData({ ...formData, details: newDetails });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setErrorMsg('');
    setSuccessMsg('');

    // Validasi details
    if (formData.details.length === 0) {
      setErrorMsg('Tambahkan minimal 1 item produk!');
      window.scrollTo(0, 0);
      return;
    }

    // Validasi setiap detail item
    for (let i = 0; i < formData.details.length; i++) {
      const detail = formData.details[i];
      if (!detail.trx_out_d_product_idf) {
        setErrorMsg(`Item ${i + 1}: Pilih produk terlebih dahulu!`);
        window.scrollTo(0, 0);
        return;
      }
    }

    try {
      const payload = {
        ...formData,
        whs_idf: parseInt(formData.whs_idf),
        trx_out_cust_idf: parseInt(formData.trx_out_cust_idf),
      };
      await api.post('/trx-out', payload);
      setSuccessMsg('Transaksi berhasil!');
      setTimeout(() => {
        setShowForm(false);
        fetchTransactions();
      }, 1000);
    } catch (err) { 
      const message = err.response?.data?.error || "Terjadi kesalahan pada server";
      setErrorMsg(message);
      window.scrollTo(0, 0);
    }
  };

  return (
    <Container>
      <div className="d-flex justify-content-between mb-3">
        <h2>Transaksi Barang Keluar</h2>
        <Button onClick={() => setShowForm(!showForm)} variant={showForm ? "secondary" : "danger"}>
          {showForm ? "Kembali" : "Tambah Transaksi"}
        </Button>
      </div>

      {showForm ? (
        <Card className="p-4">
          {errorMsg && <div className="alert alert-danger">{errorMsg}</div>}
          {successMsg && <div className="alert alert-success">{successMsg}</div>}
          <Form onSubmit={handleSubmit}>
            <Row className="mb-3">
              <Col><Form.Control placeholder="No Trx" required onChange={e => setFormData({...formData, trx_out_no: e.target.value})} /></Col>
              <Col><Form.Control type="date" value={formData.trx_out_date} onChange={e => setFormData({...formData, trx_out_date: e.target.value})} /></Col>
            </Row>
            <Row className="mb-3">
              <Col>
                <Form.Select required onChange={e => setFormData({...formData, whs_idf: e.target.value})}>
                  <option value="">Pilih Gudang</option>
                  {warehouses.map(w => <option key={w.whs_pk} value={w.whs_pk}>{w.whs_name}</option>)}
                </Form.Select>
              </Col>
              <Col>
                <Form.Select required onChange={e => setFormData({...formData, trx_out_cust_idf: e.target.value})}>
                  <option value="">Pilih Customer</option>
                  {customers.map(c => <option key={c.customer_pk} value={c.customer_pk}>{c.customer_name}</option>)}
                </Form.Select>
              </Col>
            </Row>
            <Form.Control as="textarea" placeholder="Catatan" className="mb-3" onChange={e => setFormData({...formData, trx_out_notes: e.target.value})} />
            
            <h5>Detail <Button size="sm" onClick={handleAddDetail}>+</Button></h5>
            {formData.details.map((d, idx) => (
              <Row key={idx} className="mb-2">
                <Col md={6}><Form.Select required onChange={e => handleDetailChange(idx, 'trx_out_d_product_idf', e.target.value)}><option value="">Pilih Produk</option>{products.map(p => <option key={p.product_pk} value={p.product_pk}>{p.product_name}</option>)}</Form.Select></Col>
                <Col><Form.Control type="number" placeholder="Dus" onChange={e => handleDetailChange(idx, 'trx_out_d_qty_dus', e.target.value)} /></Col>
                <Col><Form.Control type="number" placeholder="Pcs" onChange={e => handleDetailChange(idx, 'trx_out_d_qty_pcs', e.target.value)} /></Col>
              </Row>
            ))}
            <Button type="submit" className="mt-3" variant="danger">Simpan Transaksi Out</Button>
          </Form>
        </Card>
      ) : (
        <Table striped bordered hover>
          <thead><tr><th>No</th><th>Tanggal</th><th>Gudang</th><th>Customer</th><th>Notes</th></tr></thead>
          <tbody>
            {transactions.map(t => (
              <tr key={t.header.trx_out_pk}>
                <td>{t.header.trx_out_no}</td>
                <td>{new Date(t.header.trx_out_date).toLocaleDateString('id-ID')}</td>
                <td>{t.warehouse_name}</td>
                <td>{t.customer_name}</td>
                <td>{t.header.trx_out_notes}</td>
              </tr>
            ))}
          </tbody>
        </Table>
      )}
    </Container>
  );
};

export default TransactionOut;