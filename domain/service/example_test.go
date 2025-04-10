package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-hexagonal/domain/event"
	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
)

// 创建Mock存储库
type MockExampleRepo struct {
	mock.Mock
}

func (m *MockExampleRepo) Create(ctx context.Context, tr repo.Transaction, example *model.Example) (*model.Example, error) {
	// 为参数example设置Id，这样它就会有正确的Id用于生成事件
	if example.Id == 0 {
		example.Id = 1 // 设置一个默认ID
	}

	args := m.Called(ctx, tr, example)

	// 如果mock配置为返回一个Example，确保使用它作为返回值
	if e, ok := args.Get(0).(*model.Example); ok {
		return e, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockExampleRepo) Delete(ctx context.Context, tr repo.Transaction, id int) error {
	args := m.Called(ctx, tr, id)
	return args.Error(0)
}

func (m *MockExampleRepo) Update(ctx context.Context, tr repo.Transaction, example *model.Example) error {
	args := m.Called(ctx, tr, example)
	return args.Error(0)
}

func (m *MockExampleRepo) GetByID(ctx context.Context, tr repo.Transaction, id int) (*model.Example, error) {
	args := m.Called(ctx, tr, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

func (m *MockExampleRepo) FindByName(ctx context.Context, tr repo.Transaction, name string) (*model.Example, error) {
	args := m.Called(ctx, tr, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

// 创建Mock缓存存储库
type MockExampleCacheRepo struct {
	mock.Mock
}

func (m *MockExampleCacheRepo) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockExampleCacheRepo) GetByID(ctx context.Context, id int) (*model.Example, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

func (m *MockExampleCacheRepo) GetByName(ctx context.Context, name string) (*model.Example, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Example), args.Error(1)
}

func (m *MockExampleCacheRepo) Set(ctx context.Context, example *model.Example) error {
	args := m.Called(ctx, example)
	return args.Error(0)
}

func (m *MockExampleCacheRepo) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockExampleCacheRepo) Invalidate(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// 创建Mock事件总线
type MockEventBus struct {
	mock.Mock
}

// Publish 实现EventBus接口的Publish方法
func (m *MockEventBus) Publish(ctx context.Context, event event.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// Subscribe 实现EventBus接口的Subscribe方法
func (m *MockEventBus) Subscribe(handler event.EventHandler) {
	m.Called(handler)
}

// Unsubscribe 实现EventBus接口的Unsubscribe方法
func (m *MockEventBus) Unsubscribe(handler event.EventHandler) {
	m.Called(handler)
}

// 创建Mock事务对象
type MockTransaction struct {
	mock.Mock
}

func (m *MockTransaction) Begin() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransaction) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransaction) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransaction) Conn(ctx context.Context) interface{} {
	args := m.Called(ctx)
	return args.Get(0)
}

// 为测试定义一个辅助函数来设置EventBus
func withEventBus(service *ExampleService, bus event.EventBus) *ExampleService {
	service.EventBus = bus
	return service
}

// 测试ExampleService的创建
func TestNewExampleService(t *testing.T) {
	mockRepo := new(MockExampleRepo)
	mockCacheRepo := new(MockExampleCacheRepo)
	mockEventBus := new(MockEventBus)

	// 测试不带可选参数的创建
	service := NewExampleService(mockRepo, nil)
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.Repository)
	assert.Nil(t, service.CacheRepo)
	assert.Nil(t, service.EventBus)

	// 测试带缓存的创建
	service = NewExampleService(mockRepo, mockCacheRepo)
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.Repository)
	assert.Equal(t, mockCacheRepo, service.CacheRepo)
	assert.Nil(t, service.EventBus)

	// 测试设置事件总线
	service = NewExampleService(mockRepo, mockCacheRepo)
	service.EventBus = mockEventBus
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.Repository)
	assert.Equal(t, mockCacheRepo, service.CacheRepo)
	assert.Equal(t, mockEventBus, service.EventBus)

	// 测试使用辅助函数设置事件总线
	service = withEventBus(NewExampleService(mockRepo, mockCacheRepo), mockEventBus)
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.Repository)
	assert.Equal(t, mockCacheRepo, service.CacheRepo)
	assert.Equal(t, mockEventBus, service.EventBus)
}

