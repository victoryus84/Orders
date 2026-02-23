package migrations

import (
	"log"

	"gorm.io/gorm"
)

// ColumnInfo holds information about a column
type ColumnInfo struct {
	Name    string
	Type    string
	InModel bool
	InDB    bool
}

// GetDBColumns retrieves column names from database for a specific table
func GetDBColumns(db *gorm.DB, tableName string) ([]string, error) {
	var columns []string
	result := db.Raw(`
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_name = ? 
		ORDER BY ordinal_position
	`, tableName).Scan(&columns)
	return columns, result.Error
}

// AnalyzeSchemaSync compares DB schema with models and reports discrepancies
func AnalyzeSchemaSync(db *gorm.DB) {
	log.Println("\n======= 🔍 DATABASE SCHEMA ANALYSIS =======")

	tables := GetAllModels()

	for _, table := range tables {
		analyzeTable(db, table)
	}

	log.Println("\n========================================")
}

func analyzeTable(db *gorm.DB, table interface{}) {
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(table)

	tableName := stmt.Table

	// Get columns from model
	modelColumns := make(map[string]bool)
	for _, field := range stmt.Schema.Fields {
		modelColumns[field.DBName] = true
	}

	// Get columns from database
	dbColumns, err := GetDBColumns(db, tableName)
	if err != nil {
		log.Printf("❌ Error reading table %s: %v\n", tableName, err)
		return
	}

	dbColumnsMap := make(map[string]bool)
	for _, col := range dbColumns {
		dbColumnsMap[col] = true
	}

	// Find discrepancies
	orphanedColumns := []string{}
	missingColumns := []string{}

	for col := range dbColumnsMap {
		if !modelColumns[col] {
			orphanedColumns = append(orphanedColumns, col)
		}
	}

	for col := range modelColumns {
		if !dbColumnsMap[col] {
			missingColumns = append(missingColumns, col)
		}
	}

	// Print report
	if len(orphanedColumns) == 0 && len(missingColumns) == 0 {
		log.Printf("✅ %s - SYNCED (DB: %d cols, Model: %d cols)\n",
			tableName, len(dbColumnsMap), len(modelColumns))
	} else {
		log.Printf("\n⚠️  TABLE: %s\n", tableName)
		log.Printf("   DB Columns: %d | Model Columns: %d\n", len(dbColumnsMap), len(modelColumns))

		if len(orphanedColumns) > 0 {
			log.Printf("   🗑️  ORPHANED IN DB (can be deleted):\n")
			for _, col := range orphanedColumns {
				log.Printf("      - %s\n", col)
			}
		}

		if len(missingColumns) > 0 {
			log.Printf("   ❌ MISSING IN DB (need to add):\n")
			for _, col := range missingColumns {
				log.Printf("      - %s\n", col)
			}
		}
	}
}

// PrintSyncCommands generates DROP COLUMN commands for cleanup
func PrintSyncCommands(db *gorm.DB) {
	log.Println("\n======= 📋 AUTO-FIX COMMANDS =======")

	tables := GetAllModels()

	for _, table := range tables {
		stmt := &gorm.Statement{DB: db}
		stmt.Parse(table)

		tableName := stmt.Table

		// Get columns from model
		modelColumns := make(map[string]bool)
		for _, field := range stmt.Schema.Fields {
			modelColumns[field.DBName] = true
		}

		// Get columns from database
		dbColumns, err := GetDBColumns(db, tableName)
		if err != nil {
			continue
		}

		// Find orphaned columns
		for _, col := range dbColumns {
			if !modelColumns[col] {
				log.Printf("db.Migrator().DropColumn(&models.%s{}, \"%s\")\n",
					TableNameToModel(tableName), col)
			}
		}
	}

	log.Println("\n====================================")
}
