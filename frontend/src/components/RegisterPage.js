import React, { useState } from 'react';
import { registerApi } from '../services/api';

function RegisterPage({ onRegisterSuccess }) {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [username, setUsername] = useState('');

    const handleRegister = async () => {
        try {
            const response = await registerApi({ email, password, username });
            onRegisterSuccess(response.data.message);
        } catch (error) {
            console.error('Register error:', error.response.data);
        }
    };

    return (
        <div className="container mt-5">
        <div className="row justify-content-center">
            <div className="col-md-6">
            <div className="card">
                <div className="card-body">
                <h5 className="card-title">Register</h5>
                <form>
                    <div className="mb-3">
                    <label htmlFor="username" className="form-label">Username</label>
                    <input type="text" className="form-control" id="username" placeholder="Username" value={username} onChange={e => setUsername(e.target.value)} />
                    </div>
                    <div className="mb-3">
                    <label htmlFor="email" className="form-label">Email address</label>
                    <input type="email" className="form-control" id="email" placeholder="Enter email" value={email} onChange={e => setEmail(e.target.value)} />
                    </div>
                    <div className="mb-3">
                    <label htmlFor="password" className="form-label">Password</label>
                    <input type="password" className="form-control" id="password" placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} />
                    </div>
                    <button type="submit" className="btn btn-success" onClick={handleRegister}>Register</button>
                </form>
                </div>
            </div>
            </div>
        </div>
        </div>
    );
}

export default RegisterPage;
