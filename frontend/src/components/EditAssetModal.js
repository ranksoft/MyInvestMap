import React, { useState, useEffect } from 'react';
import { Modal, Button, Form, Alert } from 'react-bootstrap';
import axios from 'axios';

function EditAssetModal({ show, handleClose, asset, onAssetUpdated }) {
  const [updatedAsset, setUpdatedAsset] = useState({ ...asset });
  const [notification, setNotification] = useState({ message: '', type: '' });

  useEffect(() => {
    setUpdatedAsset({ ...asset });
  }, [asset]);

  const handleChange = (event) => {
    const { name, value } = event.target;
    if (name === 'price' || name === 'quantity') {
        if (value === '' || value.match(/^(\d+)?([.,](\d+)?)?$/)) {
          const formattedValue = value.replace(',', '.');
          setUpdatedAsset({ ...updatedAsset, [name]: formattedValue });
        }
      } else {
        setUpdatedAsset({ ...updatedAsset, [name]: value });
      }
  };

  useEffect(() => {
    if (show) {
      setNotification({ message: '', type: '' });
    }
  }, [show]);

  const handleSubmit = () => {
    updatedAsset.price = parseFloat(updatedAsset.price);
    updatedAsset.quantity = parseFloat(updatedAsset.quantity);
    axios.put(`${process.env.REACT_APP_BACKEND_URL}/api/assets/update/${asset.id}`, updatedAsset)
      .then(response => {
        onAssetUpdated();
        handleClose();
        setNotification({ message: 'Asset updating successfully!', type: 'success' });
      })
      .catch(error => {
        console.error('Error updating asset', error);
        setNotification({ message: 'Error updating asset. Please try again.', type: 'danger' });
    });
  };

  return (
    <Modal show={show} onHide={handleClose}>
      <Modal.Header closeButton>
        <Modal.Title>Edit Asset</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        {notification.message && (
            <Alert variant={notification.type}>{notification.message}</Alert>
        )}
        <Form>
          <Form.Group className="mb-3">
            <Form.Label>Stock Tag</Form.Label>
            <Form.Control
              type="text"
              name="stockTag"
              value={updatedAsset.stockTag}
              onChange={handleChange}
            />
          </Form.Group>

          <Form.Group className="mb-3">
            <Form.Label>Exchange</Form.Label>
            <Form.Control
              type="text"
              name="exchange"
              value={updatedAsset.exchange}
              onChange={handleChange}
            />
          </Form.Group>

          <Form.Group className="mb-3">
            <Form.Label>Price</Form.Label>
            <Form.Control
              type="number"
              name="price"
              value={updatedAsset.price}
              onChange={handleChange}
            />
          </Form.Group>

          <Form.Group className="mb-3">
            <Form.Label>Quantity</Form.Label>
            <Form.Control
              type="number"
              name="quantity"
              value={updatedAsset.quantity}
              onChange={handleChange}
            />
          </Form.Group>
        </Form>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={handleClose}>
          Close
        </Button>
        <Button variant="primary" onClick={handleSubmit}>
          Save Changes
        </Button>
      </Modal.Footer>
    </Modal>
  );
}

export default EditAssetModal;
