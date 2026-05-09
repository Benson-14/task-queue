package worker

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Benson-14/task-queue/internal/queue"
	"github.com/Benson-14/task-queue/internal/task"
)

type processorFn func(id int, payload json.RawMessage) error

var processors = map[string]processorFn{
	"send_email":      processSendEmail,
	"resize_image":    processResizeImage,
	"generate_report": processGenerateReport,
	"db_cleanup":      processDBCleanup,
	"batch_job":       processBatchJob,
}

// ── handler implementations ──────────────────────────────────────────────────

func processSendEmail(workerID int, payload json.RawMessage) error {
	var p struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return fmt.Errorf("bad payload: %w", err)
	}
	time.Sleep(500 * time.Millisecond) // simulate SMTP round-trip
	fmt.Printf("[worker-%d] ✉  Sent email to %q — subject: %q\n", workerID, p.To, p.Subject)
	return nil
}

func processResizeImage(workerID int, payload json.RawMessage) error {
	var p struct {
		ImageURL string `json:"image_url"`
		Width    int    `json:"width"`
		Height   int    `json:"height"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return fmt.Errorf("bad payload: %w", err)
	}
	time.Sleep(1500 * time.Millisecond) // simulate heavy image processing
	fmt.Printf("[worker-%d] 🖼  Resized %q → %dx%d\n", workerID, p.ImageURL, p.Width, p.Height)
	return nil
}

func processGenerateReport(workerID int, payload json.RawMessage) error {
	var p struct {
		UserID int    `json:"user_id"`
		Format string `json:"format"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return fmt.Errorf("bad payload: %w", err)
	}
	time.Sleep(3 * time.Second) // simulate slow PDF generation
	fmt.Printf("[worker-%d] 📊  Generated %s report for user %d\n", workerID, p.Format, p.UserID)
	return nil
}

func processDBCleanup(workerID int, payload json.RawMessage) error {
	var p struct {
		Table         string `json:"table"`
		OlderThanDays int    `json:"older_than_days"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return fmt.Errorf("bad payload: %w", err)
	}
	time.Sleep(200 * time.Millisecond)
	return fmt.Errorf("permission denied: insufficient privileges to truncate table %q", p.Table)
}

func processBatchJob(workerID int, payload json.RawMessage) error {
	var p struct {
		JobID int `json:"job_id"`
		Items int `json:"items"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return fmt.Errorf("bad payload: %w", err)
	}
	delay := time.Duration(p.Items) * time.Millisecond
	time.Sleep(delay)
	fmt.Printf("[worker-%d] 📦  Batch job #%d — processed %d items in %s\n",
		workerID, p.JobID, p.Items, delay)
	return nil
}

// ── Worker ───────────────────────────────────────────────────────────────────

type Worker struct {
	id          int
	queue       *queue.Queue
	stop        chan struct{}
	wg          *sync.WaitGroup
	processed   []string
	mu          sync.Mutex
	statusStore *task.StatusStore
}

func NewWorker(id int, queue *queue.Queue, wg *sync.WaitGroup, store *task.StatusStore) *Worker {
	return &Worker{
		id:          id,
		queue:       queue,
		stop:        make(chan struct{}),
		wg:          wg,
		processed:   make([]string, 0),
		statusStore: store,
	}
}

func (w *Worker) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.run()
	}()
}

func (w *Worker) Stop() {
	close(w.stop)
}

func (w *Worker) Wait() {
	w.wg.Wait()
}

func (w *Worker) run() {
	for {
		select {
		case <-w.stop:
			return
		case t := <-w.queue.Tasks():
			t.Status = task.StatusRunning
			w.statusStore.Set(t)

			err := w.dispatch(t)

			if err != nil {
				t.Status = task.StatusFailed
				t.Error = err.Error()
			} else {
				t.Status = task.StatusCompleted
			}
			w.statusStore.Set(t)
		}
	}
}

// dispatch routes the task to the right processor, falling back to a generic one.
func (w *Worker) dispatch(t *task.Task) error {
	w.mu.Lock()
	w.processed = append(w.processed, t.ID)
	w.mu.Unlock()

	fn, ok := processors[t.Type]
	if !ok {
		fmt.Printf("[worker-%d] ⚙️  Unknown task type %q — running generic handler\n", w.id, t.Type)
		time.Sleep(100 * time.Millisecond)
		return nil
	}
	return fn(w.id, t.Payload)
}

func (w *Worker) GetProcessed() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	result := make([]string, len(w.processed))
	copy(result, w.processed)
	return result
}
