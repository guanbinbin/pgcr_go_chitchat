package data

import (
	"time"
)

// 帖子
type Thread struct {
	Id        int       //id
	Uuid      string    //uuid
	Topic     string    //标题
	UserId    int       //用户
	CreatedAt time.Time //创建时间
}

type Post struct { //回复
	Id        int
	Uuid      string
	Body      string //内容
	UserId    int    //用户id
	ThreadId  int    //帖子id
	CreatedAt time.Time
}

// 格式化时间
func (thread *Thread) CreatedAtDate() string {
	return thread.CreatedAt.Format("2006-01-02 15:04:05")
}

func (post *Post) CreatedAtDate() string {
	return post.CreatedAt.Format("2006-01-02 15:04:05")
}

// 获取帖子的所有评论数量
func (thread *Thread) NumReplies() (count int) {
	// 执行查询
	rows, err := Db.Query("SELECT count(*) FROM posts where thread_id = $1", thread.Id)
	if err != nil {
		return
	}
	// 遍历结果，并赋值给count
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return
		}
	}
	// 关闭数据流
	rows.Close()
	return
}

// 获取一个帖子的所有回复
func (thread *Thread) Posts() (posts []Post, err error) {
	// 执行查询
	rows, err := Db.Query("SELECT id, uuid, body, user_id, thread_id, created_at FROM posts where thread_id = $1", thread.Id)
	if err != nil {
		return
	}
	// 封装查询结果
	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt); err != nil {
			return
		}
		posts = append(posts, post)
	}
	// 关闭数据流
	rows.Close()
	return
}

// 创建一篇新的帖子
func (user *User) CreateThread(topic string) (conv Thread, err error) {
	// 执行sql语句
	statement := "insert into threads (uuid, topic, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, topic, user_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	// 关闭数据流
	defer stmt.Close()
	// 执行创建，并将返回结果封装为Thread帖子对象并返回
	err = stmt.QueryRow(createUUID(), topic, user.Id, time.Now()).Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt)
	return
}

// 创建新的回复
func (user *User) CreatePost(conv Thread, body string) (post Post, err error) {
	// 创建sql语句
	statement := "insert into posts (uuid, body, user_id, thread_id, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, body, user_id, thread_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	// 关闭数据流
	defer stmt.Close()
	// 执行查询并将结果封装到评论Post对象并返回
	err = stmt.QueryRow(createUUID(), body, user.Id, conv.Id, time.Now()).Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
	return
}

// 获取数据库中所有的帖子
func Threads() (threads []Thread, err error) {
	// 执行sql语句
	rows, err := Db.Query("SELECT id, uuid, topic, user_id, created_at FROM threads ORDER BY created_at DESC")
	if err != nil {
		return
	}
	// 遍历每一行数据
	for rows.Next() {
		// 转换为帖子
		conv := Thread{}
		if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt); err != nil {
			return
		}
		// 追加
		threads = append(threads, conv)
	}
	// 关闭数据流
	rows.Close()
	return
}

// 根据UUID获取帖子
func ThreadByUUID(uuid string) (conv Thread, err error) {
	// 构建帖子对象
	conv = Thread{}
	// 执行查询并将结果封装到帖子对象
	err = Db.QueryRow("SELECT id, uuid, topic, user_id, created_at FROM threads WHERE uuid = $1", uuid).
		Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt)
	return
}

// 获取是谁创建了帖子
func (thread *Thread) User() (user User) {
	// 创建用户对象
	user = User{}
	// 根据用户id查询用户并封装到用户对象
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", thread.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

// 获取是谁回复的
func (post *Post) User() (user User) {
	// 创建用户对象
	user = User{}
	// 根据用户id查询并将查询结果封装到用户对象
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", post.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}
