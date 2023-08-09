package task_manager

type TaskManager struct {
}

func NewTaskManager() *TaskManager {

	return &TaskManager{}
}

func RegisterTaskHandler(handler TaskHandler) {

}

func SubmitTask(task Task) (err error) {
	return nil
}
