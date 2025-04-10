package example

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-hexagonal/api/dto"
	"go-hexagonal/domain/model"
	"go-hexagonal/domain/repo"
)

// MockExampleService is defined in create_test.go

// TestablGetUseCase 为测试目的修改GetUseCase
type TestablGetUseCase struct {
	GetUseCase
	txProvider func(ctx context.Context) (repo.Transaction, error)
}

func NewTestablGetUseCase(svc *MockExampleService) *TestablGetUseCase {
	return &TestablGetUseCase{
		GetUseCase: GetUseCase{
			exampleService: svc,
		},
		txProvider: CreateTestTransaction,
	}
}

// Execute 重写Execute方法以替换事务处理逻辑
func (uc *TestablGetUseCase) Execute(ctx context.Context, id int) (*dto.GetExampleResponse, error) {
	// 使用测试事务
	tx, err := uc.txProvider(ctx)
	if err != nil {
		return nil, fmt.Errorf("创建事务失败: %w", err)
	}
	defer tx.Rollback()

	// 调用领域服务
	example, err := uc.exampleService.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取样例失败: %w", err)
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	// 将领域模型转换为DTO
	result := &dto.GetExampleResponse{
		Id:        example.Id,
		Name:      example.Name,
		Alias:     example.Alias,
		CreatedAt: example.CreatedAt,
		UpdatedAt: example.UpdatedAt,
	}

	return result, nil
}

// TestGetUseCase_Success 测试通过ID获取样例的成功情况
func TestGetUseCase_Success(t *testing.T) {
	// 创建mock服务
	mockService := new(MockExampleService)

	// 测试数据
	exampleId := 1

	now := time.Now()
	expectedExample := &model.Example{
		Id:        exampleId,
		Name:      "Test Example",
		Alias:     "test",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 设置mock行为
	mockService.On("Get", mock.Anything, exampleId).Return(expectedExample, nil)

	// 使用可测试版本创建用例
	useCase := NewTestablGetUseCase(mockService)

	// 执行用例
	ctx := context.Background()
	result, err := useCase.Execute(ctx, exampleId)

	// 断言结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedExample.Id, result.Id)
	assert.Equal(t, expectedExample.Name, result.Name)
	assert.Equal(t, expectedExample.Alias, result.Alias)
	assert.Equal(t, expectedExample.CreatedAt, result.CreatedAt)
	assert.Equal(t, expectedExample.UpdatedAt, result.UpdatedAt)

	mockService.AssertExpectations(t)
}

// TestGetUseCase_Error 测试通过ID获取样例时出错的情况
func TestGetUseCase_Error(t *testing.T) {
	// 创建mock服务
	mockService := new(MockExampleService)

	// 测试数据
	exampleId := 999 // 不存在的ID

	// 设置mock行为 - 模拟错误
	expectedError := assert.AnError
	mockService.On("Get", mock.Anything, exampleId).Return(nil, expectedError)

	// 使用可测试版本创建用例
	useCase := NewTestablGetUseCase(mockService)

	// 执行用例
	ctx := context.Background()
	result, err := useCase.Execute(ctx, exampleId)

	// 断言结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "获取样例失败")

	mockService.AssertExpectations(t)
}
