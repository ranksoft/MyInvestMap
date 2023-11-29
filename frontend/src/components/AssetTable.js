import React, { useState, useEffect } from 'react';
import axios from 'axios';
import AddAssetForm from './AddAssetForm';
import SellAssetForm from './SellAssetForm';
import EditAssetModal from './EditAssetModal';
import { Table } from 'react-bootstrap';

function AssetTable() {
  const [assets, setAssets] = useState([]);
  const [editingAsset, setEditingAsset] = useState(null);
  const [showEditModal, setShowEditModal] = useState(false);
  const [stockData, setStockData] = useState({});
  const [isLoadingStockData, setIsLoadingStockData] = useState(true);


  const handleEdit = (asset) => {
    setEditingAsset(asset);
    setShowEditModal(true);
  };

  const handleCloseEditModal = () => {
    setShowEditModal(false);
    setEditingAsset(null);
  };

  const fetchAssets = () => {
    axios.get(`${process.env.REACT_APP_BACKEND_URL}/api/assets`)
      .then(response => {
        setAssets(response.data);
      })
      .catch(error => console.error('Error fetching assets:', error));
      console.log(assets);
  };
  useEffect(() => {
    fetchAssets();
  }, []);

  const deleteAsset = (id) => {
    axios.delete(`${process.env.REACT_APP_BACKEND_URL}/api/assets/delete/${id}`)
      .then(response => {
        fetchAssets();
      })
      .catch(error => console.error('Error deleting asset:', error));
  };

  const getName = (asset) => {
    return asset.name.Valid ? asset.name.String : 'N/A';
  };
  
  const getCurrentPrice = (asset) => {
    return asset.currentPrice.Valid ? asset.currentPrice.Float64 : 0;
  };
  
  const calculateInvestment = (asset) => {
      return asset.price * asset.quantity;
  };

  const calculateProfitOrLoss = (asset) => {
      const currentPrice = getCurrentPrice(asset);
      if (!currentPrice) {
        return 0;
      }
      let currentTotal = currentPrice * asset.quantity;
      let initialTotal = asset.price * asset.quantity;
      return currentTotal - initialTotal;
  };

  const calculateTotalInvestment = () => {
    return assets?.reduce((total, asset) => {
      const assetQuantity = parseFloat(asset.quantity);
      const assetPrice = parseFloat(asset.price);
      return asset.isPurchase ? total + (assetPrice * assetQuantity) : total;
    }, 0).toFixed(4);
  };
  
  const calculateTotalProfitLoss = () => {
    return assets?.reduce((total, asset) => {
      if (!asset.isPurchase) {
        return total;
      }

      const currentPrice = parseFloat(getCurrentPrice(asset));
      if (!currentPrice) {
        return 0;
      }
      const assetQuantity = parseFloat(asset.quantity);
      const assetPrice = parseFloat(asset.price);
      const profitOrLoss = (currentPrice - assetPrice) * assetQuantity;
      return total + profitOrLoss;
    }, 0).toFixed(4);
  };
  
  const calculatePortfolioValue = () => {
    return parseFloat(calculateTotalInvestment()) + parseFloat(calculateTotalProfitLoss());
  };

  return (
    <div className="container mt-4">
    <h2 className="mb-4">MyInvestMap Portfolio</h2>
    <AddAssetForm onAssetAdded={fetchAssets} />
    <SellAssetForm onAssetSold={fetchAssets} />
    <Table striped bordered hover responsive="lg">
      <thead className="bg-light">
        <tr>
          <th>Is Purchase</th>
          <th>Stock Tag</th>
          <th>Exchange</th>
          <th>Name</th>
          <th>Price</th>
          <th>Quantity</th>
          <th>Current Price</th>
          <th>Total Investment</th>
          <th>Total Profit/Loss</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
      {assets?.map(asset => (
              <tr key={asset.id}>
              {asset.isPurchase 
              ? <td className="text-primary">Purchase</td> 
              : <td className="text-warning">Sale</td>}
              <td>{asset.stockTag}</td>
              <td>{asset.exchange}</td>
              <td>{getName(asset)}</td>
              <td>{asset.price}</td>
              <td>{asset.quantity}</td>
              <td>{getCurrentPrice(asset)}</td>
              <td>{calculateInvestment(asset)}</td>
              <td className={calculateProfitOrLoss(asset) >= 0 ? 'text-success' : 'text-danger'}>
                {calculateProfitOrLoss(asset)}
                </td>
              <td>
              <button className="btn btn-info btn-spacing" onClick={() => handleEdit(asset)}>
                <i className="fas fa-edit"></i>
              </button>
              <button className="btn btn-danger" onClick={() => deleteAsset(asset.id)}>
                <i className="fas fa-trash-alt"></i>
              </button>
              </td>
            </tr>
      ))}
      </tbody>
      <tfoot className="bg-light">
        <tr className="table-primary">
          <td colSpan="8"></td>
          <td >Total investment:</td>
          <td>{formatCurrency(calculateTotalInvestment())}</td>
        </tr>
        <tr className="table-primary">
          <td colSpan="8"></td>
          <td >Total profit/loss:</td>
          <td className={parseFloat(calculateTotalProfitLoss()) >= 0 ? 'text-success' : 'text-danger'}>
            {formatCurrency(calculateTotalProfitLoss())}
          </td>
        </tr>
        <tr className="table-primary">
          <td colSpan="8"></td>
          <td>Total portfolio value:</td>
          <td>{formatCurrency(calculatePortfolioValue())}</td>
        </tr>
      </tfoot>
    </Table>
    <EditAssetModal
      show={showEditModal}
      handleClose={handleCloseEditModal}
      asset={editingAsset}
      onAssetUpdated={fetchAssets}
    />
  </div>
  );
}

function formatCurrency(value) {
  return parseFloat(value).toFixed(2).replace(/\d(?=(\d{3})+\.)/g, '$&,');
}

export default AssetTable;
