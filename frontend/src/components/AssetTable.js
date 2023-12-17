import React, { useState, useEffect, useCallback} from 'react';
import AddAssetForm from './AddAssetForm';
import SellAssetForm from './SellAssetForm';
import EditAssetModal from './EditAssetModal';
import ApiKeyForm from './ApiKeyForm';
import { Table } from 'react-bootstrap';
import { getAssetsApi, deleteAssetApi, refreshAssetsApi } from '../services/api';

function AssetTable() {
  const [assets, setAssets] = useState([]);
  const [editingAsset, setEditingAsset] = useState(null);
  const [showEditModal, setShowEditModal] = useState(false);
  const [selectedAssets, setSelectedAssets] = useState(new Set());
  const [errorMessage, setErrorMessage] = useState('');

  const handleEdit = (asset) => {
    setEditingAsset(asset);
    setShowEditModal(true);
  };

  const handleCloseEditModal = () => {
    setShowEditModal(false);
    setEditingAsset(null);
  };

  const toggleAssetSelection = (assetId) => {
    setSelectedAssets((prevSelectedAssets) => {
      const newSelectedAssets = new Set(prevSelectedAssets);
      if (newSelectedAssets.has(assetId)) {
        newSelectedAssets.delete(assetId);
      } else {
        if (newSelectedAssets.size >= 8) {
          alert("You can select a maximum of 8 assets.");
          return prevSelectedAssets;
        }
        newSelectedAssets.add(assetId);
      }
      return newSelectedAssets;
    });
  };

  function handleRefreshSelected() {
    if (selectedAssets.length > 8) {
      alert("You can select a maximum of 8 assets for refresh.");
      return;
    }

    const selectedSymbols = assets.filter(asset => selectedAssets.has(asset.id)).map(asset => asset.stockTag);
    if (selectedSymbols.length > 0) {
      const symbolsToUpdate = selectedSymbols.slice(0, 8);
      refreshAssetsApi({
        symbols: symbolsToUpdate
      })
      .then(response => {
        console.log('Assets updated:', response.data);
        fetchAssets();
        setSelectedAssets(new Set())
      })
      .catch(error => () => {
        if (error.message) {
          let errorMessage = error.message || 'Error occurred updating assets';
          setErrorMessage('Error updating assets: ' + errorMessage);
          console.error('Error updating assets:', errorMessage);
      } else {
          setErrorMessage('Error updating assets: An unexpected error occurred');
          console.error('Error updating assets:', error);
      }
      });
    }
  }

  const fetchAssets = useCallback(() => {
    getAssetsApi()
    .then(response => {
      setAssets(response.data);
    })
    .catch(error => () => {
      if (error.message) {
        let errorMessage = error.message || 'Error occurred fetching assets';
        setErrorMessage('Error fetching assets: ' + errorMessage);
        console.error('Error fetching assets:', errorMessage);
    } else {
        setErrorMessage('Error fetching assets: An unexpected error occurred');
        console.error('Error fetching assets:', error);
    }
    });
  }, []);

  useEffect(() => {
    fetchAssets();
    // const interval = setInterval(fetchAssets, 180000); 
    // return () => clearInterval(interval); 
  }, [fetchAssets]);

  const deleteAsset = (id) => {
    deleteAssetApi(id)
    .then(response => {
      fetchAssets();
    })
    .catch(error => () => {
      if (error.message) {
        let errorMessage = error.message || 'Error occurred deleting assets';
        setErrorMessage('Error deleting assets: ' + errorMessage);
        console.error('Error deleting assets:', errorMessage);
    } else {
        setErrorMessage('Error deleting assets: An unexpected error occurred');
        console.error('Error deleting assets:', error);
    }
    });
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

  const calculateProfitOrLossPercentage = (asset) => {
      return (calculateProfitOrLoss(asset) / calculateInvestment(asset)) * 100;
  }

  const calculateTotalInvestment = () => {
    return assets?.reduce((total, asset) => {
      const assetQuantity = parseFloat(asset.quantity);
      const assetPrice = parseFloat(asset.price);
      return asset.isPurchase ? total + (assetPrice * assetQuantity) : total - (assetPrice * assetQuantity);
    }, 0).toFixed(4);
  };
  
  const calculateTotalProfitLoss = () => {
    return assets?.reduce((total, asset) => {

      const currentPrice = parseFloat(getCurrentPrice(asset));
      if (!currentPrice) {
        return total + 0;
      }
      const assetQuantity = parseFloat(asset.quantity);
      const assetPrice = parseFloat(asset.price);
      const profitOrLoss = (currentPrice - assetPrice) * assetQuantity;
      if (!asset.isPurchase) {
        return total - calculateProfitOrLoss(asset);
      }
      return total + profitOrLoss;
    }, 0).toFixed(4);
  };
  
  const calculatePortfolioValue = () => {
    return parseFloat(calculateTotalInvestment()) + parseFloat(calculateTotalProfitLoss());
  };

  const calculateTotalProfitLossPercentage = () => {
    return (calculateTotalProfitLoss() / calculateTotalInvestment()) * 100;
  }

  const countUniqueAssetTags = () => {
    const uniqueTags = new Set();
    assets.forEach(asset => uniqueTags.add(asset.stockTag));
    return uniqueTags.size;
  };

  return (
    <div className="container mt-4">
    {errorMessage && <div className="alert alert-danger">{errorMessage}</div>}
    <ApiKeyForm />
    <h2 className="mb-4">MyInvestMap Portfolio</h2>

    <div className="d-flex justify-content-between mb-4">
      <AddAssetForm onAssetAdded={fetchAssets} />
      <SellAssetForm onAssetSold={fetchAssets} />
    </div>
    <Table striped bordered hover responsive="lg">
      <thead className="bg-light">
        <tr>
          <th>
            <button className="btn btn-lin" onClick={handleRefreshSelected}>
              <i className="fas fa-sync"></i>
            </button>
          </th>
          <th colspan="10" className="text-center">MyInvestMap</th>
        </tr>
        <tr>
          <th></th>
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
              <tr onClick={() => toggleAssetSelection(asset.id)}>
              <td onClick={(e) => e.stopPropagation()}>
                <input 
                  type="checkbox" 
                  className="btn-check" 
                  id={'btn-check-' + asset.id + '-outlined'}
                  checked={selectedAssets.has(asset.id)}
                  onChange={() => toggleAssetSelection(asset.id)}
                />
                <label className="btn btn-outline-secondary" for={'btn-check-' + asset.id + '-outlined'}>
                  <i className="far fa-check-square fa-lg"></i>
                </label>
              </td>
              {asset.isPurchase 
              ? <td className="text-primary">Purchase</td> 
              : <td className="text-warning">Sale</td>}
              <td>{asset.stockTag}</td>
              <td>{asset.exchange}</td>
              <td>{getName(asset)}</td>
              <td>{asset.price}</td>
              <td>{asset.quantity}</td>
              <td>{formatCurrency(getCurrentPrice(asset))}</td>
              <td>{formatCurrency(calculateInvestment(asset))}</td>
              <td className={calculateProfitOrLoss(asset) >= 0 ? 'text-success' : 'text-danger'}>
                {formatCurrency(calculateProfitOrLoss(asset))} ({formatCurrency(calculateProfitOrLossPercentage(asset))}%)
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
          <td colSpan="2" className="text-end" >Total investment:</td>
          <td>{formatCurrency(calculateTotalInvestment())}</td>
        </tr>
        <tr className="table-primary">
          <td colSpan="8"></td>
          <td colSpan="2" className="text-end" >Total profit/loss:</td>
          <td className={parseFloat(calculateTotalProfitLoss()) >= 0 ? 'text-success' : 'text-danger'}>
            {formatCurrency(calculateTotalProfitLoss())} ({formatCurrency(calculateTotalProfitLossPercentage())}%) 
          </td>
        </tr>
        <tr className="table-primary">
          <td colSpan="8"></td>
          <td colSpan="2" className="text-end" >Total portfolio value:</td>
          <td>{formatCurrency(calculatePortfolioValue())}</td>
        </tr>
        <tr className="table-secondary">
          <td colSpan="8"></td>
          <td colSpan="2" className="text-end" >Total Unique Assets:</td>
          <td>{countUniqueAssetTags()}</td>
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

function formatCurrency(value, define = 2) {
  return parseFloat(value).toFixed(define).replace(/\d(?=(\d{3})+\.)/g, '$&,');
}

export default AssetTable;