// TestExampleService_Create 测试Create方法
func TestExampleService_Create(t *testing.T) {
	// 准备
	mockRepo := new(MockExampleRepo)
	mockCacheRepo := new(MockExampleCacheRepo)
	mockEventBus := new(MockEventBus)

	service := NewExampleService(mockRepo, mockCacheRepo)
	service.EventBus = mockEventBus

	// 创建一个正确的Example实例以便获取事件
	input, err := model.NewExample("Test", "test-alias")
	assert.NoError(t, err)
	input.Id = 1 // 确保ID已设置

	// 模拟依赖行为
	mockRepo.On("Create", mock.Anything, mock.Anything, mock.AnythingOfType("*model.Example")).Run(func(args mock.Arguments) {
		// 当创建被调用时，为创建好的对象添加事件并确保ID已经设置
		example := args.Get(2).(*model.Example)
		example.Id = 1
	}).Return(input, nil) // 返回带有事件的对象

	mockCacheRepo.On("Set", mock.Anything, mock.Anything).Return(nil)
	mockEventBus.On("Publish", mock.Anything, mock.AnythingOfType("event.ExampleCreatedEvent")).Return(nil)

	// 执行
	result, err := service.Create(context.Background(), "Test", "test-alias")

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test", result.Name)
	assert.Equal(t, "test-alias", result.Alias)

	// 验证所有模拟调用
	mockRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
}

// 测试Delete方法
func TestExampleService_Delete(t *testing.T) {
	mockRepo := new(MockExampleRepo)
	mockCacheRepo := new(MockExampleCacheRepo)
	mockEventBus := new(MockEventBus)

	// 创建服务实例
	service := NewExampleService(mockRepo, mockCacheRepo)
	service.EventBus = mockEventBus

	testCases := []struct {
		name        string
		setupMocks  func()
		exampleId   int
		wantErr     bool
		expectedErr error
	}{
		{
			name: "成功删除示例",
			setupMocks: func() {
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 1).Return(&model.Example{
					Id:    1,
					Name:  "Test Example",
					Alias: "test-alias",
				}, nil)
				mockRepo.On("Delete", mock.Anything, mock.Anything, 1).Return(nil)
				mockCacheRepo.On("Delete", mock.Anything, 1).Return(nil)
				mockEventBus.On("Publish", mock.Anything, mock.Anything).Return(nil)
			},
			exampleId: 1,
			wantErr:   false,
		},
		{
			name: "示例不存在",
			setupMocks: func() {
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 2).Return(nil, repo.ErrNotFound)
			},
			exampleId:   2,
			wantErr:     true,
			expectedErr: repo.ErrNotFound,
		},
		{
			name: "删除错误",
			setupMocks: func() {
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 3).Return(&model.Example{
					Id:    3,
					Name:  "Test Example",
					Alias: "test-alias",
				}, nil)
				mockRepo.On("Delete", mock.Anything, mock.Anything, 3).Return(errors.New("delete error"))
			},
			exampleId: 3,
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 设置模拟行为
			mockRepo.ExpectedCalls = nil
			mockCacheRepo.ExpectedCalls = nil
			mockEventBus.ExpectedCalls = nil
			tc.setupMocks()

			// 执行测试
			ctx := context.Background()
			err := service.Delete(ctx, tc.exampleId)

			// 验证结果
			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}

			// 验证Mock调用
			mockRepo.AssertExpectations(t)
			mockCacheRepo.AssertExpectations(t)
			mockEventBus.AssertExpectations(t)
		})
	}
}

