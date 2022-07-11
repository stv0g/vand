package config

import (
	"fmt"
	"image/color"
	"log"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/stv0g/vand/pkg/display/widget"
	"github.com/stv0g/vand/pkg/mqtt"
	"gopkg.in/yaml.v2"
)

type Modem struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	PollInterval time.Duration `yaml:"poll_interval"`
}

type Car struct {
	Address string `yaml:"address"`

	TankVolume float32 `yaml:"tank_volume"`

	PollInterval time.Duration `yaml:"poll_interval"`
}

type Solar struct {
	Address string `yaml:"address"`

	PollInterval time.Duration `yaml:"poll_interval"`
}

type Display struct {
	Port string `yaml:"port"`

	Pins struct {
		DC    string `yaml:"dc"`
		Reset string `yaml:"reset"`
		Next  string `yaml:"next"`
	}

	Pages []DisplayPage `yaml:"pages"`
}

type DisplayPage struct {
	ID     string `yaml:"id"`
	NextID string `yaml:"next"`
	OverID string `yaml:"over"`

	Time            time.Duration   `yaml:"time"`
	BackgroundColor color.Color     `yaml:"background-color"`
	Widgets         []widget.Widget `yaml:"widgets"`
}

type Environment struct {
}

type OwnTracks struct {
	Topic string `yaml:"topic"`
}

type GPS struct {
	Port     string `yaml:"port"`
	Baudrate int    `yaml:"baudrate"`

	MinInterval time.Duration `yaml:"min_interval"`
	MinDistance float64       `yaml:"min_distance"`

	OwnTracks OwnTracks `yaml:"owntracks"`
}

type Web struct {
	// Host is the local machine IP Address to bind the HTTP Server to
	Listen string `yaml:"listen"`

	// Directory of frontend assets if not bundled
	Static string `yaml:"static"`

	// BaseURL at which the VANd web server is accessible
	BaseURL string `yaml:"base_url"`
}

type Bridge struct {
	Flatten bool `yaml:"flatten"`
	Encrypt bool `yaml:"encrypt"`
}

// Config contains the main configuration
type Config struct {
	*viper.Viper `yaml:"-"`

	Debug bool `yaml:"debug"`

	DataDir string `yaml:"data_dir"`

	Broker      mqtt.BrokerConfig `yaml:"broker"`
	BrokerCloud mqtt.BrokerConfig `yaml:"broker_cloud"`

	Display     Display     `yaml:"display"`
	Web         Web         `yaml:"web"`
	GPS         GPS         `yaml:"gps"`
	Car         Car         `yaml:"car"`
	Solar       Solar       `yaml:"solar"`
	Environment Environment `yaml:"env"`
	Modem       Modem       `yaml:"modem"`
}

func decodeOption(cfg *mapstructure.DecoderConfig) {
	cfg.DecodeHook = mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		mapstructure.TextUnmarshallerHookFunc(),
		widget.DecodeHookFunc(),
	)

	cfg.ZeroFields = false
	cfg.TagName = "yaml"
}

// NewConfig returns a new decoded Config struct
func NewConfig(configFile string) (*Config, error) {
	// Create cfg structure
	cfg := &Config{
		Viper: viper.New(),
	}

	cfg.SetDefault("data_dir", "/var/lib/vand")

	cfg.SetDefault("broker.port", 1883)
	cfg.SetDefault("broker.topic", "vand")

	cfg.SetDefault("broker_cloud.port", 1883)
	cfg.SetDefault("broker_cloud.topic", "vand")

	cfg.SetDefault("web.listen", ":8080")
	cfg.SetDefault("web.static", "./frontend/build")
	cfg.SetDefault("web.base_url", "http://localhost:8080")

	cfg.SetDefault("gps.baudrate", 9600)
	cfg.SetDefault("gps.port", "/dev/serial0")
	cfg.SetDefault("gps.min_interval", "15m")
	cfg.SetDefault("gps.min_distance", 100.0)

	replacer := strings.NewReplacer(".", "_")
	cfg.SetEnvKeyReplacer(replacer)
	cfg.SetEnvPrefix("vand")
	cfg.AutomaticEnv()

	if configFile != "" {
		cfg.SetConfigFile(configFile)

		if err := cfg.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		cfg.SetConfigName("vand")
		cfg.AddConfigPath("/etc")
		cfg.AddConfigPath("etc")
		cfg.AddConfigPath(".")

		if err := cfg.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	if err := cfg.UnmarshalExact(cfg, decodeOption); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	log.Printf("Loaded configuration:\n")
	bs, _ := yaml.Marshal(cfg)
	log.Writer().Write(bs)

	return cfg, nil
}
