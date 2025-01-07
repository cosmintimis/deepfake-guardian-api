package postgresql

import (
	"context"
	"fmt"

	"github.com/cosmintimis/deepfake-guardian-api/internal/config"
	"github.com/jackc/pgx/v5"
)

var dbConnection *pgx.Conn

func InitDB() (*pgx.Conn, error) {
	globalConfig := config.GetConfig()
	dbConnection, err := pgx.Connect(context.Background(), globalConfig.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// SQL statement to create the `media` table if it doesn't exist
	createMediaTableSQL := `
        CREATE TABLE IF NOT EXISTS media (
            id TEXT PRIMARY KEY NOT NULL,         
            title TEXT NOT NULL,                
            description TEXT,                   
            location TEXT,                     
            type TEXT NOT NULL,                 
            mimeType TEXT NOT NULL,               
            size INTEGER NOT NULL,               
            tags TEXT,                          
            mediaData BYTEA NOT NULL         
        );
    `
	_, err = dbConnection.Exec(context.Background(), createMediaTableSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to create media table: %w", err)
	}

	return dbConnection, nil
}

func GetDBConnection() *pgx.Conn {
	return dbConnection
}
