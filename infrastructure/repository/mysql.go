package repository

import (
	"database/sql"
	"g-tech.com/infrastructure/logger"
	_ "github.com/go-sql-driver/mysql"
)

/**
 * See more
 * https://pseudomuto.com/2018/01/clean-sql-transactions-in-golang/
 * https://github.com/yanpozka/sqlite_trans/blob/master/main.go
 */
type MySqlRepository struct {
	DbContext *sql.DB
}

type MySqlRowScanner interface {
	Scan(dest ... interface{}) error
}

/**
 * Sets MySql DbContext
 */
func (respository *MySqlRepository) SetDbContext(dbContext *sql.DB) {
	respository.DbContext = dbContext
}

/**
 * Handles MySql error
 */
func (respository *MySqlRepository) HandleError(err error) {
	logger.Error("[MySql]", err.Error())
}

/**
 * Initializes a MySql infrastructure
 */
func ConnectMySql(host string, userName string, password string, database string, maxOpenConnections int, maxIdleConnections int) (db *sql.DB) {
	db, err := sql.Open("mysql", userName + ":" + password + "@tcp(" + host + ")/" + database + "?multiStatements=true&parseTime=true")

	if err != nil {
		logger.Error("Failed to connect to MySql", err.Error())
		return nil
	}

	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)

	return db
}
