import { useSelector } from 'react-redux';
import { Typography, Button, Box, Paper } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

function HomePage() {
  const { isAuthenticated } = useSelector((state) => state.auth);

  return (
    <Box sx={{ textAlign: 'center', py: 8 }}>
      <Typography variant="h2" component="h1" gutterBottom>
        News Aggregator
      </Typography>
      <Typography variant="h5" color="text.secondary" paragraph>
        Собирайте новости из всех ваших любимых источников в одном месте
      </Typography>
      
      <Paper sx={{ p: 4, mt: 4, maxWidth: 600, mx: 'auto' }}>
        {isAuthenticated ? (
          <>
            <Typography variant="h6" gutterBottom>
              Добро пожаловать в ваш персональный агрегатор новостей!
            </Typography>
            <Box sx={{ mt: 3 }}>
              <Button
                variant="contained"
                component={RouterLink}
                to="/news"
                size="large"
                sx={{ mr: 2 }}
              >
                Читать новости
              </Button>
              <Button
                variant="outlined"
                component={RouterLink}
                to="/subscriptions"
                size="large"
              >
                Управлять подписками
              </Button>
            </Box>
          </>
        ) : (
          <>
            <Typography variant="h6" gutterBottom>
              Начните использовать уже сегодня
            </Typography>
            <Box sx={{ mt: 3 }}>
              <Button
                variant="contained"
                component={RouterLink}
                to="/register"
                size="large"
                sx={{ mr: 2 }}
              >
                Зарегистрироваться
              </Button>
              <Button
                variant="outlined"
                component={RouterLink}
                to="/login"
                size="large"
              >
                Войти
              </Button>
            </Box>
          </>
        )}
      </Paper>
    </Box>
  );
}

export default HomePage;