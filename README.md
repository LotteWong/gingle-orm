# gingle-orm

A simple xorm-like and gorm-like object relation mapping framework implemented by Golang.

---

## Features

- [x] Object Relation Mapping Based on Reflect

- [x] Table and Record CURD Support

- [x] Use Dialect to Adapt to Different Databases

- [x] Use Clause to Chain or Build SQLs

- [x] Hook Support

- [x] Transaction Support

- [x] Migration Support

## Quick Start

### Hook Sample

```go
type Account struct {
	ID       int `gingleorm:"PRIMARY KEY"`
	Password string
}

func (account *Account) BeforeInsert(s *Session) error {
	log.Infoln("Before Insert:", account)
	account.ID += 1000
	return nil
}

func (account *Account) AfterQuery(s *Session) error {
	log.Infoln("After Query:", account)
	account.Password = "******"
	return nil
}

func TestSession_CallMethod(t *testing.T) {
	s := NewSession().Model(&Account{})
	dropErr := s.DropTable()
	createErr := s.CreateTable()

	affected, insertErr := s.Insert(&Account{1, "123456"}, &Account{2, "clwong"})
	account := &Account{}
	firstErr := s.First(account)

	if dropErr != nil || createErr != nil || insertErr != nil || firstErr != nil {
		t.Fatal("Failed to call hooks before insert and after query")
	}
	if affected != 2 || account.ID != 1001 || account.Password != "******" {
		t.Fatal("Failed to call hooks before insert and after query")
	}
}
```

### Transaction Sample

```go
func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "gingleorm.db")
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	return engine
}

type User struct {
	Name string `gingleorm:"PRIMARY KEY"`
	Age int
}

func transactionRollback(t *testing.T) {
	e := OpenDB(t)
	defer e.Close()

	s := e.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := e.Transaction(func (s *session.Session) (res interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return nil, err
	})

	if err == nil || s.ExistTable() {
		t.Fatal("Failed to rollback")
	}
}

func transactionCommit(t *testing.T) {
	e := OpenDB()
	defer e.Close()

	s := e.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := e.Transaction(func (s *session.Session) (res interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return nil, err
	})

	user := &User{}
	_ = s.First(user)
	if err != nil || user.Name != "Tom" || user.Age != 18 {
		t.Fatal("Failed to commit")
	}
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commint", func(t *testing.T) {
		transactionCommit(t)
	})
}
```

### Migration Sample

```go
func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "gingleorm.db")
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	return engine
}

type User struct {
	Name string `gingleorm:"PRIMARY KEY"`
	Age int
}

func TestEngine_Migrate(t *testing.T) {
	e := OpenDB()
	defer e.Close()

	s := e.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text PRIMARY KEY, XXX integer);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) VALUES (?), (?)", "Tom", "Sam").Exec()
	e.Migrate(&User{})

	rows, _ := s.Raw("SELECT * FROM User").Query()
	cols, _ := rows.Columns()
	if !reflect.DeepEqual(columns, []string{"Name", "Age"}) {
		t.Fatal("Failed to migrate")
	}
}
```

## TODOs

- [ ] Complicated SQLs and Tags

- [ ] Adapt to More Databases *(right now sqlite3 only)*
- [ ] Read/Write Splitting Support
- [ ] Associated Relation Support
- [ ] Import/Export Data Support

