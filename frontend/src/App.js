import React, { useState, useEffect } from 'react';
import LoginPage from './components/LoginPage';
import RegisterPage from './components/RegisterPage';
import AssetTable from './components/AssetTable';
import { logoutApi } from './services/api';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isRegisterPage, setIsRegisterPage] = useState(false);
  
  const setToken = token => {
    localStorage.setItem('token', token);
  };
  
  const getToken = () => {
    return localStorage.getItem('token');
  };
  
  useEffect(() => {
    if (getToken()) {
      setIsAuthenticated(true);
    } else {
      setIsAuthenticated(false);
    }
  }, []);

  const onLoginSuccess = token => {
    setToken(token);
    setIsAuthenticated(true);
  };

  const onRegisterSuccess = message => {
    console.log(message);
    setIsRegisterPage(false);
  };
  
  const handleLogout = async () => {
    try {
      logoutApi();
      setIsAuthenticated(false);
      setToken('');
    } catch (error) {
      console.error('Logout error:', error);
    }
  };

  return (
    <div className="App">
      {!isAuthenticated ? (
        <>
          <div className="container mt-4">
            {isRegisterPage ? (
              <RegisterPage onRegisterSuccess={onRegisterSuccess} />
            ) : (
              <LoginPage onLoginSuccess={onLoginSuccess} />
            )}
          </div>
          <div className="text-center mt-3">
            <button onClick={() => setIsRegisterPage(!isRegisterPage)} className="btn btn-secondary">
              {isRegisterPage ? 'Go to Login' : 'Go to Register'}
            </button>
          </div>
        </>
      ) : (
        <div className="container mt-5">
          <button onClick={handleLogout} className="btn btn-danger mb-3">Logout</button>
          <AssetTable />
        </div>
      )}
    </div>
  );
}

export default App;
