package services

import (
	"sim_data_gen/entity"
	"sim_data_gen/services/snap"
	. "sim_data_gen/utils"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
)

var SnapSvr *SnapService

// SnapService 管理快照
type SnapService struct {
	snapDir string
	mu      sync.RWMutex
	tasks   map[string]*snap.SnapTask
}

func StartSnapService(varNamePrx string) {
	SnapSvr = new(SnapService)
	SnapSvr.snapDir = GetConfig().String(varNamePrx+".dir", "../snaps")
}

// StartTask 开始一个记录快照的任务
func (t *SnapService) StartTask(taskName string) (*snap.SnapTask, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	//if _, err := os.Stat(t.getTaskMainFilePath(taskName)); err == nil {
	//	return t.loadTaskFromFile(taskName)
	//}
	task := &snap.SnapTask{
		TaskName:   taskName,
		SnapData:   []string{},
		DataLoaded: true, //没啥数据要载入
	}
	t.makeSureTasksOk()
	t.tasks[taskName] = task
	if err := t.saveTask(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %v", err)
	}
	return task, nil
}

func (t *SnapService) makeSureTasksOk() {
	if t.tasks == nil {
		t.tasks = make(map[string]*snap.SnapTask)
	}
}
func (t *SnapService) GetTask(taskName string) (*snap.SnapTask, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if task, ok := t.tasks[taskName]; ok {
		if task.DataLoaded {
			return task, nil
		}
		return t.loadTaskFromFile(taskName)
	} else {
		return nil, fmt.Errorf("Task %s not found", taskName)
	}
}
func (t *SnapService) loadTaskFromFile(taskName string) (*snap.SnapTask, error) {
	taskPath := t.getTaskMainFilePath(taskName)
	if _, err := os.Stat(taskPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("task %s not found", taskName)
	}

	task, err := t.loadTask(taskName)
	if err == nil {
		return task, nil
	}
	return nil, err
}

// EndTask 结束一个记录快照的任务并保存
func (s *SnapService) EndTask(taskName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	/*
		if task, err := s.loadTask(taskName); err == nil {
			fmt.Printf("Task %s ended and saved with %d snapshots.\n", taskName, len(task.SnapData))
			// 保留历史数据文件，实际可根据需求修改
		} else {
			fmt.Printf("Task %s not found: %v\n", taskName, err)
		}*/
}

// DeleteTask 删除一个快照任务
func (s *SnapService) DeleteTask(taskName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	taskPath := s.getTaskMainFilePath(taskName)
	parentDir := filepath.Dir(taskPath)
	if err := os.Remove(parentDir); err == nil {
		Logger.Info("Task is deleted.\n", zap.String("taskName", taskName))
	} else {
		Logger.Info("Delete task failed",
			zap.String("taskName", taskName))
	}
}

// SnapToString 将快照转换为字符串
func (s *SnapService) SnapToString(taskName string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if task, err := s.loadTask(taskName); err == nil {
		data, err := json.Marshal(task.SnapData)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return "", fmt.Errorf("task %s not found", taskName)
}

// 新增辅助方法
func (t *SnapService) getTaskMainFilePath(taskName string) string {
	return filepath.Join(t.snapDir, taskName, "main.json")
}
func (t *SnapService) getTaskResultFilePath(taskName string) string {
	return filepath.Join(t.snapDir, taskName, "result.json")
}
func (t *SnapService) getTaskDir(taskName string) string {
	return filepath.Join(t.snapDir, taskName)
}

func (t *SnapService) saveTask(task *snap.SnapTask) error {
	// 确保目录存在
	taskDir := t.getTaskDir(task.TaskName)
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		return err
	}
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return os.WriteFile(t.getTaskMainFilePath(task.TaskName), data, 0644)
}

func (t *SnapService) loadTask(taskName string) (*snap.SnapTask, error) {
	data, err := os.ReadFile(t.getTaskMainFilePath(taskName))
	if err != nil {
		return nil, err
	}

	var task snap.SnapTask
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

/**
 * 生成报告和Excel
 * latestValidContent - 最近一次的准确的品味
 */
func (t *SnapService) GenerateMineReport(taskName string, latestValidContent map[string][2]float64) (string, error) {
	task, _ := SnapSvr.GetTask(taskName)
	if task == nil {
		return "", fmt.Errorf("task %s not found", taskName)
	}
	resultDir := t.getTaskResultFilePath(task.TaskName)
	data, err := os.ReadFile(resultDir)
	if err != nil {
		Logger.Info("Read resultList failed",
			zap.String("taskName", taskName), zap.Error(err))
		return "", err
	}
	results := new(entity.ContentResultList)
	err = json.Unmarshal(data, results)
	if err != nil {
		Logger.Info("Read resultList failed",
			zap.String("taskName", taskName), zap.Error(err))
		return "", err
	}

	taskDir := t.getTaskDir(task.TaskName)
	task.GenereteReport(results, path.Join(taskDir, "report.txt"))
	return task.GenerateExcel(results, latestValidContent, path.Join(taskDir, "report.xlsx"))
}

func (t *SnapService) SaveResult(taskName string, resultList *entity.ContentResultList) {
	task, _ := SnapSvr.GetTask(taskName)
	if task != nil {
		taskDir := t.getTaskResultFilePath(task.TaskName)
		data, err := json.Marshal(resultList)
		if err != nil {
			Logger.Info("Save resultList failed",
				zap.String("taskName", taskName), zap.Error(err))
		}
		err = os.WriteFile(taskDir, data, 0644)
		if err != nil {
			Logger.Info("Write task failed",
				zap.String("taskName", taskName), zap.Error(err))
		}
	}
}
