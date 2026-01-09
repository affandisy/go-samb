// src/pages/TransactionIn.jsx
import { useEffect, useState } from 'react';
import { Container, Table, Button, Form, Row, Col, Card } from 'react-bootstrap';
import api from '../utils/api';

const TransactionIn = () => {
  const [transactions, setTransactions] = useState([]);
  const [showForm, setShowForm] = useState(false);

  const [errorMsg, setErrorMsg] = useState('');
  const [successMsg, setSuccessMsg] = useState('');
  
  // State untuk Data Select (Dropdown)
  const [warehouses, setWarehouses] = useState([]);
  const [suppliers, setSuppliers] = useState([]);
  const [products, setProducts] = useState([]);

  // Form State
  const [formData, setFormData] = useState({
    trx_in_no: '',
    whs_idf: '',
    trx_in_date: new Date().toISOString().split('T')[0],
    trx_in_supp_idf: '',
    trx_in_notes: '',
    details: []
  });

  useEffect(() => {
    fetchTransactions();
    fetchMasterData();
  }, []);

  const fetchTransactions = async () => {
    try {
      const res = await api.get('/trx-in');
      setTransactions(res.data || []);
    } catch (error) { console.error(error); }
  };

  const fetchMasterData = async () => {
    const [w, s, p] = await Promise.all([api.get('/warehouses'), api.get('/suppliers'), api.get('/products')]);
    setWarehouses(w.data || []);
    setSuppliers(s.data || []);
    setProducts(p.data || []);
  };

  const handleAddDetail = () => {
    setFormData({
      ...formData,
      details: [...formData.details, { trx_in_d_product_idf: '', trx_in_d_qty_dus: 0, trx_in_d_qty_pcs: 0 }]
    });
  };

  const handleDetailChange = (index, field, value) => {
    const newDetails = [...formData.details];
    let val = parseInt(value);
    if (field.includes('qty') && val < 0) val = 0;

    newDetails[index][field] = val;
    setFormData({ ...formData, details: newDetails });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setErrorMsg('');
    setSuccessMsg('');

    try {
      const payload = {
        ...formData,
        whs_idf: parseInt(formData.whs_idf),
        trx_in_supp_idf: parseInt(formData.trx_in_supp_idf),
      };
      await api.post('/trx-in', payload);
      alert('Transaksi berhasil disimpan!');
      setShowForm(false);
      fetchTransactions();
      // Reset form logic bisa ditambahkan di sini
    } catch (error) {
      const message = error.response?.data?.error || error.message;
      setErrorMsg(message);
      window.scrollTo(0, 0);
    }
  };

  return (
    <Container>
      <div className="d-flex justify-content-between mb-3">
        <h2>Transaksi Barang Masuk</h2>
        <Button onClick={() => setShowForm(!showForm)} variant={showForm ? "secondary" : "primary"}>
          {showForm ? "Kembali ke List" : "Tambah Transaksi"}
        </Button>
      </div>

      {showForm ? (
        <Card className="p-4">
          <Form onSubmit={handleSubmit}>
            <Row className="mb-3">
              <Col><Form.Label>No. Transaksi</Form.Label><Form.Control type="text" required onChange={e => setFormData({...formData, trx_in_no: e.target.value})} /></Col>
              <Col><Form.Label>Tanggal</Form.Label><Form.Control type="date" value={formData.trx_in_date} onChange={e => setFormData({...formData, trx_in_date: e.target.value})} /></Col>
            </Row>
            <Row className="mb-3">
              <Col>
                <Form.Label>Gudang</Form.Label>
                <Form.Select required onChange={e => setFormData({...formData, whs_idf: e.target.value})}>
                  <option value="">Pilih Gudang</option>
                  {warehouses.map(w => <option key={w.whs_pk} value={w.whs_pk}>{w.whs_name}</option>)}
                </Form.Select>
              </Col>
              <Col>
                <Form.Label>Supplier</Form.Label>
                <Form.Select required onChange={e => setFormData({...formData, trx_in_supp_idf: e.target.value})}>
                  <option value="">Pilih Supplier</option>
                  {suppliers.map(s => <option key={s.supplier_pk} value={s.supplier_pk}>{s.supplier_name}</option>)}
                </Form.Select>
              </Col>
            </Row>
            <Form.Group className="mb-3"><Form.Label>Catatan</Form.Label><Form.Control as="textarea" onChange={e => setFormData({...formData, trx_in_notes: e.target.value})} /></Form.Group>

            <h5>Detail Barang <Button size="sm" onClick={handleAddDetail}>+ Tambah Item</Button></h5>
            {formData.details.map((d, idx) => (
              <Row key={idx} className="mb-2">
                <Col md={5}>
                  <Form.Select required onChange={e => handleDetailChange(idx, 'trx_in_d_product_idf', e.target.value)}>
                    <option value="">Pilih Produk</option>
                    {products.map(p => <option key={p.product_pk} value={p.product_pk}>{p.product_name}</option>)}
                  </Form.Select>
                </Col>
                <Col><Form.Control type="number" placeholder="Dus" onChange={e => handleDetailChange(idx, 'trx_in_d_qty_dus', e.target.value)} /></Col>
                <Col><Form.Control type="number" placeholder="Pcs" onChange={e => handleDetailChange(idx, 'trx_in_d_qty_pcs', e.target.value)} /></Col>
              </Row>
            ))}
            <Button type="submit" className="mt-3" variant="success">Simpan Transaksi</Button>
          </Form>
        </Card>
      ) : (
        <Table striped bordered hover>
          <thead><tr><th>No</th><th>Tanggal</th><th>Gudang</th><th>Supplier</th><th>Notes</th></tr></thead>
          <tbody>
            {transactions.map(t => (
              <tr key={t.header.trx_in_pk}>
                <td>{t.header.trx_in_no}</td>
                <td>{new Date(t.header.trx_in_date).toLocaleDateString('id-ID')}</td>
                <td>{t.warehouse_name}</td>
                <td>{t.supplier_name}</td>
                <td>{t.header.trx_in_notes}</td>
              </tr>
            ))}
          </tbody>
        </Table>
      )}
    </Container>
  );
};

export default TransactionIn;