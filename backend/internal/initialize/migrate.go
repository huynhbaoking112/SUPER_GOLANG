package initialize

import (
	"fmt"
	"go-backend-v2/global"
	"go-backend-v2/internal/models"
)

func InitMigrations() {
	db := global.DB

	err := db.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.UserAuthProvider{},
		&models.Workspace{},
		&models.WorkspaceRole{},
		&models.UserWorkspaceMembership{},
		&models.Resource{},
	)

	if err != nil {
		panic(fmt.Errorf("failed to migrate database: %v", err))
	}

	fmt.Println("Database migrations completed successfully!")
}
