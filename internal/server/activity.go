package server

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	api "github.com/Likebutter/todo-api/api/v1"
	_ "github.com/mattn/go-sqlite3"
)

const createTableQuery string = `
		CREATE TABLE IF NOT EXISTS activities (
		id INTEGER NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		time DATETIME NOT NULL
		);
`
const file string = "activities.db"

var ErrIDNotFound = errors.New("Id not found")

type Activities struct {
	db *sql.DB
}

func NewActivities() (*Activities, error) {
	db, err := sql.Open("sqlite3", file)

	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(createTableQuery); err != nil {
		return nil, err
	}

	return &Activities{
		db: db,
	}, nil
}

func (c *Activities) Insert(activity *api.Activity) (*api.Activity, error) {
	_, err := c.db.Exec("INSERT INTO activities VALUES(NULL,?,?,?);",
		activity.Name,
		activity.Description,
		activity.Time.AsTime(),
	)

	if err != nil {
		return nil, err
	}

	log.Printf("Added activity: %v", activity)

	return activity, nil
}

func (c *Activities) Get(id int) (*api.Activity, error) {
	log.Printf("Getting activity for id: %d", id)

	row := c.db.QueryRow("SELECT id, name, description, time FROM activities WHERE id=?", id)

	activity := api.Activity{}

	var err error
	var time time.Time

	if err = row.Scan(&activity.Id, &activity.Name, &activity.Description, &time); err == sql.ErrNoRows {
		log.Printf("Id not found")
		return &api.Activity{}, ErrIDNotFound
	}

	activity.Time = timestamppb.New(time)

	return &activity, err
}

func (c *Activities) List(offset int) ([]*api.Activity, error) {
	log.Printf("Getting list from offset %d\n", offset)

	rows, err := c.db.Query("SELECT * FROM activities WHERE ID > ? ORDER BY id DESC LIMIT 100", offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []*api.Activity{}
	for rows.Next() {
		i := api.Activity{}
		var time time.Time
		err = rows.Scan(&i.Id, &time, &i.Description)
		if err != nil {
			return nil, err
		}
		i.Time = timestamppb.New(time)
		data = append(data, &i)
	}
	return data, nil
}
