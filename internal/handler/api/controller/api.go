package controller

import (
	"github.com/gofiber/swagger"
	"github.com/nocturna-ta/golib/router"
	_ "github.com/nocturna-ta/vote/docs"
	"github.com/nocturna-ta/vote/internal/usecases"
	"github.com/nocturna-ta/vote/pkg/utils"
	"html/template"
	"time"
)

type API struct {
	prefix         string
	port           uint
	readTimeout    time.Duration
	writeTimeout   time.Duration
	requestTimeout time.Duration
	enableSwagger  bool
	voteUc         usecases.VoteUseCases
	electionTimeUc usecases.ElectionTimeUseCases
	otpUc          usecases.OTPUseCases
}

type Options struct {
	Prefix         string
	Port           uint
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	RequestTimeout time.Duration
	EnableSwagger  bool
	VoteUc         usecases.VoteUseCases
	ElectionTimeUc usecases.ElectionTimeUseCases
	OTPUc          usecases.OTPUseCases
}

func New(opts *Options) *API {
	return &API{
		prefix:         opts.Prefix,
		port:           opts.Port,
		readTimeout:    opts.ReadTimeout,
		writeTimeout:   opts.WriteTimeout,
		requestTimeout: opts.RequestTimeout,
		enableSwagger:  opts.EnableSwagger,
		voteUc:         opts.VoteUc,
		electionTimeUc: opts.ElectionTimeUc,
		otpUc:          opts.OTPUc,
	}
}

func (api *API) RegisterRoute() *router.FastRouter {
	myRouter := router.New(&router.Options{
		Prefix:         api.prefix,
		Port:           api.port,
		ReadTimeout:    api.readTimeout,
		WriteTimeout:   api.writeTimeout,
		RequestTimeout: api.requestTimeout,
	})

	if api.enableSwagger {
		swaggerConfig := swagger.Config{
			Title:        "API Documentation",
			DeepLinking:  true,
			DocExpansion: "list",
			CustomStyle:  template.CSS(utils.ClaudeDarkTheme),
		}

		myRouter.CustomHandler("GET", "/vote/docs/*", swagger.New(swaggerConfig), router.MustAuthorized(false))
	}

	myRouter.GET("/health", api.Ping, router.MustAuthorized(false))
	myRouter.Group("/v1", func(v1 *router.FastRouter) {
		v1.Group("/vote", func(vote *router.FastRouter) {
			vote.POST("/cast", api.CastVote, router.MustAuthorized(false))
			vote.GET("/:id/status", api.GetVoteStatus, router.MustAuthorized(false))
		})
		v1.Group("/election-time", func(electionTime *router.FastRouter) {
			electionTime.GET("/status", api.GetElectionStatus, router.MustAuthorized(false))
			electionTime.POST("/", api.CreateElectionTime, router.MustAuthorized(false))
			electionTime.GET("/:id", api.GetElectionTimeByID, router.MustAuthorized(false))
			electionTime.PUT("/:id", api.UpdateElectionTime, router.MustAuthorized(false))
			electionTime.DELETE("/:id", api.DeleteElectionTime, router.MustAuthorized(false))
			electionTime.POST("/:id/activate", api.ActivateElection, router.MustAuthorized(false))
			electionTime.POST("/sync", api.SyncElectionStatuses, router.MustAuthorized(false))
		})
		v1.Group("/otp", func(otp *router.FastRouter) {
			otp.POST("/generate", api.GenerateOTP, router.MustAuthorized(false))
			otp.POST("/verify", api.VerifyOTP, router.MustAuthorized(false))
			otp.POST("/resend", api.ResendOTP, router.MustAuthorized(false))
			otp.GET("/status", api.GetOTPStatus, router.MustAuthorized(false))
		})

	})
	return myRouter
}
