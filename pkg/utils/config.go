package utils

import (
	"net/url"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
)

// Configurations wraps all the config variables required by the auth service
type Configurations struct {
	ServerAddress              string
	MONGO_URI                  string
	DBName                     string
	DBUser                     string
	DBPass                     string
	AccessTokenPrivateKeyPath  string
	AccessTokenPublicKeyPath   string
	RefreshTokenPrivateKeyPath string
	RefreshTokenPublicKeyPath  string
	JwtExpiration              int // in minutes
	RefreshJwtExpiration       int // in minutes
	SendGridApiKey             string
	MailVerifCodeExpiration    int // in hours
	PassResetCodeExpiration    int // in minutes
	MailVerifTemplateID        string
	PassResetTemplateID        string
	AssetsUrl                  string
}

// NewConfigurations returns a new Configuration object
func NewConfigurations(logger hclog.Logger) *Configurations {

	viper.AutomaticEnv()

	logger.Debug("found database url in env, connection string is formed by parsing it")

	viper.SetDefault("SERVER_ADDRESS", "localhost:8000")
	viper.SetDefault("MONGO_URI", "mongodb://localhost:27017")
	viper.SetDefault("DB_NAME", "booksland")
	viper.SetDefault("DB_USER", "")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("ACCESS_TOKEN_PRIVATE_KEY_PATH", "./access-private.pem")
	viper.SetDefault("ACCESS_TOKEN_PUBLIC_KEY_PATH", "./access-public.pem")
	viper.SetDefault("REFRESH_TOKEN_PRIVATE_KEY_PATH", "./refresh-private.pem")
	viper.SetDefault("REFRESH_TOKEN_PUBLIC_KEY_PATH", "./refresh-public.pem")
	viper.SetDefault("JWT_EXPIRATION", 30)
	viper.SetDefault("REFRESH_JWT_EXPIRATION", 24*60*30)
	viper.SetDefault("MAIL_VERIFICATION_CODE_EXPIRATION", 30)
	viper.SetDefault("PASSWORD_RESET_CODE_EXPIRATION", 15)
	viper.SetDefault("MAIL_VERIFICATION_TEMPLATE_ID", "d-5ecbea6e38764af3b703daf03f139b48")
	viper.SetDefault("PASSWORD_RESET_TEMPLATE_ID", "d-3fc222d11809441abaa8ed459bb44319")
	viper.SetDefault("ASSETS_URL", "http://localhost:8000")

	configs := &Configurations{
		ServerAddress:              viper.GetString("SERVER_ADDRESS"),
		MONGO_URI:                  viper.GetString("MONGO_URI"),
		DBName:                     viper.GetString("DB_NAME"),
		DBUser:                     viper.GetString("DB_USER"),
		DBPass:                     viper.GetString("DB_PASSWORD"),
		JwtExpiration:              viper.GetInt("JWT_EXPIRATION"),
		RefreshJwtExpiration:       viper.GetInt("REFRESH_JWT_EXPIRATION"),
		AccessTokenPrivateKeyPath:  viper.GetString("ACCESS_TOKEN_PRIVATE_KEY_PATH"),
		AccessTokenPublicKeyPath:   viper.GetString("ACCESS_TOKEN_PUBLIC_KEY_PATH"),
		RefreshTokenPrivateKeyPath: viper.GetString("REFRESH_TOKEN_PRIVATE_KEY_PATH"),
		RefreshTokenPublicKeyPath:  viper.GetString("REFRESH_TOKEN_PUBLIC_KEY_PATH"),
		SendGridApiKey:             viper.GetString("SENDGRID_API_KEY"),
		MailVerifCodeExpiration:    viper.GetInt("MAIL_VERIFICATION_CODE_EXPIRATION"),
		PassResetCodeExpiration:    viper.GetInt("PASSWORD_RESET_CODE_EXPIRATION"),
		MailVerifTemplateID:        viper.GetString("MAIL_VERIFICATION_TEMPLATE_ID"),
		PassResetTemplateID:        viper.GetString("PASSWORD_RESET_TEMPLATE_ID"),
		AssetsUrl:                  viper.GetString("ASSETS_URL"),
	}

	// reading heroku provided port to handle deployment with heroku
	port := viper.GetString("PORT")
	if port != "" {
		logger.Debug("using the port allocated by heroku", port)
		configs.ServerAddress = "0.0.0.0:" + port
	}

	logger.Debug("serve port", configs.ServerAddress)
	logger.Debug("MONGO_URI", configs.MONGO_URI)
	logger.Debug("db name", configs.DBName)
	logger.Debug("jwt expiration", configs.JwtExpiration)

	return configs
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (config *Configurations) AppendUrl(url string) string {
	if IsUrl(url) {
		return url
	}
	return config.AssetsUrl + url
}
