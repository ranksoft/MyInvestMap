import React, { useState, useEffect } from 'react';
import { getApiKey, saveApiKey } from '../services/api';
import { Button, Form, FormGroup, FormControl, Container, Row, Col } from 'react-bootstrap';


function ApiKeyForm() {
  const [apiKey, setApiKey] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  useEffect(() => {
    getApiKey(apiKey)
    .then(response => {
        setApiKey(response.data.api_key)
    })
    .catch(error => {
      if (error.message) {
        let errorMessage = error.message || 'Error occurred fetching API key';
        setErrorMessage('Error fetching API key: ' + errorMessage);
        console.error('Error fetching API key:', errorMessage);
      } else {
          setErrorMessage('Error fetching API key: An unexpected error occurred');
          console.error('Error fetching API key:', error);
      }
    });
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    await saveApiKey(apiKey)
    .then(response => {
        setApiKey(response.data.api_key)
    })
    .catch(error => {
      if (error.message) {
        let errorMessage = error.message || 'Error occurred saving API key';
        setErrorMessage('Error saving API key: ' + errorMessage);
        console.error('Error saving API key:', errorMessage);
      } else {
          setErrorMessage('Error saving API key: An unexpected error occurred');
          console.error('Error saving API key:', error);
      }
    });
  };

  return (
    <form className='my-4' onSubmit={handleSubmit}>
      <div className='mb-3'>
       <FormGroup className='mb-3'>
            <FormControl
                type="text"
                value={apiKey}
                onChange={(e) => setApiKey(e.target.value)}
                placeholder="Enter API Key"
            />
        </FormGroup>
      <Button className='btn btn-primary' type="submit">Save API Key</Button>
      </div>
    </form>
  );
}

export default ApiKeyForm;
