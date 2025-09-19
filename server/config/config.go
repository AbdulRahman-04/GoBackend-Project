package config

type Config struct {
	AppName string
	Port    int
	DBURI   string
	URL     string
	JWTKEY  string
	Email   EmailConfig
	Phone   PhoneConfig
	Redis   RedisConfig // 🔥 add this
}

type EmailConfig struct {
	User string
	Pass string
}

type PhoneConfig struct {
	Sid   string
	Token string
	Phone string
}

type RedisConfig struct { // 🔥 add this
	Host     string
	Password string
	DB       int
}

var AppConfig = &Config{
	AppName: "Event_Booking",
	Port:    4040,
	DBURI:   "your mongodb url",
	URL:     "http://localhost:4040",
	JWTKEY:  "your_jwt_secret_here",
	Email: EmailConfig{
		User: "your gmail id",
		Pass: "your app password",
	},
	Phone: PhoneConfig{
		Sid:   "your_twilio_sid_here",
		Token: "your_twilio_token_here",
		Phone: "+1234567890",
	},
	Redis: RedisConfig{ // 🔥 add this
		Host:     "localhost:6379",
		Password: "",
		DB:       0,
	},
}
