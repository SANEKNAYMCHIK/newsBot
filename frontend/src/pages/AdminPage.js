import { useState, useEffect, useCallback } from 'react';
import { useSelector } from 'react-redux';
import {
  Container,
  Paper,
  Typography,
  Box,
  Tabs,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Button,
  IconButton,
  Chip,
  TextField,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
  CircularProgress,
  Card,
  CardContent,
  Grid,
  Pagination,
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings';
import PersonRemoveIcon from '@mui/icons-material/PersonRemove';
import AddIcon from '@mui/icons-material/Add';
import { newsApi } from '../api/newsApi';

function AdminPage() {
  const [tabValue, setTabValue] = useState(0);
  
  // Состояние для пользователей с пагинацией
  const [users, setUsers] = useState([]);
  const [usersPage, setUsersPage] = useState(1);
  const [usersPageSize, setUsersPageSize] = useState(20);
  const [totalUsers, setTotalUsers] = useState(0);
  const [totalUsersPages, setTotalUsersPages] = useState(1);
  const [loadingUsers, setLoadingUsers] = useState(false);
  
  // Состояние для источников с пагинацией
  const [sources, setSources] = useState([]);
  const [sourcesPage, setSourcesPage] = useState(1);
  const [sourcesPageSize, setSourcesPageSize] = useState(20);
  const [totalSources, setTotalSources] = useState(0);
  const [totalSourcesPages, setTotalSourcesPages] = useState(1);
  const [loadingSources, setLoadingSources] = useState(false);
  
  // Состояние для категорий
  const [categories, setCategories] = useState([]);
  const [loadingCategories, setLoadingCategories] = useState(false);
  
  const [error, setError] = useState('');
  const [successMessage, setSuccessMessage] = useState('');
  
  // Диалоги
  const [editSourceDialog, setEditSourceDialog] = useState({ open: false, source: null });
  const [addCategoryDialog, setAddCategoryDialog] = useState(false);
  const [newCategory, setNewCategory] = useState({ name: '' });
  
  const { user: currentUser } = useSelector((state) => state.auth);

  // Функция загрузки общих количеств при загрузке компонента
  const loadInitialCounts = useCallback(async () => {
    if (!currentUser || currentUser.role !== 'admin') {
      return;
    }

    try {
      // Загружаем категории сразу (они без пагинации)
      const categoriesRes = await newsApi.getCategories();
      let categoriesData = [];
      
      if (categoriesRes.data && Array.isArray(categoriesRes.data.data)) {
        categoriesData = categoriesRes.data.data;
      } else if (Array.isArray(categoriesRes.data)) {
        categoriesData = categoriesRes.data;
      } else if (categoriesRes.data && Array.isArray(categoriesRes.data.items)) {
        categoriesData = categoriesRes.data.items;
      }
      
      setCategories(categoriesData);
      
      // Загружаем общее количество источников (только для подсчета)
      // Используем минимальный pageSize для быстрой загрузки
      const sourcesRes = await newsApi.getAllSources(1, 1);
      let totalSourcesCount = 0;
      
      if (sourcesRes.data && Array.isArray(sourcesRes.data.data)) {
        totalSourcesCount = sourcesRes.data.total || sourcesRes.data.count || 0;
      } else if (sourcesRes.data && Array.isArray(sourcesRes.data.items)) {
        totalSourcesCount = sourcesRes.data.total || sourcesRes.data.count || 0;
      } else if (sourcesRes.data && sourcesRes.data.total) {
        totalSourcesCount = sourcesRes.data.total;
      }
      
      setTotalSources(totalSourcesCount);
      
      // Загружаем общее количество пользователей (только для подсчета)
      const usersRes = await newsApi.getAllUsers(1, 1);
      let totalUsersCount = 0;
      
      if (usersRes.data && Array.isArray(usersRes.data.data)) {
        totalUsersCount = usersRes.data.total || usersRes.data.count || 0;
      } else if (usersRes.data && Array.isArray(usersRes.data.items)) {
        totalUsersCount = usersRes.data.total || usersRes.data.count || 0;
      } else if (usersRes.data && usersRes.data.total) {
        totalUsersCount = usersRes.data.total;
      }
      
      setTotalUsers(totalUsersCount);
      
    } catch (err) {
      console.error('❌ Initial counts load error:', err);
    }
  }, [currentUser]);

  // Функция загрузки пользователей с пагинацией
  const loadUsers = useCallback(async () => {
    if (!currentUser || currentUser.role !== 'admin') {
      return;
    }
    
    setLoadingUsers(true);
    setError('');
    
    try {
      const response = await newsApi.getAllUsers(usersPage, usersPageSize);
      console.log('✅ Users response:', response.data);
      
      let usersData = [];
      let total = 0;
      let totalPages = 1;
      
      if (response.data && Array.isArray(response.data.data)) {
        usersData = response.data.data;
        total = response.data.total || response.data.count || 0;
        totalPages = response.data.total_pages || Math.ceil(total / usersPageSize);
      } else if (Array.isArray(response.data)) {
        usersData = response.data;
        total = response.data.length;
        totalPages = 1;
      } else if (response.data && Array.isArray(response.data.items)) {
        usersData = response.data.items;
        total = response.data.total || response.data.count || 0;
        totalPages = response.data.total_pages || Math.ceil(total / usersPageSize);
      }
      
      setUsers(usersData);
      setTotalUsers(total);
      setTotalUsersPages(totalPages);
      
    } catch (err) {
      console.error('❌ Users load error:', err);
      setError(err.response?.data?.error || 'Ошибка загрузки пользователей');
    } finally {
      setLoadingUsers(false);
    }
  }, [currentUser, usersPage, usersPageSize]);

  // Функция загрузки источников с пагинацией
  const loadSources = useCallback(async () => {
    if (!currentUser || currentUser.role !== 'admin') {
      return;
    }
    
    setLoadingSources(true);
    setError('');
    
    try {
      const response = await newsApi.getAllSources(sourcesPage, sourcesPageSize);
      console.log('✅ Sources response:', response.data);
      
      let sourcesData = [];
      let total = 0;
      let totalPages = 1;
      
      if (response.data && Array.isArray(response.data.data)) {
        sourcesData = response.data.data;
        total = response.data.total || response.data.count || 0;
        totalPages = response.data.total_pages || Math.ceil(total / sourcesPageSize);
      } else if (Array.isArray(response.data)) {
        sourcesData = response.data;
        total = response.data.length;
        totalPages = 1;
      } else if (response.data && Array.isArray(response.data.items)) {
        sourcesData = response.data.items;
        total = response.data.total || response.data.count || 0;
        totalPages = response.data.total_pages || Math.ceil(total / sourcesPageSize);
      }
      
      setSources(sourcesData);
      setTotalSources(total);
      setTotalSourcesPages(totalPages);
      
    } catch (err) {
      console.error('❌ Sources load error:', err);
      setError(err.response?.data?.error || 'Ошибка загрузки источников');
    } finally {
      setLoadingSources(false);
    }
  }, [currentUser, sourcesPage, sourcesPageSize]);

  // Функция загрузки категорий
  const loadCategories = useCallback(async () => {
    if (!currentUser || currentUser.role !== 'admin') {
      return;
    }
    
    setLoadingCategories(true);
    
    try {
      const response = await newsApi.getCategories();
      console.log('✅ Categories response:', response.data);
      
      let categoriesData = [];
      
      if (response.data && Array.isArray(response.data.data)) {
        categoriesData = response.data.data;
      } else if (Array.isArray(response.data)) {
        categoriesData = response.data;
      } else if (response.data && Array.isArray(response.data.items)) {
        categoriesData = response.data.items;
      }
      
      setCategories(categoriesData);
      
    } catch (err) {
      console.error('❌ Categories load error:', err);
      setError('Ошибка загрузки категорий');
    } finally {
      setLoadingCategories(false);
    }
  }, [currentUser]);

  // Инициализация при загрузке компонента
  useEffect(() => {
    if (!currentUser || currentUser.role !== 'admin') return;
    
    // Загружаем общие количества и данные для активной вкладки
    loadInitialCounts();
    
    switch (tabValue) {
      case 0: // Пользователи
        loadUsers();
        break;
      case 1: // Источники
        loadSources();
        break;
      case 2: // Категории
        // Категории уже загружены в loadInitialCounts
        break;
    }
  }, [currentUser, tabValue, loadInitialCounts, loadUsers, loadSources]);

  // Загрузка данных при изменении пагинации пользователей
  useEffect(() => {
    if (tabValue === 0) {
      loadUsers();
    }
  }, [usersPage, usersPageSize, tabValue, loadUsers]);

  // Загрузка данных при изменении пагинации источников
  useEffect(() => {
    if (tabValue === 1) {
      loadSources();
    }
  }, [sourcesPage, sourcesPageSize, tabValue, loadSources]);

  // Реальные обработчики с вызовом API
  const handleMakeAdmin = async (userId) => {
    try {
      await newsApi.makeAdmin(userId);
      setSuccessMessage('Пользователь назначен администратором');
      loadUsers(); // Перезагружаем список
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка назначения администратора');
    }
  };

  const handleRemoveAdmin = async (userId) => {
    try {
      await newsApi.removeAdmin(userId);
      setSuccessMessage('Администраторские права сняты');
      loadUsers(); // Перезагружаем список
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка снятия прав администратора');
    }
  };

  const handleUpdateSource = async () => {
    try {
      await newsApi.updateSource(editSourceDialog.source.id, editSourceDialog.source);
      setSuccessMessage('Источник успешно обновлен');
      setEditSourceDialog({ open: false, source: null });
      loadSources(); // Перезагружаем список
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка обновления источника');
    }
  };

  const handleDeleteSource = async (sourceId) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот источник?')) return;
    
    try {
      await newsApi.deleteSource(sourceId);
      setSuccessMessage('Источник успешно удален');
      loadSources(); // Перезагружаем список
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка удаления источника');
    }
  };

  const handleAddCategory = async () => {
    try {
      await newsApi.addCategory(newCategory);
      setSuccessMessage('Категория успешно добавлена');
      setAddCategoryDialog(false);
      setNewCategory({ name: '' });
      loadCategories(); // Перезагружаем список
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка добавления категории');
    }
  };

  const handleTabChange = (event, newValue) => {
    setTabValue(newValue);
  };

  const handleUsersPageChange = (event, value) => {
    setUsersPage(value);
  };

  const handleSourcesPageChange = (event, value) => {
    setSourcesPage(value);
  };

  const handleUsersPageSizeChange = (event) => {
    setUsersPageSize(event.target.value);
    setUsersPage(1); // Сбрасываем на первую страницу
  };

  const handleSourcesPageSizeChange = (event) => {
    setSourcesPageSize(event.target.value);
    setSourcesPage(1); // Сбрасываем на первую страницу
  };

  // Проверка прав доступа
  if (!currentUser || currentUser.role !== 'admin') {
    return (
      <Container maxWidth="lg">
        <Box sx={{ mt: 8, textAlign: 'center' }}>
          <Alert severity="error" sx={{ maxWidth: 600, margin: '0 auto' }}>
            <Typography variant="h6">Доступ запрещен</Typography>
            <Typography>Эта страница доступна только администраторам</Typography>
          </Alert>
        </Box>
      </Container>
    );
  }

  // Определяем общее состояние загрузки для активной вкладки
  const getLoadingState = () => {
    switch (tabValue) {
      case 0: return loadingUsers;
      case 1: return loadingSources;
      case 2: return loadingCategories;
      default: return false;
    }
  };

  const isLoading = getLoadingState();

  return (
    <Container maxWidth="lg">
      <Box sx={{ mt: 4, mb: 4 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
          <AdminPanelSettingsIcon sx={{ fontSize: 40, color: 'primary.main', mr: 2 }} />
          <Typography variant="h4" component="h1">
            Административная панель
          </Typography>
        </Box>

        {error && (
          <Alert severity="error" sx={{ mb: 3 }} onClose={() => setError('')}>
            {error}
          </Alert>
        )}

        {successMessage && (
          <Alert severity="success" sx={{ mb: 3 }} onClose={() => setSuccessMessage('')}>
            {successMessage}
          </Alert>
        )}

        {/* Статистика */}
        <Grid container spacing={3} sx={{ mb: 4 }}>
          <Grid item xs={12} md={4}>
            <Card>
              <CardContent>
                <Typography color="text.secondary" gutterBottom>
                  Пользователей
                </Typography>
                <Typography variant="h4" component="div">
                  {totalUsers}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} md={4}>
            <Card>
              <CardContent>
                <Typography color="text.secondary" gutterBottom>
                  Источников
                </Typography>
                <Typography variant="h4" component="div">
                  {totalSources}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} md={4}>
            <Card>
              <CardContent>
                <Typography color="text.secondary" gutterBottom>
                  Категорий
                </Typography>
                <Typography variant="h4" component="div">
                  {categories.length}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>

        <Tabs value={tabValue} onChange={handleTabChange} sx={{ mb: 3 }}>
          <Tab label={`Пользователи (${totalUsers})`} />
          <Tab label={`Источники (${totalSources})`} />
          <Tab label={`Категории (${categories.length})`} />
        </Tabs>

        {isLoading ? (
          <Box display="flex" justifyContent="center" alignItems="center" minHeight="40vh">
            <CircularProgress />
            <Typography sx={{ ml: 2 }}>Загрузка данных...</Typography>
          </Box>
        ) : (
          <>
            {tabValue === 0 && (
              <Box>
                {/* Панель управления пагинацией для пользователей */}
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                  <Typography variant="h6">Пользователи</Typography>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                    <Typography variant="body2" color="text.secondary">
                      Страница {usersPage} из {totalUsersPages}
                    </Typography>
                    <FormControl size="small" sx={{ minWidth: 100 }}>
                      <InputLabel>На странице</InputLabel>
                      <Select
                        value={usersPageSize}
                        label="На странице"
                        onChange={handleUsersPageSizeChange}
                      >
                        <MenuItem value={10}>10</MenuItem>
                        <MenuItem value={20}>20</MenuItem>
                        <MenuItem value={50}>50</MenuItem>
                        <MenuItem value={100}>100</MenuItem>
                      </Select>
                    </FormControl>
                  </Box>
                </Box>

                <TableContainer component={Paper}>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>ID</TableCell>
                        <TableCell>Email</TableCell>
                        <TableCell>Имя</TableCell>
                        <TableCell>Telegram Username</TableCell>
                        <TableCell>Роль</TableCell>
                        <TableCell>Действия</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {users.map((userItem) => (
                        <TableRow key={userItem.id}>
                          <TableCell>{userItem.id}</TableCell>
                          <TableCell>{userItem.email || '-'}</TableCell>
                          <TableCell>{userItem.tg_first_name || '-'}</TableCell>
                          <TableCell>
                            {userItem.tg_username ? (
                              <Chip 
                                label={`@${userItem.tg_username}`}
                                size="small"
                                color="primary"
                                variant="outlined"
                              />
                            ) : (
                              <Typography variant="body2" color="text.secondary">
                                Не указан
                              </Typography>
                            )}
                          </TableCell>
                          <TableCell>
                            <Chip 
                              label={userItem.role === 'admin' ? 'Админ' : 'Пользователь'} 
                              color={userItem.role === 'admin' ? 'success' : 'default'}
                              size="small"
                            />
                          </TableCell>
                          <TableCell>
                            {userItem.role !== 'admin' ? (
                              <Button
                                size="small"
                                startIcon={<AdminPanelSettingsIcon />}
                                onClick={() => handleMakeAdmin(userItem.id)}
                                color="success"
                              >
                                Сделать админом
                              </Button>
                            ) : userItem.id !== currentUser?.id ? (
                              <Button
                                size="small"
                                startIcon={<PersonRemoveIcon />}
                                onClick={() => handleRemoveAdmin(userItem.id)}
                                color="error"
                              >
                                Убрать админа
                              </Button>
                            ) : (
                              <Typography variant="caption" color="text.secondary">
                                Текущий пользователь
                              </Typography>
                            )}
                          </TableCell>
                        </TableRow>
                      ))}
                      {users.length === 0 && (
                        <TableRow>
                          <TableCell colSpan={6} align="center" sx={{ py: 3 }}>
                            <Typography variant="body2" color="text.secondary">
                              Нет данных о пользователях
                            </Typography>
                          </TableCell>
                        </TableRow>
                      )}
                    </TableBody>
                  </Table>
                </TableContainer>

                {/* Пагинация для пользователей */}
                {totalUsersPages > 1 && (
                  <Box display="flex" justifyContent="center" sx={{ mt: 3 }}>
                    <Pagination 
                      count={totalUsersPages} 
                      page={usersPage} 
                      onChange={handleUsersPageChange} 
                      color="primary"
                      showFirstButton
                      showLastButton
                    />
                  </Box>
                )}
              </Box>
            )}

            {tabValue === 1 && (
              <Box>
                {/* Панель управления пагинацией для источников */}
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                  <Typography variant="h6">Источники</Typography>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                    <Typography variant="body2" color="text.secondary">
                      Страница {sourcesPage} из {totalSourcesPages}
                    </Typography>
                    <FormControl size="small" sx={{ minWidth: 100 }}>
                      <InputLabel>На странице</InputLabel>
                      <Select
                        value={sourcesPageSize}
                        label="На странице"
                        onChange={handleSourcesPageSizeChange}
                      >
                        <MenuItem value={10}>10</MenuItem>
                        <MenuItem value={20}>20</MenuItem>
                        <MenuItem value={50}>50</MenuItem>
                        <MenuItem value={100}>100</MenuItem>
                      </Select>
                    </FormControl>
                  </Box>
                </Box>

                <TableContainer component={Paper} sx={{ mb: 3 }}>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>ID</TableCell>
                        <TableCell>Название</TableCell>
                        <TableCell>URL</TableCell>
                        <TableCell>Категория</TableCell>
                        <TableCell>Статус</TableCell>
                        <TableCell>Действия</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {sources.map((source) => (
                        <TableRow key={source.id}>
                          <TableCell>{source.id}</TableCell>
                          <TableCell>{source.name || source.title}</TableCell>
                          <TableCell>
                            <a href={source.url} target="_blank" rel="noopener noreferrer" style={{ textDecoration: 'none' }}>
                              {source.url?.substring(0, 30)}...
                            </a>
                          </TableCell>
                          <TableCell>{source.category_name || source.category_id}</TableCell>
                          <TableCell>
                            <Chip 
                              label={source.is_active ? 'Активен' : 'Неактивен'} 
                              color={source.is_active ? 'success' : 'error'}
                              size="small"
                            />
                          </TableCell>
                          <TableCell>
                            <IconButton
                              size="small"
                              onClick={() => setEditSourceDialog({ open: true, source })}
                              color="primary"
                            >
                              <EditIcon />
                            </IconButton>
                            <IconButton
                              size="small"
                              onClick={() => handleDeleteSource(source.id)}
                              color="error"
                            >
                              <DeleteIcon />
                            </IconButton>
                          </TableCell>
                        </TableRow>
                      ))}
                      {sources.length === 0 && (
                        <TableRow>
                          <TableCell colSpan={6} align="center" sx={{ py: 3 }}>
                            <Typography variant="body2" color="text.secondary">
                              Нет данных об источниках
                            </Typography>
                          </TableCell>
                        </TableRow>
                      )}
                    </TableBody>
                  </Table>
                </TableContainer>

                {/* Пагинация для источников */}
                {totalSourcesPages > 1 && (
                  <Box display="flex" justifyContent="center" sx={{ mt: 3 }}>
                    <Pagination 
                      count={totalSourcesPages} 
                      page={sourcesPage} 
                      onChange={handleSourcesPageChange} 
                      color="primary"
                      showFirstButton
                      showLastButton
                    />
                  </Box>
                )}
              </Box>
            )}

            {tabValue === 2 && (
              <Box>
                <Box sx={{ display: 'flex', justifyContent: 'flex-end', mb: 2 }}>
                  <Button
                    variant="contained"
                    startIcon={<AddIcon />}
                    onClick={() => setAddCategoryDialog(true)}
                  >
                    Добавить категорию
                  </Button>
                </Box>
                
                <TableContainer component={Paper}>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>ID</TableCell>
                        <TableCell>Название</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {categories.map((category) => (
                        <TableRow key={category.id}>
                          <TableCell>{category.id}</TableCell>
                          <TableCell>{category.name || category.title}</TableCell>
                        </TableRow>
                      ))}
                      {categories.length === 0 && (
                        <TableRow>
                          <TableCell colSpan={2} align="center" sx={{ py: 3 }}>
                            <Typography variant="body2" color="text.secondary">
                              Нет данных о категориях
                            </Typography>
                          </TableCell>
                        </TableRow>
                      )}
                    </TableBody>
                  </Table>
                </TableContainer>
              </Box>
            )}
          </>
        )}
      </Box>

      {/* Диалог редактирования источника */}
      <Dialog open={editSourceDialog.open} onClose={() => setEditSourceDialog({ open: false, source: null })}>
        <DialogTitle>Редактировать источник</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Название"
            fullWidth
            value={editSourceDialog.source?.name || ''}
            onChange={(e) => setEditSourceDialog({
              ...editSourceDialog,
              source: { ...editSourceDialog.source, name: e.target.value }
            })}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="URL"
            fullWidth
            value={editSourceDialog.source?.url || ''}
            onChange={(e) => setEditSourceDialog({
              ...editSourceDialog,
              source: { ...editSourceDialog.source, url: e.target.value }
            })}
            sx={{ mb: 2 }}
          />
          <FormControl fullWidth>
            <InputLabel>Статус</InputLabel>
            <Select
              value={editSourceDialog.source?.is_active ? 'active' : 'inactive'}
              label="Статус"
              onChange={(e) => setEditSourceDialog({
                ...editSourceDialog,
                source: { ...editSourceDialog.source, is_active: e.target.value === 'active' }
              })}
            >
              <MenuItem value="active">Активен</MenuItem>
              <MenuItem value="inactive">Неактивен</MenuItem>
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditSourceDialog({ open: false, source: null })}>Отмена</Button>
          <Button onClick={handleUpdateSource} variant="contained">Сохранить</Button>
        </DialogActions>
      </Dialog>

      {/* Диалог добавления категории */}
      <Dialog open={addCategoryDialog} onClose={() => setAddCategoryDialog(false)}>
        <DialogTitle>Добавить категорию</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Название категории"
            fullWidth
            value={newCategory.name}
            onChange={(e) => setNewCategory({ name: e.target.value })}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setAddCategoryDialog(false)}>Отмена</Button>
          <Button onClick={handleAddCategory} variant="contained">Добавить</Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
}

export default AdminPage;