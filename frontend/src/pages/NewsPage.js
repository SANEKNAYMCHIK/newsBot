import React, { useState, useEffect, useCallback } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useLocation, useNavigate } from 'react-router-dom';
import {
  Container,
  Card,
  CardContent,
  CardActions,
  Typography,
  Button,
  Box,
  CircularProgress,
  Alert,
  Grid,
  Pagination,
  Chip,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Breadcrumbs,
  Link,
  Tooltip,
  Snackbar,
} from '@mui/material';
import HomeIcon from '@mui/icons-material/Home';
import NewspaperIcon from '@mui/icons-material/Newspaper';
import RefreshIcon from '@mui/icons-material/Refresh';
import { newsApi } from '../api/newsApi';
import { setCredentials } from '../store/authSlice';

function NewsPage() {
  const dispatch = useDispatch();
  const location = useLocation();
  const navigate = useNavigate();
  const [news, setNews] = useState([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState('');
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalNews, setTotalNews] = useState(0);
  const [sources, setSources] = useState([]);
  const [selectedSource, setSelectedSource] = useState('');
  const [selectedSourceName, setSelectedSourceName] = useState('');
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });
  
  const { token } = useSelector((state) => state.auth);

  // Восстановление токена из localStorage
  useEffect(() => {
    const storedToken = localStorage.getItem('token');
    const storedUser = localStorage.getItem('user');
    
    if (!token && storedToken && storedUser) {
      dispatch(setCredentials({
        user: JSON.parse(storedUser),
        token: storedToken
      }));
    }
  }, [dispatch, token]);

  // Обработка query параметров из URL
  useEffect(() => {
    const queryParams = new URLSearchParams(location.search);
    const sourceFromQuery = queryParams.get('source');
    
    if (sourceFromQuery) {
      console.log('📌 Source from URL query:', sourceFromQuery);
      setSelectedSource(sourceFromQuery);
    }
  }, [location]);

  // Загружаем подписки
  useEffect(() => {
    const loadSubscriptions = async () => {
      const currentToken = localStorage.getItem('token');
      
      if (!currentToken) {
        console.log('⚠️ No token in localStorage');
        return;
      }
      
      try {
        const response = await newsApi.getUserSubscriptions();
        console.log('✅ Subscriptions for filter:', response.data);
        
        let subscriptionsArray = [];
        
        if (Array.isArray(response.data)) {
          subscriptionsArray = response.data;
        } else if (response.data && Array.isArray(response.data.data)) {
          subscriptionsArray = response.data.data;
        } else if (response.data && Array.isArray(response.data.items)) {
          subscriptionsArray = response.data.items;
        }
        
        setSources(subscriptionsArray);
        
        if (selectedSource && subscriptionsArray.length > 0) {
          const source = subscriptionsArray.find(s => 
            (s.id || s.source_id)?.toString() === selectedSource.toString()
          );
          if (source) {
            setSelectedSourceName(source.name || source.source_name || source.title);
          }
        }
        
      } catch (err) {
        console.error('❌ Failed to load subscriptions:', err);
      }
    };
    
    loadSubscriptions();
  }, [selectedSource]);

  // Загружаем новости (useCallback для стабильной ссылки)
  const fetchNews = useCallback(async () => {
    const currentToken = localStorage.getItem('token');
    if (!currentToken) {
      setError('Пользователь не авторизован. Пожалуйста, войдите снова.');
      setLoading(false);
      return;
    }
    
    setLoading(true);
    setError('');
    
    try {
      let response;
      if (selectedSource) {
        response = await newsApi.getNewsBySource(selectedSource, page, 20);
      } else {
        response = await newsApi.getUserNews(page, 20);
      }
      
      console.log('✅ News response:', response.data);
      
      let newsArray = [];
      let pages = 1;
      let total = 0;
      
      if (response.data && Array.isArray(response.data.data)) {
        newsArray = response.data.data;
        pages = response.data.total_pages || 1;
        total = response.data.total || 0;
      } else if (Array.isArray(response.data)) {
        newsArray = response.data;
        pages = 1;
        total = response.data.length;
      } else if (response.data && Array.isArray(response.data.items)) {
        newsArray = response.data.items;
        pages = response.data.total_pages || 1;
        total = response.data.total || 0;
      }
      
      setNews(newsArray);
      setTotalPages(pages);
      setTotalNews(total);
      
    } catch (err) {
      console.error('❌ News fetch error:', err);
      if (err.response) {
        if (err.response.status === 401) {
          setError('Ошибка авторизации. Пожалуйста, войдите снова.');
          localStorage.removeItem('token');
          localStorage.removeItem('user');
        } else {
          setError(`Ошибка: ${err.response.status}. ${err.response.data?.message || ''}`);
        }
      } else {
        setError('Не удалось загрузить новости. ' + err.message);
      }
    } finally {
      setLoading(false);
    }
  }, [selectedSource, page]);

  useEffect(() => {
    fetchNews();
  }, [fetchNews]); // Теперь используем только fetchNews

  // Обновление новостей вручную
  const handleRefreshNews = async () => {
    const currentToken = localStorage.getItem('token');
    if (!currentToken) {
      setError('Пользователь не авторизован');
      return;
    }
    
    setRefreshing(true);
    try {
      await newsApi.refreshNews();
      setSnackbar({
        open: true,
        message: 'Новости успешно обновлены!',
        severity: 'success'
      });
      
      setTimeout(() => {
        fetchNews();
      }, 2000);
      
    } catch (err) {
      console.error('❌ Refresh news error:', err);
      setSnackbar({
        open: true,
        message: err.response?.data?.error || 'Ошибка при обновлении новостей',
        severity: 'error'
      });
    } finally {
      setRefreshing(false);
    }
  };

  const handleSourceChange = (event) => {
    const sourceId = event.target.value;
    setSelectedSource(sourceId);
    setPage(1);
    
    const source = sources.find(s => 
      (s.id || s.source_id)?.toString() === sourceId.toString()
    );
    setSelectedSourceName(source?.name || source?.source_name || source?.title || '');
    
    if (sourceId) {
      navigate(`/news?source=${sourceId}`, { replace: true });
    } else {
      navigate('/news', { replace: true });
    }
  };

  const handlePageChange = (event, value) => {
    setPage(value);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const getSourceId = (source) => {
    return source.id || source.source_id;
  };

  const handleCloseSnackbar = () => {
    setSnackbar({ ...snackbar, open: false });
  };

  if (loading && page === 1) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh" flexDirection="column">
        <CircularProgress />
        <Typography sx={{ ml: 2, mt: 2 }}>Загрузка новостей...</Typography>
      </Box>
    );
  }

  return (
    <Container maxWidth="lg">
      <Box sx={{ mt: 4, mb: 4 }}>
        <Breadcrumbs sx={{ mb: 3 }}>
          <Link
            underline="hover"
            color="inherit"
            href="/"
            sx={{ display: 'flex', alignItems: 'center' }}
          >
            <HomeIcon sx={{ mr: 0.5 }} fontSize="inherit" />
            Главная
          </Link>
          <Typography color="text.primary" sx={{ display: 'flex', alignItems: 'center' }}>
            <NewspaperIcon sx={{ mr: 0.5 }} fontSize="inherit" />
            Новости
          </Typography>
        </Breadcrumbs>
        
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Box>
            <Typography variant="h4" component="h1" gutterBottom>
              Новости
              {selectedSourceName && (
                <Typography component="span" variant="h5" color="primary" sx={{ ml: 2 }}>
                  : {selectedSourceName}
                </Typography>
              )}
            </Typography>
          </Box>
          
          <Tooltip title="Обновить новости вручную">
            <Button
              variant="outlined"
              startIcon={<RefreshIcon />}
              onClick={handleRefreshNews}
              disabled={refreshing}
            >
              {refreshing ? 'Обновление...' : 'Обновить новости'}
            </Button>
          </Tooltip>
        </Box>
        
        <Box sx={{ mb: 3, p: 2, bgcolor: 'grey.50', borderRadius: 2, border: '1px solid', borderColor: 'divider' }}>
          <Grid container spacing={2} alignItems="center">
            <Grid item xs={12} md={6}>
              <FormControl fullWidth>
                <InputLabel id="source-select-label">Выберите источник</InputLabel>
                <Select
                  labelId="source-select-label"
                  value={selectedSource}
                  label="Выберите источник"
                  onChange={handleSourceChange}
                  sx={{ minWidth: 250 }}
                >
                  <MenuItem value="">Все подписанные источники</MenuItem>
                  {sources.map((source) => (
                    <MenuItem key={getSourceId(source)} value={getSourceId(source)}>
                      {source.name || source.source_name || source.title || `Источник ${getSourceId(source)}`}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={6}>
              <Typography variant="body2" color="text.secondary" align="right">
                Найдено новостей: {totalNews}
                {selectedSource && ` (страница ${page} из ${totalPages})`}
              </Typography>
            </Grid>
          </Grid>
        </Box>

        {error && (
          <Alert severity="error" sx={{ mb: 3 }}>
            {error}
            {error.includes('авторизации') && (
              <Box sx={{ mt: 1 }}>
                <Button 
                  variant="contained" 
                  size="small" 
                  href="/login"
                >
                  Войти
                </Button>
              </Box>
            )}
          </Alert>
        )}

        {news.length === 0 && !loading ? (
          <Box textAlign="center" py={4}>
            <Typography variant="h6" color="text.secondary" gutterBottom>
              Новостей пока нет
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {sources.length === 0 
                ? 'У вас нет подписок на источники. Добавьте источники на странице "Мои подписки".' 
                : 'Попробуйте выбрать другой источник или подождите обновления новостей.'}
            </Typography>
            {sources.length === 0 && (
              <Button 
                variant="contained" 
                sx={{ mt: 2 }}
                href="/subscriptions"
              >
                Перейти к подпискам
              </Button>
            )}
          </Box>
        ) : (
          <>

          <Grid container spacing={3} sx={{ 
  display: 'grid',
  gridTemplateColumns: {
    xs: '1fr',
    sm: 'repeat(2, 1fr)',
    md: 'repeat(3, 1fr)'
  },
  gridAutoRows: '1fr', // Это заставляет все карточки в ряду быть одинаковой высоты
  gap: 3
}}>
  {news.map((item, index) => (
    <Box 
      key={item.id || index}
      sx={{ 
        display: 'flex',
        height: '100%'
      }}
    >
      <Card 
        variant="outlined"
        sx={{ 
          width: '100%',
          display: 'flex',
          flexDirection: 'column',
          height: '100%', // Занимает всю высоту родителя
          '&:hover': { 
            boxShadow: 3,
            transform: 'translateY(-4px)',
            transition: 'all 0.3s ease'
          }
        }}
      >
        <CardContent sx={{ 
          flexGrow: 1,
          display: 'flex',
          flexDirection: 'column',
          overflow: 'hidden',
          pb: 1
        }}>
          {/* Заголовок с фиксированной высотой */}
          <Typography 
            variant="h6" 
            component="h2" 
            gutterBottom 
            sx={{ 
              lineHeight: 1.3,
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              display: '-webkit-box',
              WebkitLineClamp: 3,
              WebkitBoxOrient: 'vertical',
              minHeight: '4.5em',
              maxHeight: '4.5em',
              flexShrink: 0
            }}
          >
            {item.title || 'Без названия'}
          </Typography>
          
          {/* Описание - с очисткой пробелов и нормализацией */}
          <Box sx={{ 
            flexGrow: 1,
            overflow: 'hidden',
            mb: 1
          }}>
            <Typography 
              variant="body2" 
              color="text.secondary" 
              sx={{
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                display: '-webkit-box',
                WebkitLineClamp: 5,
                WebkitBoxOrient: 'vertical',
                lineHeight: 1.5
              }}
            >
              {/* Улучшенная очистка текста */}
              {(item.content || item.description || 'Нет описания')
                .replace(/<[^>]*>/g, '') // Удаляем HTML теги
                .replace(/&nbsp;/g, ' ') // Заменяем неразрывные пробелы
                .replace(/\s+/g, ' ') // Заменяем множественные пробелы на один
                .trim() // Убираем пробелы в начале и конце
                .substring(0, 250)}
              {(item.content || item.description)?.length > 250 ? '...' : ''}
            </Typography>
          </Box>
          
          {/* Нижняя часть - прижимаем к низу */}
          <Box sx={{ 
            mt: 'auto',
            display: 'flex', 
            justifyContent: 'space-between', 
            alignItems: 'center', 
            pt: 2,
            borderTop: '1px solid',
            borderColor: 'divider',
            flexShrink: 0
          }}>
            <Chip 
              label={item.source_name || 'Источник'} 
              size="small" 
              color="primary"
              variant="outlined"
              sx={{ maxWidth: '120px' }}
            />
            <Typography variant="caption" color="text.secondary" noWrap>
              {item.published_at 
                ? new Date(item.published_at).toLocaleDateString('ru-RU', {
                    day: 'numeric',
                    month: 'short',
                    year: 'numeric'
                  })
                : 'Дата неизвестна'}
            </Typography>
          </Box>
        </CardContent>
        <CardActions sx={{ pt: 0, pb: 2, px: 2 }}>
          <Button 
            size="small" 
            component="a" 
            href={item.url || '#'} 
            target="_blank" 
            rel="noopener noreferrer"
            disabled={!item.url}
            fullWidth
            variant="contained"
            sx={{ py: 1 }}
          >
            Читать полностью
          </Button>
        </CardActions>
      </Card>
    </Box>
  ))}
</Grid>

            {/* <Grid container spacing={3} sx={{ display: 'flex', flexWrap: 'wrap' }}>
              {news.map((item, index) => (
                <Grid 
                  item 
                  xs={12} 
                  md={6} 
                  lg={4} 
                  key={item.id || index}
                  sx={{ display: 'flex' }}
                >
                  <Card 
                    variant="outlined"
                    sx={{ 
                      flex: 1,
                      display: 'flex',
                      flexDirection: 'column',
                      '&:hover': { 
                        boxShadow: 3,
                        transform: 'translateY(-4px)',
                        transition: 'all 0.3s ease'
                      }
                    }}
                  >
                    <CardContent sx={{ 
                      flexGrow: 1,
                      display: 'flex',
                      flexDirection: 'column',
                      overflow: 'hidden'
                    }}>
                      <Typography variant="h6" component="h2" gutterBottom sx={{ 
                        lineHeight: 1.3,
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        display: '-webkit-box',
                        WebkitLineClamp: 3,
                        WebkitBoxOrient: 'vertical',
                        minHeight: '4.5em',
                        flexShrink: 0
                      }}>
                        {item.title || 'Без названия'}
                      </Typography>
                      <Box sx={{ flexGrow: 1, overflow: 'hidden' }}>
                        <Typography variant="body2" color="text.secondary" paragraph sx={{
                          overflow: 'hidden',
                          textOverflow: 'ellipsis',
                          display: '-webkit-box',
                          WebkitLineClamp: 4,
                          WebkitBoxOrient: 'vertical',
                        }}>
                          {(item.content || item.description || '')
                            .replace(/<[^>]*>/g, '')
                            .replace(/&nbsp;/g, ' ')
                            .substring(0, 200)}...
                        </Typography>
                      </Box>
                      <Box sx={{ 
                        mt: 'auto',
                        display: 'flex', 
                        justifyContent: 'space-between', 
                        alignItems: 'center', 
                        pt: 2,
                        borderTop: '1px solid',
                        borderColor: 'divider'
                      }}>
                        <Chip 
                          label={item.source_name || 'Источник'} 
                          size="small" 
                          color="primary"
                          variant="outlined"
                        />
                        <Typography variant="caption" color="text.secondary">
                          {item.published_at 
                            ? new Date(item.published_at).toLocaleDateString('ru-RU', {
                                day: 'numeric',
                                month: 'long',
                                year: 'numeric'
                              })
                            : 'Дата неизвестна'}
                        </Typography>
                      </Box>
                    </CardContent>
                    <CardActions sx={{ pt: 0 }}>
                      <Button 
                        size="small" 
                        component="a" 
                        href={item.url || '#'} 
                        target="_blank" 
                        rel="noopener noreferrer"
                        disabled={!item.url}
                        fullWidth
                        variant="contained"
                      >
                        Читать полностью
                      </Button>
                    </CardActions>
                  </Card>
                </Grid>
              ))}
            </Grid> */}

            {totalPages > 1 && (
              <Box display="flex" justifyContent="center" sx={{ mt: 4, mb: 4 }}>
                <Pagination 
                  count={totalPages} 
                  page={page} 
                  onChange={handlePageChange} 
                  color="primary"
                  size="large"
                  showFirstButton
                  showLastButton
                  siblingCount={1}
                  boundaryCount={1}
                />
              </Box>
            )}
          </>
        )}
      </Box>

      <Snackbar
        open={snackbar.open}
        autoHideDuration={6000}
        onClose={handleCloseSnackbar}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert onClose={handleCloseSnackbar} severity={snackbar.severity} sx={{ width: '100%' }}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Container>
  );
}

export default NewsPage;