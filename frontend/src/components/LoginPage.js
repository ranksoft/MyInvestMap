import React, { useState } from 'react';
import { loginApi } from '../services/api';

function LoginPage({ onLoginSuccess }) {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [errorMessage, setErrorMessage] = useState('');

    const handleLogin = async (e) => {
        e.preventDefault();
        try {
            const response = await loginApi({ email, password });
            onLoginSuccess(response.data.token);
        } catch (error) {
            if (error.message) {
                let errorMessage = error.message || 'Error occurred during login';
                setErrorMessage('Login error: ' + errorMessage);
                console.error('Login error:', errorMessage);
            } else {
                setErrorMessage('Login error: An unexpected error occurred');
                console.error('Login error:', error);
            }
        }
    };

    return (
        <div className="container mt-5">
        <div className="row justify-content-center">
            <div className="col-md-6">
            <div className="card">
                <div className="card-body">
                <h5 className="card-title">Login</h5>
                <form>
                    <div className="mb-3">
                    <label htmlFor="email" className="form-label">Email address</label>
                    <input type="email" className="form-control" id="email" placeholder="Enter email" value={email} onChange={e => setEmail(e.target.value)} />
                    </div>
                    <div className="mb-3">
                    <label htmlFor="password" className="form-label">Password</label>
                    <input type="password" className="form-control" id="password" placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} />
                    </div>
                    {errorMessage && <div className="alert alert-danger">{errorMessage}</div>}
                    <button type="submit" className="btn btn-primary" onClick={handleLogin}>Login</button>
                </form>
                </div>
            </div>
            </div>
        </div>
        </div>
    );
}

export default LoginPage;
