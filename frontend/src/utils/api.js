// src/utils/api.js
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api', // Sesuaikan dengan port backend Go kamu
});

export default api;