// 测试Update方法
func TestExampleService_Update(t *testing.T) {
	mockRepo := new(MockExampleRepo)
	mockCacheRepo := new(MockExampleCacheRepo)
	mockEventBus := new(MockEventBus)

	// 创建服务实例
	service := NewExampleService(mockRepo, mockCacheRepo)
	service.EventBus = mockEventBus

	testCases := []struct {
		name        string
		setupMocks  func()
		exampleId   int
		newName     string
		newAlias    string
		wantErr     bool
		expectedErr error
	}{
		{
			name: "成功更新示例",
			setupMocks: func() {
				example := &model.Example{
					Id:    1,
					Name:  "Original Name",
					Alias: "original-alias",
				}
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 1).Return(example, nil)
				mockRepo.On("Update", mock.Anything, mock.Anything, mock.AnythingOfType("*model.Example")).Return(nil)
				mockCacheRepo.On("Set", mock.Anything, mock.AnythingOfType("*model.Example")).Return(nil)
				mockEventBus.On("Publish", mock.Anything, mock.Anything).Return(nil)
			},
			exampleId: 1,
			newName:   "Updated Name",
			newAlias:  "updated-alias",
			wantErr:   false,
		},
		{
			name: "示例不存在",
			setupMocks: func() {
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 2).Return(nil, repo.ErrNotFound)
			},
			exampleId:   2,
			newName:     "Updated Name",
			newAlias:    "updated-alias",
			wantErr:     true,
			expectedErr: repo.ErrNotFound,
		},
		{
			name: "更新验证失败",
			setupMocks: func() {
				example := &model.Example{
					Id:    3,
					Name:  "Original Name",
					Alias: "original-alias",
				}
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 3).Return(example, nil)
				// 不会调用Update，因为验证失败
			},
			exampleId:   3,
			newName:     "", // 空名称会导致验证错误
			newAlias:    "updated-alias",
			wantErr:     true,
			expectedErr: model.ErrEmptyExampleName,
		},
		{
			name: "更新存储错误",
			setupMocks: func() {
				example := &model.Example{
					Id:    4,
					Name:  "Original Name",
					Alias: "original-alias",
				}
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 4).Return(example, nil)
				mockRepo.On("Update", mock.Anything, mock.Anything, mock.AnythingOfType("*model.Example")).Return(errors.New("update error"))
			},
			exampleId: 4,
			newName:   "Updated Name",
			newAlias:  "updated-alias",
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 设置模拟行为
			mockRepo.ExpectedCalls = nil
			mockCacheRepo.ExpectedCalls = nil
			mockEventBus.ExpectedCalls = nil
			tc.setupMocks()

			// 执行测试
			ctx := context.Background()
			err := service.Update(ctx, tc.exampleId, tc.newName, tc.newAlias)

			// 验证结果
			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}

			// 验证Mock调用
			mockRepo.AssertExpectations(t)
			mockCacheRepo.AssertExpectations(t)
			mockEventBus.AssertExpectations(t)
		})
	}
}

