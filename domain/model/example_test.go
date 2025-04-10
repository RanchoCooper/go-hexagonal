package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewExample(t *testing.T) {
	tests := []struct {
		name    string
		inName  string
		inAlias string
		wantErr bool
		errType error
	}{
		{
			name:    "should create a valid example",
			inName:  "Valid Name",
			inAlias: "valid-alias",
			wantErr: false,
		},
		{
			name:    "should fail with empty name",
			inName:  "",
			inAlias: "valid-alias",
			wantErr: true,
			errType: ErrEmptyExampleName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			example, err := NewExample(tt.inName, tt.inAlias)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
				assert.Nil(t, example)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, example)
				assert.Equal(t, tt.inName, example.Name)
				assert.Equal(t, tt.inAlias, example.Alias)
				assert.NotEmpty(t, example.CreatedAt)
				assert.NotEmpty(t, example.UpdatedAt)

				// 检查事件
				events := example.Events()
				assert.Len(t, events, 1)
				assert.Equal(t, "example.created", events[0].EventType())

				// 确保事件已被消费（清空）
				assert.Empty(t, example.Events())
			}
		})
	}
}

func TestExample_Validate(t *testing.T) {
	tests := []struct {
		name    string
		example *Example
		wantErr bool
		errType error
	}{
		{
			name: "valid example",
			example: &Example{
				Id:    1,
				Name:  "Valid Name",
				Alias: "valid-alias",
			},
			wantErr: false,
		},
		{
			name: "invalid id",
			example: &Example{
				Id:    -1,
				Name:  "Valid Name",
				Alias: "valid-alias",
			},
			wantErr: true,
			errType: ErrInvalidExampleID,
		},
		{
			name: "empty name",
			example: &Example{
				Id:    1,
				Name:  "",
				Alias: "valid-alias",
			},
			wantErr: true,
			errType: ErrEmptyExampleName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.example.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExample_Update(t *testing.T) {
	tests := []struct {
		name     string
		example  *Example
		newName  string
		newAlias string
		wantErr  bool
		errType  error
	}{
		{
			name: "valid update",
			example: &Example{
				Id:     1,
				Name:   "Original Name",
				Alias:  "original-alias",
				events: make([]DomainEvent, 0),
			},
			newName:  "Updated Name",
			newAlias: "updated-alias",
			wantErr:  false,
		},
		{
			name: "empty name",
			example: &Example{
				Id:     1,
				Name:   "Original Name",
				Alias:  "original-alias",
				events: make([]DomainEvent, 0),
			},
			newName:  "",
			newAlias: "updated-alias",
			wantErr:  true,
			errType:  ErrEmptyExampleName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 记录更新前的时间
			beforeUpdate := time.Now()
			time.Sleep(10 * time.Millisecond) // 确保时间戳有变化

			err := tt.example.Update(tt.newName, tt.newAlias)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
				// 验证值未变化
				assert.NotEqual(t, tt.newName, tt.example.Name)
				assert.NotEqual(t, tt.newAlias, tt.example.Alias)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newName, tt.example.Name)
				assert.Equal(t, tt.newAlias, tt.example.Alias)
				assert.True(t, tt.example.UpdatedAt.After(beforeUpdate), "UpdatedAt should be updated")

				// 检查事件
				events := tt.example.Events()
				assert.Len(t, events, 1)
				assert.Equal(t, "example.updated", events[0].EventType())

				// 确保事件已被消费（清空）
				assert.Empty(t, tt.example.Events())
			}
		})
	}
}

func TestExample_MarkDeleted(t *testing.T) {
	example := &Example{
		Id:     1,
		Name:   "Test Example",
		Alias:  "test-alias",
		events: make([]DomainEvent, 0),
	}

	example.MarkDeleted()

	// 验证事件
	events := example.Events()
	assert.Len(t, events, 1)
	assert.Equal(t, "example.deleted", events[0].EventType())

	// 验证事件已清空
	assert.Empty(t, example.Events())

	// 类型断言事件以验证详情
	deletedEvent, ok := events[0].(ExampleDeletedEvent)
	assert.True(t, ok)
	assert.Equal(t, example.Id, deletedEvent.ExampleID)
}

func TestExample_TableName(t *testing.T) {
	example := Example{}
	assert.Equal(t, "example", example.TableName())
}

func TestExampleEvents(t *testing.T) {
	// 测试创建事件
	t.Run("Created Event", func(t *testing.T) {
		example := &Example{
			Id:    1,
			Name:  "Test Example",
			Alias: "test-alias",
		}

		event := NewExampleCreatedEvent(example)
		assert.Equal(t, "example.created", event.EventType())
		assert.Equal(t, example.Id, event.ExampleID)
		assert.Equal(t, example.Name, event.Name)
		assert.Equal(t, example.Alias, event.Alias)
	})

	// 测试更新事件
	t.Run("Updated Event", func(t *testing.T) {
		example := &Example{
			Id:    1,
			Name:  "Test Example",
			Alias: "test-alias",
		}

		event := NewExampleUpdatedEvent(example)
		assert.Equal(t, "example.updated", event.EventType())
		assert.Equal(t, example.Id, event.ExampleID)
		assert.Equal(t, example.Name, event.Name)
		assert.Equal(t, example.Alias, event.Alias)
	})

	// 测试删除事件
	t.Run("Deleted Event", func(t *testing.T) {
		example := &Example{
			Id:    1,
			Name:  "Test Example",
			Alias: "test-alias",
		}

		event := NewExampleDeletedEvent(example)
		assert.Equal(t, "example.deleted", event.EventType())
		assert.Equal(t, example.Id, event.ExampleID)
	})
}

func TestExample_addEvent(t *testing.T) {
	example := &Example{
		Id:     1,
		Name:   "Test Example",
		Alias:  "test-alias",
		events: make([]DomainEvent, 0),
	}

	// 创建一个事件
	event := NewExampleCreatedEvent(example)

	// 手动调用addEvent
	example.addEvent(event)

	// 验证事件是否已添加
	assert.Len(t, example.events, 1)
	assert.Equal(t, event, example.events[0])

	// 再添加一个事件
	updateEvent := NewExampleUpdatedEvent(example)
	example.addEvent(updateEvent)

	// 验证两个事件都存在
	assert.Len(t, example.events, 2)
	assert.Equal(t, event, example.events[0])
	assert.Equal(t, updateEvent, example.events[1])
}

func TestDomainEvents(t *testing.T) {
	// 测试事件收集和清空
	example, err := NewExample("Test Example", "test-alias")
	require.NoError(t, err)

	// 初始应该有一个创建事件
	events := example.Events()
	assert.Len(t, events, 1)
	assert.Equal(t, "example.created", events[0].EventType())

	// 消费后应该被清空
	assert.Empty(t, example.Events())

	// 更新触发新事件
	err = example.Update("Updated Name", "updated-alias")
	require.NoError(t, err)

	// 验证更新事件
	events = example.Events()
	assert.Len(t, events, 1)
	assert.Equal(t, "example.updated", events[0].EventType())

	// 清空后
	assert.Empty(t, example.Events())

	// 标记删除
	example.MarkDeleted()

	// 验证删除事件
	events = example.Events()
	assert.Len(t, events, 1)
	assert.Equal(t, "example.deleted", events[0].EventType())
}
