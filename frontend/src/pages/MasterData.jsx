// src/pages/MasterData.jsx
import { useEffect, useState } from 'react';
import { Container, Row, Col, Card, ListGroup } from 'react-bootstrap';
import api from '../utils/api';

const MasterData = () => {
  const [data, setData] = useState({
    suppliers: [],
    customers: [],
    products: [],
    warehouses: []
  });

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [supp, cust, prod, whs] = await Promise.all([
          api.get('/suppliers'),
          api.get('/customers'),
          api.get('/products'),
          api.get('/warehouses')
        ]);
        
        setData({
          suppliers: supp.data || [],
          customers: cust.data || [],
          products: prod.data || [],
          warehouses: whs.data || []
        });
      } catch (error) {
        console.error("Gagal mengambil data master", error);
      }
    };
    fetchData();
  }, []);

  const renderCard = (title, items, keyName) => (
    <Col md={3} className="mb-3">
      <Card>
        <Card.Header className="fw-bold">{title}</Card.Header>
        <ListGroup variant="flush">
          {items.map((item, idx) => (
            <ListGroup.Item key={idx}>{item[keyName]}</ListGroup.Item>
          ))}
          {items.length === 0 && <ListGroup.Item>Tidak ada data</ListGroup.Item>}
        </ListGroup>
      </Card>
    </Col>
  );

  return (
    <Container>
      <h2>Master Data Dashboard</h2>
      <Row>
        {renderCard("Suppliers", data.suppliers, "supplier_name")}
        {renderCard("Customers", data.customers, "customer_name")}
        {renderCard("Products", data.products, "product_name")}
        {renderCard("Warehouses", data.warehouses, "whs_name")}
      </Row>
    </Container>
  );
};

export default MasterData;