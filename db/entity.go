package db

import (
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
	"time"
)

// Entity is an identifiable database entity.
type Entity interface {
	// GetId must return the ID of the entity.
	GetID() int64

	// SetId must set the ID of the entity.
	SetID(int64)
}

// BaseEntity is the base for all database entities.
type BaseEntity struct {
	ID      int64     `json:"id"`
	DefTime time.Time `db:"def_time" json:"def_time"`
	ModTime time.Time `db:"mod_time" json:"mod_time"`
}

// GetID returns the ID of this entity.
func (entity *BaseEntity) GetID() int64 {
	return entity.ID
}

// SetID sets the ID to this entity.
func (entity *BaseEntity) SetID(id int64) {
	entity.ID = id
}

// SaveEntity saves the given entity using the insert or update statement depending on whether the ID is set or not.
// It will rollback the transaction and panic in case of an error.
func (connection *connection) SaveEntity(tx *sqlx.Tx, entity Entity, insert, update string) {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer Commit(tx)
	}

	var err error

	if entity.GetID() == 0 {
		var rows *sqlx.Rows
		rows, err = tx.NamedQuery(insert, entity)

		if err == nil {
			defer CloseRows(rows)
			rows.Next()
			var id int64

			if err := rows.Scan(&id); err != nil {
				Rollback(tx)
				logbuch.Fatal("Error scanning entity id", logbuch.Fields{"err": err, "entity": entity})
			}

			entity.SetID(id)
		}
	} else {
		_, err = tx.NamedExec(update, entity)
	}

	if err != nil {
		Rollback(tx)
		logbuch.Fatal("Error saving entity", logbuch.Fields{"err": err, "entity": entity})
	}
}
