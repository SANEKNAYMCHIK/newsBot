import { useState, useEffect } from 'react';
import { useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
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
  Tabs,
  Tab,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import { newsApi } from '../api/newsApi';

function SubscriptionsPage() {
  const navigate = useNavigate();
  const [mySubscriptions, setMySubscriptions] = useState([]);
  const [allSources, setAllSources] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [tabValue, setTabValue] = useState(0);
  const [addSourceDialogOpen, setAddSourceDialogOpen] = useState(false);
  const [newSource, setNewSource] = useState({
    name: '',
    url: '',
    category_id: '',
  });
  const [categories, setCategories] = useState([]);
  const { user } = useSelector((state) => state.auth);

  // Функция для получения ID источника (обрабатывает оба варианта)
  const getSourceId = (source) => {
    return source.id || source.source_id;
  };

  // Функция для получения имени источника
  const getSourceName = (source) => {
    return source.name || source.title || source.source_name || `Источник ${getSourceId(source)}`;
  };

  // Загружаем подписки, источники и категории
  useEffect(() => {
    console.log('🔄 useEffect triggered for data loading');
    
    const loadData = async () => {
      if (!user) {
        console.log('⚠️ No user, skipping data load');
        setLoading(false);
        return;
      }
      
      setLoading(true);
      setError('');
      console.log('📡 Loading subscriptions, sources and categories...');
      
      try {
        // Параллельно загружаем подписки, источники и категории
        const [subscriptionsRes, sourcesRes, categoriesRes] = await Promise.all([
          newsApi.getUserSubscriptions(),
          newsApi.getSources(1, 100),
          newsApi.getCategories()
        ]);
        
        console.log('✅ Subscriptions response:', subscriptionsRes);
        console.log('📊 Subscriptions data structure:', subscriptionsRes.data);
        
        console.log('✅ Sources response:', sourcesRes);
        console.log('📊 Sources data structure:', sourcesRes.data);
        
        console.log('✅ Categories response:', categoriesRes);
        console.log('📊 Categories data:', categoriesRes.data);
        
        // Обработка подписок
        let subscriptionsArray = [];
        const subscriptionsData = subscriptionsRes.data;
        
        if (subscriptionsData && Array.isArray(subscriptionsData)) {
          subscriptionsArray = subscriptionsData;
        } else if (subscriptionsData && Array.isArray(subscriptionsData.data)) {
          subscriptionsArray = subscriptionsData.data;
        } else if (subscriptionsData && Array.isArray(subscriptionsData.items)) {
          subscriptionsArray = subscriptionsData.items;
        }
        
        console.log('📦 Subscriptions loaded:', subscriptionsArray.length);
        setMySubscriptions(subscriptionsArray);
        
        // Обработка всех источников
        let sourcesArray = [];
        const sourcesData = sourcesRes.data;
        
        if (sourcesData && Array.isArray(sourcesData)) {
          sourcesArray = sourcesData;
        } else if (sourcesData && Array.isArray(sourcesData.data)) {
          sourcesArray = sourcesData.data;
        } else if (sourcesData && Array.isArray(sourcesData.items)) {
          sourcesArray = sourcesData.items;
        }
        
        console.log('📦 Sources loaded:', sourcesArray.length);
        setAllSources(sourcesArray);
        
        // Обработка категорий
        let categoriesArray = [];
        const categoriesData = categoriesRes.data;
        
        if (categoriesData && Array.isArray(categoriesData)) {
          categoriesArray = categoriesData;
        } else if (categoriesData && Array.isArray(categoriesData.data)) {
          categoriesArray = categoriesData.data;
        } else if (categoriesData && Array.isArray(categoriesData.items)) {
          categoriesArray = categoriesData.items;
        }
        
        console.log('📦 Categories loaded:', categoriesArray.length);
        setCategories(categoriesArray);
        
        console.log('✅ Data loaded successfully');
        
      } catch (err) {
        console.error('❌ Load data error:', err);
        if (err.response) {
          console.error('📡 Response status:', err.response.status);
          console.error('📡 Response data:', err.response.data);
          
          if (err.response.status === 401) {
            setError('Ошибка авторизации. Пожалуйста, войдите снова.');
          } else {
            setError(`Ошибка: ${err.response.status}. ${err.response.data?.message || ''}`);
          }
        } else {
          setError('Не удалось загрузить данные. ' + err.message);
        }
      } finally {
        setLoading(false);
        console.log('✅ Loading finished');
      }
    };
    
    loadData();
  }, [user]);

  // Обработка подписки
  const handleSubscribe = async (source) => {
    const sourceId = getSourceId(source);
    console.log('➕ Subscribing to source:', sourceId, getSourceName(source));
    
    try {
      await newsApi.subscribe(sourceId);
      
      // Добавляем в подписки
      setMySubscriptions(prev => [...prev, source]);
      console.log('✅ Subscribed successfully');
    } catch (err) {
      console.error('❌ Subscribe error:', err);
      alert('Не удалось подписаться. Попробуйте снова.');
    }
  };

  // Обработка отписки
  const handleUnsubscribe = async (source) => {
    const sourceId = getSourceId(source);
    console.log('➖ Unsubscribing from source:', sourceId, getSourceName(source));
    
    try {
      await newsApi.unsubscribe(sourceId);
      
      // Удаляем из подписок
      setMySubscriptions(prev => prev.filter(s => getSourceId(s) !== sourceId));
      console.log('✅ Unsubscribed successfully');
    } catch (err) {
      console.error('❌ Unsubscribe error:', err);
      alert('Не удалось отписаться. Попробуйте снова.');
    }
  };

  // Проверка подписки
  const isSubscribed = (source) => {
    const sourceId = getSourceId(source);
    return mySubscriptions.some(s => getSourceId(s) === sourceId);
  };

  // Переход к новостям источника
  const handleReadNews = (source) => {
    const sourceId = getSourceId(source);
    navigate(`/news?source=${sourceId}`);
  };

  // Добавление нового источника
  const handleAddSource = async () => {
    if (!newSource.name || !newSource.url || !newSource.category_id) {
      alert('Пожалуйста, заполните все поля');
      return;
    }
    
    try {
      await newsApi.addSource(newSource);
      alert('Источник успешно добавлен!');
      setAddSourceDialogOpen(false);
      setNewSource({ name: '', url: '', category_id: '' });
      
      // Обновляем список источников
      const sourcesRes = await newsApi.getSources(1, 100);
      const sourcesData = sourcesRes.data;
      
      if (sourcesData && Array.isArray(sourcesData)) {
        setAllSources(sourcesData);
      } else if (sourcesData && Array.isArray(sourcesData.data)) {
        setAllSources(sourcesData.data);
      }
    } catch (err) {
      console.error('❌ Add source error:', err);
      alert(`Ошибка при добавлении источника: ${err.response?.data?.error || err.message}`);
    }
  };

  const handleTabChange = (event, newValue) => {
    console.log('🔄 Tab changed to:', newValue);
    setTabValue(newValue);
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
        <CircularProgress />
        <Typography sx={{ ml: 2 }}>Загрузка подписок...</Typography>
      </Box>
    );
  }

  return (
    <Container maxWidth="lg">
      <Box sx={{ mt: 4, mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Мои подписки
        </Typography>
        
        {/* Кнопка добавления источника */}
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Tabs value={tabValue} onChange={handleTabChange} sx={{ flexGrow: 1 }}>
            <Tab label={`Мои подписки (${mySubscriptions.length})`} />
            <Tab label={`Все источники (${allSources.length})`} />
          </Tabs>
          
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => setAddSourceDialogOpen(true)}
            sx={{ ml: 2 }}
          >
            Добавить источник
          </Button>
        </Box>

        {error && (
          <Alert severity="error" sx={{ mb: 3 }}>
            {error}
          </Alert>
        )}

        {tabValue === 0 ? (
          // Вкладка "Мои подписки"
          mySubscriptions.length === 0 ? (
            <Box textAlign="center" py={4}>
              <Typography variant="h6" color="text.secondary" gutterBottom>
                У вас еще нет подписок
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Добавьте источники во вкладке "Все источники"
              </Typography>
            </Box>
          ) : (
            <Grid container spacing={3}>
              {mySubscriptions.map((source, index) => (
                <Grid item xs={12} md={6} lg={4} key={getSourceId(source) || index}>
                  <Card variant="outlined" sx={{ height: '100%' }}>
                    <CardContent>
                      <Typography variant="h6" gutterBottom>
                        {getSourceName(source)}
                      </Typography>
                      <Typography variant="caption" color="text.secondary" display="block">
                        ID: {getSourceId(source)}
                      </Typography>
                    </CardContent>
                    <CardActions>
                      <Button 
                        size="small" 
                        color="error"
                        onClick={() => handleUnsubscribe(source)}
                      >
                        Отписаться
                      </Button>
                      <Button 
                        size="small" 
                        onClick={() => handleReadNews(source)}
                      >
                        Читать новости
                      </Button>
                    </CardActions>
                  </Card>
                </Grid>
              ))}
            </Grid>
          )
        ) : (
          // Вкладка "Все источники"
          allSources.length === 0 ? (
            <Box textAlign="center" py={4}>
              <Typography variant="h6" color="text.secondary" gutterBottom>
                Нет доступных источников
              </Typography>
            </Box>
          ) : (
            <Grid container spacing={3}>
              {allSources.map((source, index) => {
                const subscribed = isSubscribed(source);
                return (
                  <Grid item xs={12} md={6} lg={4} key={getSourceId(source) || index}>
                    <Card variant="outlined" sx={{ height: '100%' }}>
                      <CardContent>
                        <Typography variant="h6" gutterBottom>
                          {getSourceName(source)}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" paragraph>
                          {source.description || source.url || 'Описание отсутствует'}
                        </Typography>
                        <Typography variant="caption" display="block" 
                          color={subscribed ? "success.main" : "text.secondary"}>
                          Статус: {subscribed ? '✓ Подписан' : 'Не подписан'}
                        </Typography>
                      </CardContent>
                      <CardActions>
                        {subscribed ? (
                          <Button 
                            size="small" 
                            color="error"
                            onClick={() => handleUnsubscribe(source)}
                          >
                            Отписаться
                          </Button>
                        ) : (
                          <Button 
                            size="small" 
                            color="primary"
                            onClick={() => handleSubscribe(source)}
                          >
                            Подписаться
                          </Button>
                        )}
                        <Button 
                          size="small" 
                          onClick={() => handleReadNews(source)}
                        >
                          Новости
                        </Button>
                      </CardActions>
                    </Card>
                  </Grid>
                );
              })}
            </Grid>
          )
        )}
      </Box>

      {/* Диалог добавления источника */}
      <Dialog open={addSourceDialogOpen} onClose={() => setAddSourceDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Добавить новый источник RSS</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Название источника"
            fullWidth
            value={newSource.name}
            onChange={(e) => setNewSource({...newSource, name: e.target.value})}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="URL RSS-ленты"
            fullWidth
            value={newSource.url}
            onChange={(e) => setNewSource({...newSource, url: e.target.value})}
            sx={{ mb: 2 }}
            helperText="Пример: https://example.com/rss"
          />
          <FormControl fullWidth>
            <InputLabel>Категория</InputLabel>
            <Select
              value={newSource.category_id}
              label="Категория"
              onChange={(e) => setNewSource({...newSource, category_id: e.target.value})}
            >
              <MenuItem value="">Выберите категорию</MenuItem>
              {categories.map((category) => (
                <MenuItem key={category.id} value={category.id}>
                  {category.name || category.title}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setAddSourceDialogOpen(false)}>Отмена</Button>
          <Button onClick={handleAddSource} variant="contained">Добавить</Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
}

export default SubscriptionsPage;