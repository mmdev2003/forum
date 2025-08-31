package server

import (
	"context"

	"forum-thread/internal/model"
	"forum-thread/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

import (
	messageAmqpHandler "forum-thread/internal/controller/amqp/handler/message"
	subthreadAmqpHandler "forum-thread/internal/controller/amqp/handler/subthread"
	topicAmqpHandler "forum-thread/internal/controller/amqp/handler/topic"
	accountStatisticHttpHandler "forum-thread/internal/controller/http/handler/account-statistic"
	messageHttpHandler "forum-thread/internal/controller/http/handler/message"
	subthreadHttpHandler "forum-thread/internal/controller/http/handler/subthread"
	threadHttpHandler "forum-thread/internal/controller/http/handler/thread"
	topicHttpHandler "forum-thread/internal/controller/http/handler/topic"
	httpMiddleware "forum-thread/internal/controller/http/middleware"
)

func Run(
	db model.IDatabase,
	fullTextSearchEngine model.IFullTextSearchEngine,
	messageBroker model.IMessageBroker,
	threadService model.IThreadService,
	subthreadService model.ISubthreadService,
	topicService model.ITopicService,
	messageService model.IMessageService,
	accountStatisticService model.IAccountStatisticService,
	authorizationClient model.IAuthorizationClient,
	serverPort string,
) {
	server := echo.New()
	server.Validator = validator.New()

	includeSystemHttpHandler(server, db, fullTextSearchEngine)
	includeHttpMiddleware(server, authorizationClient)

	includeThreadHttpHandler(server, threadService)

	includeSubthreadHttpHandler(server, subthreadService)
	includeSubthreadAmqpHandler(messageBroker, subthreadService)

	includeTopicHttpHandler(server, topicService)
	includeTopicAmqpHandler(messageBroker, topicService)

	includeMessageHttpHandler(server, messageService)
	includeMessageAmqpHandler(messageBroker, messageService)

	includeAccountStatisticHandler(server, accountStatisticService)

	server.Logger.Fatal(server.Start(":" + serverPort))
}

func includeThreadHttpHandler(
	server *echo.Echo,
	threadService model.IThreadService,
) {
	server.POST(model.PREFIX+"/create", threadHttpHandler.CreateThread(threadService))
	server.GET(model.PREFIX+"/all", threadHttpHandler.AllThread(threadService))
}

func includeSubthreadHttpHandler(
	server *echo.Echo,
	subthreadService model.ISubthreadService,
) {
	server.POST(model.PREFIX+"/subthread/create", subthreadHttpHandler.CreateSubthread(subthreadService))
	server.POST(model.PREFIX+"/subthread/view", subthreadHttpHandler.AddViewToSubthread(subthreadService))
	server.GET(model.PREFIX+"/subthread/:threadID", subthreadHttpHandler.SubthreadsByThreadID(subthreadService))
}

func includeSubthreadAmqpHandler(
	messageBroker model.IMessageBroker,
	subthreadService model.ISubthreadService,
) {
	go messageBroker.Subscribe(model.AddViewToSubthreadQueue, subthreadAmqpHandler.AddViewToSubthreadPostprocessing(subthreadService))
}

func includeTopicHttpHandler(
	server *echo.Echo,
	topicService model.ITopicService,
) {
	server.POST(model.PREFIX+"/topic/create", topicHttpHandler.CreateTopic(topicService))
	server.POST(model.PREFIX+"/topic/view", topicHttpHandler.AddViewToTopic(topicService))
	server.POST(model.PREFIX+"/topic/approve/:topicID", topicHttpHandler.ApproveTopic(topicService))
	server.POST(model.PREFIX+"/topic/reject/:topicID", topicHttpHandler.RejectTopic(topicService))
	server.POST(model.PREFIX+"/topic/close", topicHttpHandler.CloseTopic(topicService))
	server.POST(model.PREFIX+"/topic/priority/change", topicHttpHandler.ChangeTopicPriority(topicService))
	server.GET(model.PREFIX+"/topic/subthreadID/:subthreadID", topicHttpHandler.TopicsBySubthreadID(topicService))
	server.GET(model.PREFIX+"/topic/accountID/:accountID", topicHttpHandler.TopicsByAccountID(topicService))
	server.GET(model.PREFIX+"/topic/moderation", topicHttpHandler.TopicsOnModeration(topicService))
	server.POST(model.PREFIX+"/topic/avatar/update", topicHttpHandler.UpdateTopicAvatar(topicService))
	server.GET(model.PREFIX+"/topic/avatar/download", topicHttpHandler.DownloadTopicAvatar(topicService))
}

func includeTopicAmqpHandler(
	messageBroker model.IMessageBroker,
	topicService model.ITopicService,
) {
	go messageBroker.Subscribe(model.CreateTopicQueue, topicAmqpHandler.CreateTopicPostprocessing(topicService))
	go messageBroker.Subscribe(model.AddViewToTopicQueue, topicAmqpHandler.AddViewToTopicPostprocessing(topicService))
	go messageBroker.Subscribe(model.CloseTopicQueue, topicAmqpHandler.CloseTopicPostprocessing(topicService))
	go messageBroker.Subscribe(model.ChangeTopicPriorityQueue, topicAmqpHandler.ChangeTopicPriorityPostprocessing(topicService))
}

func includeMessageHttpHandler(
	server *echo.Echo,
	messageService model.IMessageService,
) {
	server.POST(model.PREFIX+"/message/send", messageHttpHandler.SendMessageToTopic(messageService))
	server.POST(model.PREFIX+"/message/like", messageHttpHandler.LikeMessage(messageService))
	server.POST(model.PREFIX+"/message/unlike", messageHttpHandler.UnlikeMessage(messageService))
	server.POST(model.PREFIX+"/message/report", messageHttpHandler.ReportMessage(messageService))
	server.POST(model.PREFIX+"/message/edit", messageHttpHandler.EditMessage(messageService))
	server.GET(model.PREFIX+"/message/topicID/:topicID", messageHttpHandler.MessagesByTopicID(messageService))
	server.GET(model.PREFIX+"/message/accountID/:accountID", messageHttpHandler.MessagesByAccountID(messageService))
	server.GET(model.PREFIX+"/message/search/:text", messageHttpHandler.MessagesByText(messageService))
	server.GET(model.PREFIX+"/file/download/:fileURL", messageHttpHandler.DownloadFile(messageService))
}

func includeMessageAmqpHandler(
	messageBroker model.IMessageBroker,
	messageService model.IMessageService,
) {
	go messageBroker.Subscribe(model.SendMessageToTopicQueue, messageAmqpHandler.SendMessageToTopicPostprocessing(messageService))
	go messageBroker.Subscribe(model.LikeMessageQueue, messageAmqpHandler.LikeMessagePostprocessing(messageService))
	go messageBroker.Subscribe(model.ReportMessageQueue, messageAmqpHandler.ReportMessagePostprocessing(messageService))
}

func includeAccountStatisticHandler(
	server *echo.Echo,
	accountStatisticService model.IAccountStatisticService,
) {
	server.POST(model.PREFIX+"/statistic/create", accountStatisticHttpHandler.CreateAccountStatistic(accountStatisticService))
	server.GET(model.PREFIX+"/statistic/:accountID", accountStatisticHttpHandler.StatisticByAccountID(accountStatisticService))
}

func includeHttpMiddleware(
	server *echo.Echo,
	authorizationClient model.IAuthorizationClient,
) {
	server.Use(httpMiddleware.LoggerMiddleware())
	pathForSkipAuthorization := []string{
		model.PREFIX + "table/create",
		model.PREFIX + "table/drop",
	}
	server.Use(httpMiddleware.AuthorizationMiddleware(authorizationClient, pathForSkipAuthorization))
}

func includeSystemHttpHandler(
	server *echo.Echo,
	db model.IDatabase,
	fullTextSearchEngine model.IFullTextSearchEngine,
) {
	server.GET(model.PREFIX+"/table/create", createTable(db))
	server.GET(model.PREFIX+"/table/drop", dropTable(db))
	server.GET(model.PREFIX+"/fullTextSearchEngine/create", createFullTextSearchEngineIndex(fullTextSearchEngine))
}

func createFullTextSearchEngineIndex(fullTextSearchEngine model.IFullTextSearchEngine) echo.HandlerFunc {
	return func(request echo.Context) error {
		indexes := []string{
			model.MessageFullTextSearchIndex,
		}
		err := fullTextSearchEngine.CreateIndexes(indexes)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, "indexes created")
	}
}

func createTable(db model.IDatabase) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := context.Background()

		err := db.CreateTable(ctx, model.CreateTableQuery)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, "table created")
	}
}

func dropTable(db model.IDatabase) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := context.Background()

		err := db.DropTable(ctx, model.DropTableQuery)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, "table dropped")
	}
}
