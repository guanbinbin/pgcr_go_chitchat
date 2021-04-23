package data

import (
	"time"
)

type User struct {
	Id        int
	Uuid      string
	Name      string //昵称
	Email     string //邮箱
	Password  string //密码
	CreatedAt time.Time
}

type Session struct {
	Id        int
	Uuid      string
	Email     string    //邮箱
	UserId    int       //用户id
	CreatedAt time.Time //创建时间
}

// 为已存在的用户创建一个新的session对象
func (user *User) CreateSession() (session Session, err error) {
	// 创建sql语句
	statement := "insert into sessions (uuid, email, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, email, user_id, created_at"
	// 执行sql语句
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	// 关闭处理器对象
	defer stmt.Close()
	// 取出查询到的数据，并存入到session中
	err = stmt.QueryRow(createUUID(), user.Email, user.Id, time.Now()).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	// 相当于 return session, err
	return
}

// 获取已经存在的用户的session对象
func (user *User) Session() (session Session, err error) {
	session = Session{}
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE user_id = $1", user.Id).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

// 检查数据库中的会话唯一ID是否存在
func (session *Session) Check() (valid bool, err error) {
	// 查询并赋值
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = $1", session.Uuid).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	// 不存在，返回false和err
	if err != nil {
		valid = false
		return
	}
	// 存在，返回true和nil
	if session.Id != 0 {
		valid = true
	}
	return
}

// 根据UUID删除session
func (session *Session) DeleteByUUID() (err error) {
	statement := "delete from sessions where uuid = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(session.Uuid)
	return
}

// 根据session获取用户信息
func (session *Session) User() (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", session.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

// 删除数据库中所有的session信息
func SessionDeleteAll() (err error) {
	statement := "delete from sessions"
	_, err = Db.Exec(statement)
	return
}

// 创建一个新用户
func (user *User) Create() (err error) {
	// 插入用户并指定要返回的数据
	statement := "insert into users (uuid, name, email, password, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// 执行插入并将返回值封装到用户信息中
	err = stmt.QueryRow(createUUID(), user.Name, user.Email, Encrypt(user.Password), time.Now()).Scan(&user.Id, &user.Uuid, &user.CreatedAt)
	return
}

// 从数据库中删除用户
func (user *User) Delete() (err error) {
	statement := "delete from users where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Id)
	return
}

// 更新用户信息
func (user *User) Update() (err error) {
	// 准备sql语句
	statement := "update users set name = $2, email = $3 where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// 执行更新
	_, err = stmt.Exec(user.Id, user.Name, user.Email)
	return
}

// 删除数据库中所有的用户
func UserDeleteAll() (err error) {
	statement := "delete from users"
	_, err = Db.Exec(statement)
	return
}

// 获取数据库中所有的用户
func Users() (users []User, err error) {
	// 查询数据
	rows, err := Db.Query("SELECT id, uuid, name, email, password, created_at FROM users")
	if err != nil {
		return
	}
	// 封装数据
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt); err != nil {
			return
		}
		users = append(users, user)
	}
	rows.Close()
	return
}

// 根据邮箱获取用户信息
func UserByEmail(email string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE email = $1", email).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

// 根据UUID获取用户信息
func UserByUUID(uuid string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = $1", uuid).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}
