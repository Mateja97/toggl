package toggl

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func (t *Toggl) dbGetAllQuestions(ctx context.Context) ([]Question, error) {

	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return []Question{}, err
	}

	questions := make([]Question, 0)
	rows, err := t.db.Query("SELECT id, body FROM question")
	if err != nil {
		return []Question{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var question Question
		if err := rows.Scan(&question.ID, &question.Body); err != nil {
			return []Question{}, err
		}

		options, err := t.dbGetOptionsByQuestion(question.ID)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Println("[ERROR]dbGetAllQuestions unsuccessful rollback error: ", err)
			}
			return []Question{}, err
		}
		if len(options) == 0 {
			continue
		}
		question.Options = options

		questions = append(questions, question)
	}
	return questions, tx.Commit()
}
func (t *Toggl) dbGetQuestionsByID(ctx context.Context, question RequestQuestion) ([]Question, error) {
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return []Question{}, err
	}

	questions := make([]Question, 0)
	rows, err := t.db.Query(`SELECT id, body 
							FROM question where id < $1 
							ORDER BY id DESC
							LIMIT $2`, question.offset, question.Limit)
	if err != nil {
		return []Question{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var question Question
		if err := rows.Scan(&question.ID, &question.Body); err != nil {
			return []Question{}, err
		}

		options, err := t.dbGetOptionsByQuestion(question.ID)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Println("[ERROR]dbGetAllQuestions unsuccessful rollback error: ", err)
			}
			return []Question{}, err
		}
		if len(options) == 0 {
			continue
		}
		question.Options = options

		questions = append(questions, question)
	}
	return questions, tx.Commit()
}

func (t *Toggl) dbGetOptionsByQuestion(questionID int64) ([]Option, error) {
	options := make([]Option, 0)
	rows, err := t.db.Query("SELECT body,correct FROM option where questionid = $1", questionID)
	if err != nil {
		return []Option{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var option Option
		if err := rows.Scan(&option.Body, &option.Correct); err != nil {
			return []Option{}, err
		}
		options = append(options, option)
	}
	return options, nil
}

func (t *Toggl) dbInsertQuestion(ctx context.Context, question Question) (Question, error) {
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return Question{}, err
	}
	result, err := t.db.Exec("INSERT INTO question (body) VALUES (?)", question.Body)
	if err != nil {
		return Question{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Question{}, err
	}
	question.ID = id
	for _, option := range question.Options {
		err := t.dbInsertOption(option, id)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Println("[ERROR]dbInsertQuestion unsuccessful rollback error: ", err)

			}
			return Question{}, err
		}
	}
	return question, tx.Commit()
}

func (t *Toggl) dbInsertOption(option Option, qID int64) error {
	_, err := t.db.Exec("INSERT INTO option (questionid, body, correct) VALUES (?,?,?)", qID, option.Body, option.Correct)
	if err != nil {
		return err
	}

	return nil
}

func (t *Toggl) dbUpdateQuestion(ctx context.Context, question Question, ID int64) (Question, error) {
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return Question{}, err
	}
	_, err = t.db.Exec("UPDATE question SET body = ? WHERE id = ?", question.Body, ID)
	if err != nil {
		return Question{}, err
	}
	_, err = t.db.Exec("DELETE FROM option WHERE questionid = ?", ID)
	if err != nil {
		return Question{}, err
	}

	for _, option := range question.Options {
		err := t.dbInsertOption(option, ID)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Println("[ERROR]dbUpdateQuestion unsuccessful rollback error: ", err)

			}
			log.Println("[ERROR]dbUpdateQuestion dbInsertOption error: ", err)
			return Question{}, err

		}
	}
	question.ID = ID
	return question, tx.Commit()
}

func (t *Toggl) dbDeleteQuestion(ctx context.Context, ID int64) (int64, error) {
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}
	res, err := t.db.Exec("DELETE FROM question WHERE id = ?", ID)
	if err != nil {
		return 0, err
	}
	resNum, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	_, err = t.db.Exec("DELETE FROM option WHERE questionid = ?", ID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Println("[ERROR]dbDeleteQuestion unsuccessful rollback error: ", err)
		}
		return 0, err
	}
	return resNum, tx.Commit()
}
