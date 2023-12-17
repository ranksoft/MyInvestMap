import React, { useState, useEffect } from 'react';
import LoginPage from './components/LoginPage';
import RegisterPage from './components/RegisterPage';
import AssetTable from './components/AssetTable';
import { jwtDecode } from "jwt-decode";
import { logoutApi } from './services/api';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isRegisterPage, setIsRegisterPage] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');

  useEffect(() => {
    const token = localStorage.getItem('token');
    setIsAuthenticated(token && !isTokenExpired(token));
  }, []);

  const isTokenExpired = (token) => {
    if (!token) return true;
    const decoded = jwtDecode(token);
    return decoded.exp < Date.now() / 1000;
  };

  const handleLoginSuccess = (token) => {
    localStorage.setItem('token', token);
    setIsAuthenticated(true);
    setErrorMessage('');
  };

  const handleRegisterSuccess = () => {
    setIsRegisterPage(false);
    setErrorMessage('');
  };

  const handleLogout = async () => {
    try {
      if (isAuthenticated) {
        await logoutApi();
        localStorage.removeItem('token');
        setIsAuthenticated(false);
        setIsRegisterPage(false);
        setErrorMessage('');
      }
    } catch (error) {
      setErrorMessage('Logout error: ' + error.message);
      console.error('Logout error:', error);
    }
  };

  return (
    <div className="App">
      {errorMessage && <div className="alert alert-danger">{errorMessage}</div>}
      {!isAuthenticated ? (
        <div className="container mt-4">
          {isRegisterPage ? 
            <RegisterPage onRegisterSuccess={handleRegisterSuccess} /> : 
            <LoginPage onLoginSuccess={handleLoginSuccess} />}
          <div className="text-center mt-3">
            <button onClick={() => setIsRegisterPage(!isRegisterPage)} className="btn btn-secondary">
              {isRegisterPage ? 'Go to Login' : 'Go to Register'}
            </button>
          </div>
        </div>
      ) : (
        <div className="container mt-5">
          <div className="d-flex justify-content-center mb-4">
            <button onClick={handleLogout} className="btn btn-danger mb-3">Logout</button>
          </div>
          <AssetTable />
        </div>
      )}
    </div>
  );
}

export default App;
