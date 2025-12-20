import React, { useEffect } from 'react';
import { Routes, Route, Navigate, useLocation } from 'react-router-dom';
import { Container } from '@mui/material';
import { useDispatch, useSelector } from 'react-redux';
import Header from './components/Header';
import HomePage from './pages/HomePage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import NewsPage from './pages/NewsPage';
import SubscriptionsPage from './pages/SubscriptionsPage';
import AdminPage from './pages/AdminPage';
import PrivateRoute from './components/PrivateRoute';
import { getCurrentUser } from './store/authSlice';

function App() {
  const dispatch = useDispatch();
  const location = useLocation();
  const { isAuthenticated } = useSelector((state) => state.auth);

  useEffect(() => {
    console.log('🚀 App mounted, checking auth...');
    
    const token = localStorage.getItem('token');
    console.log('🔍 Token in localStorage:', token ? 'Exists' : 'Missing');
    
    if (token) {
      console.log('🔄 Dispatching getCurrentUser...');
      dispatch(getCurrentUser());
    } else {
      console.log('⚠️ No token found, user is not authenticated');
    }
  }, [dispatch]);

  // Логируем изменение пути
  useEffect(() => {
    console.log('📍 Route changed to:', location.pathname);
    console.log('🔐 Auth state:', isAuthenticated ? 'Authenticated' : 'Not authenticated');
  }, [location, isAuthenticated]);

  return (
    <>
      <Header />
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/news" element={
            <PrivateRoute>
              <NewsPage />
            </PrivateRoute>
          } />
          <Route path="/subscriptions" element={
            <PrivateRoute>
              <SubscriptionsPage />
            </PrivateRoute>
          } />
          <Route path="/admin" element={
            <PrivateRoute adminOnly>
              <AdminPage />
            </PrivateRoute>
          } />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Container>
    </>
  );
}

export default App;