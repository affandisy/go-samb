// src/components/MyNavbar.jsx
import { Navbar, Container, Nav } from 'react-bootstrap';
import { Link } from 'react-router-dom';

const MyNavbar = () => {
  return (
    <Navbar bg="dark" variant="dark" expand="lg" className="mb-4">
      <Container>
        <Navbar.Brand as={Link} to="/">Go SAMB</Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="me-auto">
            <Nav.Link as={Link} to="/">Master Data</Nav.Link>
            <Nav.Link as={Link} to="/trx-in">Barang Masuk</Nav.Link>
            <Nav.Link as={Link} to="/trx-out">Barang Keluar</Nav.Link>
            <Nav.Link as={Link} to="/stock">Laporan Stok</Nav.Link>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
};

export default MyNavbar;