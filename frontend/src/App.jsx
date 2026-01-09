// src/App.jsx
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import MyNavbar from './components/MyNavbar';
import MasterData from './pages/MasterData';
import TransactionIn from './pages/TransactionIn';
import TransactionOut from './pages/TransactionOut';
import StockReport from './pages/StockReport';

function App() {
  return (
    <Router>
      <MyNavbar />
      <Routes>
        <Route path="/" element={<MasterData />} />
        <Route path="/trx-in" element={<TransactionIn />} />
        <Route path="/trx-out" element={<TransactionOut />} />
        <Route path="/stock" element={<StockReport />} />
      </Routes>
    </Router>
  );
}

export default App;