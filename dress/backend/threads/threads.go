package threads

import "github.com/ss49919201/ai-kata/dress/backend/database/mysql"

type Post struct {
	ID      int
	Content string
}

type Thread struct {
	ID    int
	Title string
	Posts []*Post
}

func CreateThread(title string) (*Thread, error) {
	mysqlClient, err := mysql.NewClient(
		"root",
		"password",
		"localhost:3306",
		"dress",
	)
	if err != nil {
		return nil, err
	}

	defer mysqlClient.Close()

	result, err := mysqlClient.Exec(
		"INSERT INTO threads (title) VALUES (?)",
		title,
	)
	if err != nil {
		return nil, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Thread{
		ID:    int(lastInsertId),
		Title: title,
		Posts: []*Post{},
	}, nil
}

func GetThread(id int) (*Thread, error) {
	mysqlClient, err := mysql.NewClient(
		"root",
		"password",
		"localhost:3306",
		"dress",
	)
	if err != nil {
		return nil, err
	}

	defer mysqlClient.Close()

	// スレッドを1件取得
	// 合わせて、スレッドの投稿を取得する

	rows, err := mysqlClient.Query(
		"SELECT threads.id, threads.title, posts.id, posts.content FROM threads LEFT JOIN posts ON threads.id = posts.thread_id WHERE threads.id = ?",
		id,
	)
	if err != nil {
		return nil, err
	}

	thread := &Thread{}
	for rows.Next() {
		var post *Post

		if err := rows.Scan(
			&thread.ID,
			&thread.Title,
			&post.ID,
			&post.Content,
		); err != nil {
			return nil, err
		}

		thread.Posts = append(thread.Posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return thread, nil
}

func CreatePost(threadId int, content string) (*Post, error) {
	mysqlClient, err := mysql.NewClient(
		"root",
		"password",
		"localhost:3306",
		"dress",
	)
	if err != nil {
		return nil, err
	}

	defer mysqlClient.Close()

	result, err := mysqlClient.Exec(
		"INSERT INTO posts (thread_id, content) VALUES (?, ?)",
		threadId,
		content,
	)
	if err != nil {
		return nil, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Post{
		ID:      int(lastInsertId),
		Content: content,
	}, nil
}
