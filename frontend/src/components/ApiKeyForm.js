import React, { useState, useEffect } from 'react';
import { getApiKey, saveApiKey } from '../services/api';

function ApiKeyForm() {
  const [apiKey, setApiKey] = useState('');

  useEffect(() => {
    getApiKey(apiKey)
    .then(response => {
        setApiKey(response.data.api_key)
    })
    .catch(error => {
        console.error('Error fetching API key:', error);
    });
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    await saveApiKey(apiKey)
    .then(response => {
        setApiKey(response.data.api_key)
    })
    .catch(error => {
        console.error('Error saving API key:', error);
    });
  };

  return (
    <form className='my-4' onSubmit={handleSubmit}>
      <div className='mb-3'>
      <input 
        className='form-control'
        type="text" 
        value={apiKey} 
        onChange={(e) => setApiKey(e.target.value)} 
        placeholder="Enter API Key" 
      />
      <button className='btn btn-primary' type="submit">Save API Key</button>
      </div>
    </form>
  );
}

export default ApiKeyForm;