// 测试Get方法
func TestExampleService_Get(t *testing.T) {
	mockRepo := new(MockExampleRepo)
	mockCacheRepo := new(MockExampleCacheRepo)

	// 创建服务实例
	service := NewExampleService(mockRepo, mockCacheRepo)

	testCases := []struct {
		name        string
		setupMocks  func()
		exampleId   int
		wantErr     bool
		expectedErr error
	}{
		{
			name: "从缓存获取示例",
			setupMocks: func() {
				// 模拟缓存命中
				mockCacheRepo.On("GetByID", mock.Anything, 1).Return(&model.Example{
					Id:    1,
					Name:  "Cached Example",
					Alias: "cached-alias",
				}, nil)
				// 存储库不应该被调用
			},
			exampleId: 1,
			wantErr:   false,
		},
		{
			name: "从存储库获取示例（缓存未命中）",
			setupMocks: func() {
				// 模拟缓存未命中
				mockCacheRepo.On("GetByID", mock.Anything, 2).Return(nil, errors.New("cache miss"))
				// 从存储库获取
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 2).Return(&model.Example{
					Id:    2,
					Name:  "Database Example",
					Alias: "db-alias",
				}, nil)
				// 更新缓存
				mockCacheRepo.On("Set", mock.Anything, mock.AnythingOfType("*model.Example")).Return(nil)
			},
			exampleId: 2,
			wantErr:   false,
		},
		{
			name: "缓存和存储库都未找到示例",
			setupMocks: func() {
				// 模拟缓存未命中
				mockCacheRepo.On("GetByID", mock.Anything, 3).Return(nil, errors.New("cache miss"))
				// 存储库也未找到
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 3).Return(nil, repo.ErrNotFound)
			},
			exampleId:   3,
			wantErr:     true,
			expectedErr: repo.ErrNotFound,
		},
		{
			name: "无缓存情况下从存储库获取",
			setupMocks: func() {
				// 创建不带缓存的服务
				service = NewExampleService(mockRepo, nil)
				// 只使用存储库
				mockRepo.On("GetByID", mock.Anything, mock.Anything, 4).Return(&model.Example{
					Id:    4,
					Name:  "Database Example",
					Alias: "db-alias",
				}, nil)
			},
			exampleId: 4,
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 设置模拟行为
			mockRepo.ExpectedCalls = nil
			mockCacheRepo.ExpectedCalls = nil
			tc.setupMocks()

			// 执行测试
			ctx := context.Background()
			result, err := service.Get(ctx, tc.exampleId)

			// 验证结果
			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.exampleId, result.Id)
			}

			// 验证Mock调用
			mockRepo.AssertExpectations(t)
			mockCacheRepo.AssertExpectations(t)
		})
	}
}

// 测试FindByName方法
func TestExampleService_FindByName(t *testing.T) {
	mockRepo := new(MockExampleRepo)
	mockCacheRepo := new(MockExampleCacheRepo)

	// 创建服务实例
	service := NewExampleService(mockRepo, mockCacheRepo)

	testCases := []struct {
		name        string
		setupMocks  func()
		searchName  string
		wantErr     bool
		expectedErr error
	}{
		{
			name: "从缓存获取示例",
			setupMocks: func() {
				// 模拟缓存命中
				mockCacheRepo.On("GetByName", mock.Anything, "test-name").Return(&model.Example{
					Id:    1,
					Name:  "test-name",
					Alias: "cached-alias",
				}, nil)
				// 存储库不应该被调用
			},
			searchName: "test-name",
			wantErr:    false,
		},
		{
			name: "从存储库获取示例（缓存未命中）",
			setupMocks: func() {
				// 模拟缓存未命中
				mockCacheRepo.On("GetByName", mock.Anything, "db-name").Return(nil, errors.New("cache miss"))
				// 从存储库获取
				mockRepo.On("FindByName", mock.Anything, mock.Anything, "db-name").Return(&model.Example{
					Id:    2,
					Name:  "db-name",
					Alias: "db-alias",
				}, nil)
				// 更新缓存
				mockCacheRepo.On("Set", mock.Anything, mock.AnythingOfType("*model.Example")).Return(nil)
			},
			searchName: "db-name",
			wantErr:    false,
		},
		{
			name: "缓存和存储库都未找到示例",
			setupMocks: func() {
				// 模拟缓存未命中
				mockCacheRepo.On("GetByName", mock.Anything, "missing-name").Return(nil, errors.New("cache miss"))
				// 存储库也未找到
				mockRepo.On("FindByName", mock.Anything, mock.Anything, "missing-name").Return(nil, repo.ErrNotFound)
			},
			searchName:  "missing-name",
			wantErr:     true,
			expectedErr: repo.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 设置模拟行为
			mockRepo.ExpectedCalls = nil
			mockCacheRepo.ExpectedCalls = nil
			tc.setupMocks()

			// 执行测试
			ctx := context.Background()
			result, err := service.FindByName(ctx, tc.searchName)

			// 验证结果
			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.searchName, result.Name)
			}

			// 验证Mock调用
			mockRepo.AssertExpectations(t)
			mockCacheRepo.AssertExpectations(t)
		})
	}
}
