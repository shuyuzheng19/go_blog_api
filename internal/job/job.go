package job

import (
	"blog/pkg/logger"
	"log"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Job struct {
	Hour        int // 24表示每天凌晨执行
	Description string
	Eq          bool
	Job         func()
}

var (
	jobs  = []Job{}
	mutex sync.Mutex
)

func AddJob(job Job) {
	mutex.Lock()
	defer mutex.Unlock()
	jobs = append(jobs, job)
}

func StartJob() {
	go func() {
		// 打印任务列表
		log.Println("==================================任务列表==================================")
		mutex.Lock()
		for _, job := range jobs {
			if job.Eq {
				log.Printf("任务：%s - 每天到%d点执行", job.Description, job.Hour)
			} else {
				log.Printf("任务：%s - 每%d小时执行", job.Description, job.Hour)
			}
		}
		mutex.Unlock()
		log.Println("==================================任务列表==================================")

		// 等待到下一个整点
		now := time.Now()
		next := now.Truncate(time.Hour).Add(time.Hour)
		time.Sleep(time.Until(next))

		// 每小时检查任务
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()

		for {
			hour := time.Now().Hour()

			mutex.Lock()
			for _, job := range jobs {
				if (job.Eq && job.Hour == hour) || (!job.Eq && hour%job.Hour == 0) {
					logger.Info("执行任务", zap.String("desc", job.Description), zap.Bool("eq", job.Eq), zap.Int("hour", job.Hour))
					go job.Job()
				}
			}
			mutex.Unlock()

			<-ticker.C
		}
	}()
}
