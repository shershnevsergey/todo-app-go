package repository

import (
	"fmt"
	"github.com/Shv-sergey70/todo-app-go"
	"github.com/jmoiron/sqlx"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var listId int
	createListQuery := fmt.Sprintf("INSERT INTO %s(title, description) VALUES($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&listId); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s(user_id, list_id) VALUES($1, $2)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, listId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return listId, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList
	query := fmt.Sprintf("SELECT tl.* FROM %s tl INNER JOIN %s ul ON tl.id = ul.list_id WHERE ul.user_id = $1", todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList
	query := fmt.Sprintf(`SELECT tl.* FROM %s tl 
								INNER JOIN %s ul ON tl.id = ul.list_id WHERE tl.id = $1 AND ul.user_id = $2`,
		todoListsTable,
		usersListsTable)
	err := r.db.Get(&list, query, listId, userId)

	return list, err
}

func (r *TodoListPostgres) Delete(userId, listId int) error {
	removeUserListQuery := fmt.Sprintf(`DELETE FROM %s tl USING %s ul 
		WHERE tl.id = ul.list_id AND ul.list_id = $1 AND ul.user_id = $2`,
		todoListsTable, usersListsTable)

	_, err := r.db.Exec(removeUserListQuery, listId, userId)

	return err
}
