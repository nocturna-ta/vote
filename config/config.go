package config

import (
	"github.com/nocturna-ta/golib/config"
	"github.com/nocturna-ta/golib/log"
	"time"
)

type (
	MainConfig struct {
		Server     ServerConfig     `yaml:"Server"`
		API        APIConfig        `yaml:"API"`
		Database   DBConfig         `yaml:"Database"`
		Blockchain BlockchainConfig `yaml:"BlockchainConfig"`
		JWT        JWTConfig        `yaml:"JWT"`
		Kafka      KafkaConfig      `yaml:"Kafka"`
		Encryption EncryptionConfig `yaml:"Encryption"`
		GrpcServer GrpcServerConfig `yaml:"GrpcServer"`
		Redis      RedisConfig      `yaml:"RedisConfig"`
		OTP        OTPConfig        `yaml:"OTP"`
		SMS        SMSConfig        `yaml:"SMS"`
	}

	ServerConfig struct {
		Port         uint          `yaml:"Port" env:"SERVER_PORT"`
		WriteTimeout time.Duration `yaml:"WriteTimeout" env:"SERVER_WRITE_TIMEOUT"`
		ReadTimeout  time.Duration `yaml:"ReadTimeout" env:"SERVER_READ_TIMEOUT"`
	}

	APIConfig struct {
		BasePath      string        `yaml:"BasePath" env:"API_BASE_PATH"`
		APITimeout    time.Duration `yaml:"APITimeout" env:"API_TIMEOUT"`
		EnableSwagger bool          `yaml:"EnableSwagger" env:"ENABLE_SWAGGER" default:"false"`
	}

	DBConfig struct {
		SlaveDSN        string `yaml:"SlaveDSN" env:"DB_SLAVE_DSN"`
		MasterDSN       string `yaml:"MasterDSN" env:"DB_MASTER_DSN"`
		RetryInterval   int    `yaml:"RetryInterval" env:"DB_RETRY_INTERVAL"`
		MaxIdleConn     int    `yaml:"MaxIdleConn" env:"DB_MAX_IDLE_CONN"`
		MaxConn         int    `yaml:"MaxConn" env:"DB_MAX_CONN"`
		ConnMaxLifetime string `yaml:"ConnMaxLifetime" env:"DB_CONN_MAX_LIFETIME"`
	}

	BlockchainConfig struct {
		GanacheURL             string `yaml:"GanacheURL"`
		VotechainAddress       string `yaml:"VotechainAddress" `
		VotechainBaseAddress   string `yaml:"VotechainBaseAddress"`
		KPUManagerAddress      string `yaml:"KPUManagerAddress"`
		VoterManagerAddress    string `yaml:"VoterManagerAddress"`
		ElectionManagerAddress string `yaml:"ElectionManagerAddress"`
	}

	JWTConfig struct {
		Secret string `yaml:"Secret" env:"JWT_SECRET"`
	}

	EncryptionConfig struct {
		Key string `yaml:"Key" env:"ENCRYPTION_KEY"`
	}

	GrpcServerConfig struct {
		Port uint `yaml:"Port"`
	}

	KafkaConfig struct {
		Producer KafkaProducerConfig `yaml:"Producer"`
		Consumer KafkaConsumerConfig `yaml:"Consumer"`
		Topics   KafkaTopics         `yaml:"Topics"`
	}

	KafkaProducerConfig struct {
		Brokers    []string `yaml:"Brokers" env:"KAFKA_BROKERS"`
		Idempotent bool     `yaml:"Idempotent" env:"KAFKA_IDEMPOTENT"`
		MaxAttempt int      `yaml:"MaxAttempt" env:"KAFKA_MAX_ATTEMPTS"`
	}

	KafkaConsumerConfig struct {
		Brokers        []string    `yaml:"Brokers"`
		ClusterVersion string      `yaml:"ClusterVersion"`
		ConsumerGroup  string      `yaml:"ConsumerGroup"`
		MaxRetries     int         `yaml:"MaxRetries"`
		WorkerPoolSize int         `yaml:"WorkerPoolSize"`
		MaxAttempt     int         `yaml:"MaxAttempt"`
		Retry          RetryConfig `yaml:"Retry"`
	}

	RetryConfig struct {
		MaxRetry          int             `yaml:"MaxRetry"`
		RetryInitialDelay time.Duration   `yaml:"RetryInitialDelay"`
		MaxJitter         time.Duration   `yaml:"MaxJitter"`
		HandlerTimeout    time.Duration   `yaml:"HandlerTimeout"`
		BackOffConfig     []time.Duration `yaml:"BackOffConfig"`
	}

	KafkaTopics struct {
		VoteSubmitData KafkaTopicConfig `yaml:"VoteSubmitData"`
		VoteProcessed  KafkaTopicConfig `yaml:"VoteProcessed"`
		VoteDLQ        KafkaTopicConfig `yaml:"VoteDLQ"`
	}

	KafkaTopicConfig struct {
		Value        string `yaml:"Value" env:"KAFKA_TOPIC_VALUE"`
		ErrorHandler string `yaml:"ErrorHandler"`
		WithBackOff  bool   `yaml:"WithBackOff"`
	}
	RedisConfig struct {
		Connection string `yaml:"Connection"`
	}

	OTPConfig struct {
		Length     int           `yaml:"Length" env:"OTP_LENGTH" default:"6"`
		TTL        time.Duration `yaml:"TTL" env:"OTP_TTL" default:"5m"`
		MaxRetries int           `yaml:"MaxRetries" env:"OTP_MAX_RETRIES" default:"3"`
		Enabled    bool          `yaml:"Enabled" env:"OTP_ENABLED" default:"true"`
		SendSMS    bool          `yaml:"SendSMS" env:"OTP_SEND_SMS" default:"false"`
	}

	SMSConfig struct {
		Provider  string       `yaml:"Provider" env:"SMS_PROVIDER"`
		Enabled   bool         `yaml:"Enabled" env:"SMS_ENABLED" default:"false"`
		Twilio    TwilioConfig `yaml:"Twilio"`
		Templates SMSTemplates `yaml:"Templates"`
	}

	TwilioConfig struct {
		AccountSID string `yaml:"AccountSID" env:"TWILIO_ACCOUNT_SID"`
		AuthToken  string `yaml:"AuthToken" env:"TWILIO_AUTH_TOKEN"`
		FromNumber string `yaml:"FromNumber" env:"TWILIO_FROM_NUMBER"`
	}

	SMSTemplates struct {
		OTPMessage string `yaml:"OTPMessage"`
	}
)

func ReadConfig(cfg any, configLocation string) {
	if configLocation == "" {
		configLocation = "file://config/files/config.yaml"
	}

	if err := config.ReadConfig(cfg, configLocation, true); err != nil {
		log.WithFields(log.Fields{
			"error":           err,
			"config-location": configLocation,
		}).Fatal("Failed to read config")
	}
}
