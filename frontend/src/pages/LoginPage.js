import { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate, Link as RouterLink, useLocation } from 'react-router-dom';
import {
  Container,
  Paper,
  TextField,
  Button,
  Typography,
  Box,
  Alert,
  Link,
  CircularProgress,
} from '@mui/material';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { loginUser, clearError } from '../store/authSlice';

const schema = yup.object({
  email: yup.string().email('Некорректный email').required('Email обязателен'),
  password: yup.string().min(6, 'Минимум 6 символов').required('Пароль обязателен'),
});

function LoginPage() {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const location = useLocation();
  const { loading, error, isAuthenticated } = useSelector((state) => state.auth);
  
  const [showSuccess, setShowSuccess] = useState(false);
  const [loginAttempted, setLoginAttempted] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    resolver: yupResolver(schema),
  });

  // Проверяем, если пользователь уже авторизован
  useEffect(() => {
    console.log('🔍 LoginPage useEffect - isAuthenticated:', isAuthenticated);
    
    if (isAuthenticated && loginAttempted) {
      console.log('✅ Login successful, redirecting...');
      setShowSuccess(true);
      
      // Редирект через 1 секунду
      const timer = setTimeout(() => {
        const from = location.state?.from?.pathname || '/news';
        console.log('🚀 Redirecting to:', from);
        navigate(from, { replace: true });
      }, 1000);
      
      return () => clearTimeout(timer);
    }
  }, [isAuthenticated, navigate, location, loginAttempted]);

  const onSubmit = async (data) => {
    console.log('📝 Login form submitted:', data.email);
    setLoginAttempted(true);
    dispatch(clearError());
    
    try {
      const result = await dispatch(loginUser(data));
      console.log('📊 Login result:', result);
      
      if (result.meta.requestStatus === 'fulfilled') {
        console.log('🎉 Login successful in onSubmit');
        // Редирект будет выполнен в useEffect
      } else {
        console.log('❌ Login failed in onSubmit');
        setLoginAttempted(false);
      }
    } catch (err) {
      console.error('💥 Login error:', err);
      setLoginAttempted(false);
    }
  };

  // Показываем успешный вход
  if (showSuccess) {
    return (
      <Container maxWidth="sm">
        <Box sx={{ mt: 8, textAlign: 'center' }}>
          <Paper sx={{ p: 4 }}>
            <CircularProgress sx={{ mb: 2 }} />
            <Typography variant="h5" gutterBottom>
              Вход выполнен успешно!
            </Typography>
            <Typography variant="body1" color="text.secondary">
              Перенаправляем на страницу новостей...
            </Typography>
          </Paper>
        </Box>
      </Container>
    );
  }

  return (
    <Container maxWidth="sm">
      <Box sx={{ mt: 8 }}>
        <Paper sx={{ p: 4 }}>
          <Typography variant="h4" component="h1" align="center" gutterBottom>
            Вход в систему
          </Typography>
          
          {/* Отладочная информация */}
          {/* <Box sx={{ mb: 2, p: 2, bgcolor: 'grey.100', borderRadius: 1 }}>
            <Typography variant="caption" display="block">
              Отладка: 
              <br />• Токен в localStorage: {localStorage.getItem('token') ? 'Есть' : 'Нет'}
              <br />• Авторизован: {isAuthenticated ? 'Да' : 'Нет'}
              <br />• Бэкенд: https://localhost:8443
            </Typography>
          </Box> */}

          {error && (
            <Alert 
              severity="error" 
              sx={{ mb: 2 }} 
              onClose={() => dispatch(clearError())}
            >
              <Typography variant="body2">{error}</Typography>
            </Alert>
          )}

          {location.search.includes('sessionExpired') && (
            <Alert severity="warning" sx={{ mb: 2 }}>
              Ваша сессия истекла. Пожалуйста, войдите снова.
            </Alert>
          )}

          <form onSubmit={handleSubmit(onSubmit)}>
            <TextField
              fullWidth
              label="Email"
              type="email"
              {...register('email')}
              error={!!errors.email}
              helperText={errors.email?.message}
              margin="normal"
              disabled={loading}
              autoComplete="email"
            />
            
            <TextField
              fullWidth
              label="Пароль"
              type="password"
              {...register('password')}
              error={!!errors.password}
              helperText={errors.password?.message}
              margin="normal"
              disabled={loading}
              autoComplete="current-password"
            />
            
            <Box sx={{ mt: 3, mb: 2 }}>
              <Button
                type="submit"
                variant="contained"
                fullWidth
                size="large"
                disabled={loading}
              >
                {loading ? (
                  <>
                    <CircularProgress size={24} sx={{ mr: 1 }} />
                    Вход...
                  </>
                ) : 'Войти'}
              </Button>
            </Box>
            
            <Box sx={{ textAlign: 'center', mt: 2 }}>
              <Typography variant="body2">
                Нет аккаунта?{' '}
                <Link component={RouterLink} to="/register">
                  Зарегистрироваться
                </Link>
              </Typography>
            </Box>
          </form>
        </Paper>
        
        {/* Кнопка для тестирования */}
        {/* <Box sx={{ mt: 2, textAlign: 'center' }}>
          <Button
            variant="outlined"
            size="small"
            onClick={() => {
              console.log('🧪 Debug:');
              console.log('Token:', localStorage.getItem('token'));
              console.log('User:', localStorage.getItem('user'));
              console.log('Backend reachable?');
              
              // Проверяем доступность бэкенда
              fetch('http://localhost:8080/auth/login', {
                method: 'OPTIONS'
              })
              .then(res => {
                console.log('Backend OPTIONS response:', res.status);
                if (res.ok) {
                  console.log('✅ Backend is reachable');
                } else {
                  console.log('❌ Backend not reachable');
                }
              })
              .catch(err => {
                console.error('Backend check failed:', err);
              });
            }}
          >
            Проверить подключение к бэкенду
          </Button>
        </Box> */}
      </Box>
    </Container>
  );
}

export default LoginPage;