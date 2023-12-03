import React, { useState, useEffect } from 'react';
import { Modal, Button, Form, Alert } from 'react-bootstrap';
import { addSellAssetApi } from '../services/api';

function SellAssetForm({ onAssetSold }) {
  const [asset, setAsset] = useState({ stockTag: '', exchange: '', price: 0, quantity: 0 });
  const [showModal, setShowModal] = useState(false);
  const [notification, setNotification] = useState({ message: '', type: '' });

  const handleChange = (event) => {
    const { name, value } = event.target;
    if (name === 'price' || name === 'quantity') {
      if (value === '' || value.match(/^(\d+)?([.,](\d+)?)?$/)) {
        const formattedValue = value.replace(',', '.');
        setAsset({ ...asset, [name]: formattedValue });
      }
    } else {
      setAsset({ ...asset, [name]: value });
    }
  };

  useEffect(() => {
    if (showModal) {
      setNotification({ message: '', type: '' });
    }
  }, [showModal]);

  const handleSubmit = (event) => {
    asset.price = parseFloat(asset.price);
    asset.quantity = parseFloat(asset.quantity);
    event.preventDefault();
    addSellAssetApi(asset)
    .then(response => {
        onAssetSold();
        setShowModal(false);
        setAsset({ stockTag: '', exchange: '', price: '', quantity: '' });
        setNotification({ message: 'Asset sold successfully!', type: 'success' });
    })
    .catch(error => {
        console.error('Error selling asset:', error);
        setNotification({ message: 'Error selling asset. Please try again.', type: 'danger' });
    });
  };

  return (
    <>
      <Button className="btn-bottom-spacing" variant="warning" onClick={() => setShowModal(true)}>
        Sell Asset
      </Button>

      <Modal show={showModal} onHide={() => setShowModal(false)}>
        <Modal.Header closeButton>
          <Modal.Title>Sell Asset</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {notification.message && (
            <Alert variant={notification.type}>{notification.message}</Alert>
          )}
          <Form onSubmit={handleSubmit}>
            <Form.Group className="mb-3">
              <Form.Label>Stock Tag</Form.Label>
              <Form.Control 
                type="text" 
                name="stockTag" 
                value={asset.stockTag} 
                onChange={handleChange} 
                placeholder="Stock Tag" 
              />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Exchange</Form.Label>
              <Form.Control 
                type="text" 
                name="exchange" 
                value={asset.exchange} 
                onChange={handleChange} 
                placeholder="Exchange" 
              />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Price</Form.Label>
              <Form.Control 
                type="text" 
                name="price" 
                value={asset.price} 
                onChange={handleChange} 
                placeholder="Price" 
              />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Quantity</Form.Label>
              <Form.Control 
                type="text" 
                name="quantity" 
                value={asset.quantity} 
                onChange={handleChange} 
                placeholder="Quantity" 
              />
            </Form.Group>
            <Button variant="warning" type="submit">
              Sell Asset
            </Button>
          </Form>
        </Modal.Body>
      </Modal>
    </>
  );
}

export default SellAssetForm;
