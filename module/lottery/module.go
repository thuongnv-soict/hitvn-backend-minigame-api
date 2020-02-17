package lottery

import (
	"database/sql"
	"encoding/json"
	"g-tech.com/constant"
	"g-tech.com/dto"
	"g-tech.com/infrastructure/broker"
	"g-tech.com/infrastructure/cache"
	"g-tech.com/infrastructure/logger"
	"g-tech.com/module/lottery/summary"
	"g-tech.com/module/minigame/service"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

var mLotteryResultService summary.LotterySummaryService

var wg sync.WaitGroup
var mRbChannel 		   		*amqp.Channel
var mRbResultConQueue 		amqp.Queue

func Initialize(rbChannel *amqp.Channel, dbContext *sql.DB, cache cache.CacheManager, timeout time.Duration){

	redisService := service.NewRedisService(dbContext, cache, timeout)
	configService := service.NewConfigService(dbContext, cache, redisService, timeout)
	mLotteryResultService = summary.NewLotterySummaryService(dbContext, cache, redisService, configService, timeout)

	mRbChannel = rbChannel
	// Creates a queue to consume to crawl post
	mRbResultConQueue = broker.CreateQueue(rbChannel, constant.RbRouteResult)
	err := rbChannel.QueueBind (
		mRbResultConQueue.Name,
		constant.RbRouteResult,
		constant.RbSuperExchange,
		false,
		nil,
	)
	if err != nil {
		logger.Error(err.Error())
	}
}

func Execute()  {
	forever := make(chan bool)

	wg.Add(1)
	go consumeLotteryResult()
	wg.Wait()

	<-forever
}

/**
 * Consumes posts
 */
func consumeLotteryResult()  {
	results, err := mRbChannel.Consume(
		mRbResultConQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("Failed to consume post", err.Error())
	}

	for item := range results {
		var lotteryResult dto.LotteryResult
		err := json.Unmarshal(item.Body, &lotteryResult)
		if err == nil {
			// Then crawl it
			//mPostService.CrawlPost(mRbChannel, &mRbResultPubQueue, &mRbPostETLQueue, crawlPost.Data)
			err = mLotteryResultService.SummaryResult(lotteryResult)
		}
		// Return Ack
		_ = item.Ack(false)
	}
	wg.Done()
}