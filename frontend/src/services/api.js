import secureAxios from '../axios';
import axios from 'axios';

const notSecureinstance = axios.create({
    baseURL: process.env.REACT_APP_BACKEND_URL
});

const loginApi = (data) => {
    console.log(data);
    return notSecureinstance.post('/api/login', data);
};

const registerApi = (data) => {
    return notSecureinstance.post('/api/register', data);
};

const getApiKey = async () => {
    return secureAxios.get('/api/api-key');
};
  
const saveApiKey = async (apiKey) => {
    return secureAxios.post('/api/api-key', { api_key: apiKey });
}

const logoutApi = () => {
    return secureAxios.get('/api/logout');
};

const addAssetApi = (asset) => {
    return secureAxios.post('/api/assets/add', asset);
};

const addSellAssetApi = (asset) => {
    return secureAxios.post('/api/assets/sell', asset);
};

const deleteAssetApi = (assetId) => {
    return secureAxios.delete(`/api/assets/delete/${assetId}`);
};

const updateAssetApi = (assetId, asset) => {
    return secureAxios.put(`/api/assets/update/${assetId}`, asset);
};

const getAssetsApi = () => {
    return secureAxios.get('/api/assets');
};

const refreshAssetsApi = (data) => {
    return secureAxios.post('/api/refresh-assets', data);
};


export { getApiKey, saveApiKey, loginApi, registerApi, logoutApi, addAssetApi, addSellAssetApi, deleteAssetApi, updateAssetApi, getAssetsApi, refreshAssetsApi };
