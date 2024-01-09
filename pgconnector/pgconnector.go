// pgconnector package encapsulates PostgreSQL connection management functionality.

package pgconnector

import (
    "context"
    "fmt"
    "log"
    "time"
    "github.com/jackc/pgx/v5"
)

// PGClientConfig represents the PostgreSQL client configuration.
type PGClientConfig struct {
  Hostname         string  `json:"Hostname"`
  Port             string  `json:"Port"`
  Database         string  `json:"Database"`
  Username         string  `json:"Username"`
  Password         string  `json:"Password"`
  Sslmode          string  `json:"Sslmode"`
  ApplicationName  string  `json:"ApplicationName"`
}


// PGManager represents the PostgreSQL connection manager.
type PGManager struct {
	conn   *pgx.Conn
	Config *PGClientConfig
}

// NewPGManager creates a new PGManager instance with the given configuration.
func NewPGManager(configPath string) (*PGManager, error) {
	config, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &PGManager{
		Config: config,
	}, nil
}

func loadConfig(configPath string) (*PGClientConfig, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &PGClientConfig{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// PGConnect establishes a new connection to the PostgreSQL database.
func (pm *PGManager) PGConnect() (*pgx.Conn, error) {
  connStr := fmt.Sprintf("hostname=%s port=%d dbname=%s user=%s password=%s sslmode=%s applicationame=%s", pm.Config.Hostname, pm.Config.Port, pm.Config.Database, pm.Config.Username, pm.Config.Password, pm.Config.Sslmode, pm.Config.ApplicationName)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	pm.conn = conn
	return conn, nil
}

// PGReconnectWithTimeout attempts to reconnect to the PostgreSQL database within a specified timeout.
func (pm *PGManager) PGReconnectWithTimeout(timeout time.Duration) error {
	startTime := time.Now()

	for time.Since(startTime) < timeout {
		log.Println("Attempting to reconnect...")

		err := pm.pgConnectWithRetry()
		if err == nil {
			log.Println("Reconnected successfully.")
			return nil
		}

		log.Printf("Reconnection attempt failed: %v\n", err)
		time.Sleep(2 * time.Second) // Adjust the retry interval based on your needs
	}

	return fmt.Errorf("failed to reconnect within %s", timeout)
}

// pgConnectWithRetry tries to connect to the PostgreSQL database.
func (pm *PGManager) pgConnectWithRetry() error {
	conn, err := pm.PGConnect()
	if err != nil {
		return err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return err
	}

	pm.conn = conn
	return nil
}
