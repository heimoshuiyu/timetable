package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type GetRangeRequest struct {
	Category string `json:"category"`
}

type Category struct {
	Name string `json:"name"`
}

type Token struct {
	Token string `json:"token"`
}

type Range struct {
	ID       int64  `json:"id"`
	Category string `json:"category"`
	Start    int64  `json:"start"`
	End      int64  `json:"end"`
	User     string `json:"user"`
}

type SetRangeUserRequest struct {
	Token string `json:"token"`
	Range Range  `json:"range"`
}

type AddRangeRequest struct {
	Token string `json:"token"`
	Range Range  `json:"range"`
}

type SetLimitRequest struct {
	Limit int    `json:"limit"`
	Token string `json:"token"`
}

type DeleteRangeRequest struct {
	Token string `json:"token"`
	ID    int64  `json:"id"`
}

type ClearRangeRequest struct {
	Token string `json:"token"`
	ID    int64  `json:"id"`
}

var Dbname string
var MyToken string

func init() {
	flag.StringVar(&Dbname, "dbname", "db.sqlite3", "database name")
	flag.StringVar(&MyToken, "token", "woshimima", "token")
}

func main() {
	flag.Parse()
	// global variables
	limit := 1
	dbWriteLock := sync.Mutex{}
	running := false

	// connect database
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		panic(err)
	}

	// create table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS range (
		id INTEGER PRIMARY KEY,
		category TEXT NOT NULL,
		start INTEGER NOT NULL,
		end INTEGER NOT NULL,
		user TEXT NOT NULL DEFAULT ''
	);`)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// server web folder http
	http.Handle("/", http.FileServer(http.Dir("./web/build")))

	// hanle hello api
	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})


	// handle get all range api
	http.HandleFunc("/api/range", func(w http.ResponseWriter, r *http.Request) {
		var err error
		var rows *sql.Rows
		log.Println("get range")
		if r.Method == "POST" {
			req := GetRangeRequest{}
			err = json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				HandleError(w, err)
				return
			}
			rows, err = db.Query("SELECT * FROM range WHERE category = ?", req.Category)
		} else {
			rows, err = db.Query("SELECT * FROM range ORDER BY category,start ASC")
		}
		if err != nil {
			HandleError(w, err)
			return
		}
		defer rows.Close()

		var ranges []Range = make([]Range, 0)
		for rows.Next() {
			var i Range
			err := rows.Scan(&i.ID, &i.Category, &i.Start, &i.End, &i.User)
			if err != nil {
				HandleError(w, err)
				return
			}
			ranges = append(ranges, i)
		}

		json.NewEncoder(w).Encode(ranges)
	})

	// handle add range api
	http.HandleFunc("/api/range/add", func(w http.ResponseWriter, r *http.Request) {
		var i AddRangeRequest
		err := json.NewDecoder(r.Body).Decode(&i)
		if err != nil {
			HandleError(w, err)
			return
		}

		if i.Token != MyToken {
			HandleError(w, errors.New("invalid token"))
			return
		}

		result, err := db.Exec("INSERT INTO range (category, start, end) VALUES (?, ?, ?)", i.Range.Category, i.Range.Start, i.Range.End)
		if err != nil {
			HandleError(w, err)
			return
		}

		i.Range.ID, err = result.LastInsertId()

		json.NewEncoder(w).Encode(i)
	})

	// handle set limit api
	http.HandleFunc("/api/setlimit", func(w http.ResponseWriter, r *http.Request) {
		var i SetLimitRequest
		err := json.NewDecoder(r.Body).Decode(&i)
		if err != nil {
			HandleError(w, err)
			return
		}

		if i.Token != MyToken {
			HandleError(w, errors.New("invalid token"))
			return
		}

		limit = i.Limit
		json.NewEncoder(w).Encode(i)
	})

	// handle set range's user api
	http.HandleFunc("/api/range/setuser", func(w http.ResponseWriter, r *http.Request) {
		var i SetRangeUserRequest
		err := json.NewDecoder(r.Body).Decode(&i)
		if err != nil {
			HandleError(w, err)
			return
		}

		log.Println("setuser", i.Range.User, i.Range.ID)

		if i.Token != MyToken && !running {
			HandleError(w, errors.New("还没到开始时间哦"))
			return
		}

		dbWriteLock.Lock()

		// check if user is exist
		username := ""
		err = db.QueryRow("SELECT user FROM range WHERE id = ?", i.Range.ID).Scan(&username)
		if err != nil {
			dbWriteLock.Unlock()
			HandleError(w, err)
			return
		}
		if username != "" && i.Token != MyToken {
			dbWriteLock.Unlock()
			HandleError(w, errors.New("这个时间段已经有人啦"))
			return
		}

		// count user in range
		count := 0
		db.QueryRow("SELECT COUNT(*) FROM range WHERE user = ? LIMIT 1", i.Range.User).Scan(&count)
		if count >= limit && i.Token != MyToken {
			dbWriteLock.Unlock()
			HandleError(w, errors.New("你的排班已经超过了限制"))
			return
		}

		_, err = db.Exec("UPDATE range SET user = ? WHERE id = ?", i.Range.User, i.Range.ID)
		if err != nil {
			dbWriteLock.Unlock()
			HandleError(w, err)
			return
		}

		dbWriteLock.Unlock()

		json.NewEncoder(w).Encode(i)
	})

	// handle delete range api
	http.HandleFunc("/api/range/delete", func(w http.ResponseWriter, r *http.Request) {
		var i DeleteRangeRequest
		err := json.NewDecoder(r.Body).Decode(&i)
		if err != nil {
			HandleError(w, err)
			return
		}

		if i.Token != MyToken {
			HandleError(w, errors.New("invalid token"))
			return
		}

		_, err = db.Exec("DELETE FROM range WHERE id = ?", i.ID)
		if err != nil {
			HandleError(w, err)
			return
		}

		json.NewEncoder(w).Encode(i)
	})

	// handle clear range api
	http.HandleFunc("/api/range/clear", func(w http.ResponseWriter, r *http.Request) {
		var i ClearRangeRequest
		err := json.NewDecoder(r.Body).Decode(&i)
		if err != nil {
			HandleError(w, err)
			return
		}

		if i.Token != MyToken {
			HandleError(w, errors.New("invalid token"))
			return
		}

		_, err = db.Exec("UPDATE range SET user = '' WHERE id = ?", i.ID)
		if err != nil {
			HandleError(w, err)
			return
		}

		json.NewEncoder(w).Encode(i)
	})

	// handle start request
	http.HandleFunc("/api/start", func(w http.ResponseWriter, r *http.Request) {
		var i Token
		err := json.NewDecoder(r.Body).Decode(&i)
		if err != nil {
			HandleError(w, err)
			return
		}

		if i.Token != MyToken {
			HandleError(w, errors.New("invalid token"))
			return
		}

		dbWriteLock.Lock()
		running = true
		dbWriteLock.Unlock()

		json.NewEncoder(w).Encode(i)
	})

	// handle stop request
	http.HandleFunc("/api/stop", func(w http.ResponseWriter, r *http.Request) {
		var i Token
		err := json.NewDecoder(r.Body).Decode(&i)
		if err != nil {
			HandleError(w, err)
			return
		}

		if i.Token != MyToken {
			HandleError(w, errors.New("invalid token"))
			return
		}

		dbWriteLock.Lock()
		running = false
		dbWriteLock.Unlock()

		json.NewEncoder(w).Encode(i)
	})

	// handle get cateogry api
	http.HandleFunc("/api/category", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT category FROM range GROUP BY category")
		if err != nil {
			HandleError(w, err)
			return
		}

		var categories []Category
		for rows.Next() {
			var c Category
			err = rows.Scan(&c.Name)
			if err != nil {
				HandleError(w, err)
				return
			}
			categories = append(categories, c)
		}

		json.NewEncoder(w).Encode(categories)
	})

	// start http server
	log.Println("start http server at :8080")
	http.ListenAndServe(":8080", nil)
}

func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
}
