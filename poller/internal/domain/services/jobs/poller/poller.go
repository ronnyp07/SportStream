package pooller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/ronnyp07/SportStream/internal/domain/models"
	ports "github.com/ronnyp07/SportStream/internal/domain/ports/jobbuilder"
	metrics_port "github.com/ronnyp07/SportStream/internal/domain/ports/metrics"
	"github.com/ronnyp07/SportStream/internal/pkg/config"
	"github.com/ronnyp07/SportStream/internal/pkg/infaestructure/log"
	natsQueue "github.com/ronnyp07/SportStream/internal/pkg/infaestructure/msgqueue"
	"github.com/sts-solutions/base-code/cchttp"
	"github.com/sts-solutions/base-code/ccretry"
)

const (
	name            = "poller"
	externalService = "articleservice"
	natsSubject     = "SPORTSTREAM.status.updated"
	pageSize        = 2
)

type Job struct {
	scheduler     gocron.Scheduler
	jobConfig     config.Job
	jobBuilder    ports.JobBuilder
	httpclient    cchttp.Client
	metrics       metrics_port.SchedulerMetricsHandler
	msgQueueServ  natsQueue.MsgQueueService
	currentPage   int
	maxPages      int
	lastFetchTime time.Time
	stateMutex    sync.Mutex
}

func New(
	scheduler gocron.Scheduler,
	jobBuilder ports.JobBuilder,
	httpclient cchttp.Client,
	metrics metrics_port.SchedulerMetricsHandler,
	msgQueueServ natsQueue.MsgQueueService,
) *Job {
	return &Job{
		scheduler:    scheduler,
		jobBuilder:   jobBuilder,
		httpclient:   httpclient,
		metrics:      metrics,
		msgQueueServ: msgQueueServ,
	}
}

func (j *Job) Configure(ctx context.Context, jobConfig config.Job) error {
	_, err := j.jobBuilder.BuildJob(ctx, name, jobConfig, j.runTask)
	if err != nil {
		log.Logger().Error(ctx, fmt.Sprintf("An error occurred configuring the job %s due to %s", name, err.Error()))
		return err
	}
	j.jobConfig = jobConfig
	return nil
}

func (j *Job) Name() string {
	return name
}

func (j *Job) runTask(ctx context.Context) {

	j.metrics.ReportScheduleOfJob(j.Name())
	log.Logger().Info(ctx, fmt.Sprintf("executing job %s at %s", j.Name(), time.Now()))

	j.stateMutex.Lock()
	defer j.stateMutex.Unlock()

	url := fmt.Sprintf("%s/?page=%d&pageSize=%d", j.jobConfig.ExternalAddrs, j.currentPage, pageSize)
	reqMethod := http.MethodGet

	rb := cchttp.NewRequestBuilder().
		WithHTTPMethod(reqMethod).
		WithURL(url).
		WithCorrelationIDHeaderFromContext(ctx).
		WithHTTPClient(j.httpclient)

	req, err := rb.Build()
	if err != nil {
		log.Logger().Error(ctx, fmt.Sprintf("unable to build the request for job %s due to %s", name, err.Error()))
	}

	retryDuration, err := time.ParseDuration(j.jobConfig.Retry.Duration)
	if err != nil {
		log.Logger().Error(ctx, fmt.Sprintf("invalid retry duration for job %s due to %s", name, err.Error()))
	}
	var responseCode int
	var actualCallStart, actualCallEnd time.Time

	retry := ccretry.NewRetry(func() error {
		actualCallStart = time.Now()
		externalResponse, err := j.httpclient.Do(req.HTTPRequest())
		actualCallEnd = time.Now()

		if externalResponse != nil {
			responseCode = externalResponse.StatusCode

			bodyBytes, err := io.ReadAll(externalResponse.Body)
			if err != nil {
				return fmt.Errorf("error reading response body: %s", err.Error())
			}

			if closedErr := externalResponse.Body.Close(); closedErr != nil {
				return fmt.Errorf("error closing response body: %s", closedErr.Error())
			}
			if externalResponse.StatusCode != http.StatusOK {
				return fmt.Errorf("invalid response code %d", externalResponse.StatusCode)
			}

			// Unmarshal into our struct
			var apiResponse models.ArticleResponse
			if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
				return fmt.Errorf("error unmarshaling response: %s", err.Error())
			}

			// Update max pages if needed
			if apiResponse.PageInfo.NumPages > j.maxPages {
				j.maxPages = apiResponse.PageInfo.NumPages
			}

			// Convert back to clean JSON
			cleanJSON, err := json.Marshal(apiResponse.Content)
			if err != nil {
				return fmt.Errorf("error marshaling to JSON: %s", err.Error())
			}

			j.msgQueueServ.PublishMessage(ctx, "SPORTSTREAM.DOCKER.status.updated", bodyBytes)

			if _, err := j.msgQueueServ.PublishMessage(ctx, natsSubject, cleanJSON); err != nil {
				return fmt.Errorf("error publishing to NATS: %s", err.Error())
			}

			log.Logger().Debug(ctx, fmt.Sprintf("Successfully published articles to NATS %v",
				map[string]interface{}{
					"article_count": len(apiResponse.Content), "subject": natsSubject,
				}))
		}

		if err != nil {
			return err
		}

		return nil
	}).WithMaxAttempts(j.jobConfig.Retry.MaxAttempts).
		WithSleep(retryDuration)

	defer func() {
		j.handleMetrics(actualCallStart, reqMethod, url, responseCode, err)
	}()

	retryResponse, err := retry.Run()
	if err != nil {
		log.Logger().Error(ctx, fmt.Sprintf("job execution failed %v", map[string]interface{}{
			"name":             j.Name(),
			"numberOfAttempts": retryResponse.NumberOfAttempts(),
			"error":            err.Error(),
		}))

		return
	}

	if retryResponse.NumberOfAttempts() > 1 {
		log.Logger().Info(ctx, fmt.Sprintf("retry succeeded", retryResponse.String()))
	}

	actualCallLatency := actualCallEnd.Sub(actualCallStart)
	j.lastFetchTime = time.Now()

	// Move to next page or wrap around
	j.currentPage++
	if j.currentPage >= j.maxPages {
		j.currentPage = 0
	}
	log.Logger().Info(ctx, fmt.Sprintf("job execution completed %v", map[string]interface{}{
		"name":          j.Name(),
		"time":          time.Now(),
		"latency":       actualCallLatency,
		"latencyString": actualCallLatency.String(),
		"request":       req,
	}))
}

func (j *Job) handleMetrics(start time.Time, method string, url string, responseCode int, err error) {
	j.metrics.SchedulerOutgoingHttpRequest(start, method, url, responseCode, externalService)
	if err != nil {
		j.metrics.JobErrorInc(name, err.Error())
	}
	j.metrics.SchedulerHTTPClientCall(start, method, url, responseCode, externalService)
}
