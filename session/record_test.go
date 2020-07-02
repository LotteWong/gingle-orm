package session

import "testing"

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func testRecordInit(t *testing.T) *Session {
	t.Helper()
	s := NewSession().Model(&User{})
	dropErr := s.DropTable()
	createErr := s.CreateTable()
	affected, insertErr := s.Insert(user1, user2)
	if dropErr != nil || createErr != nil || insertErr != nil || affected != 1 {
		t.Fatal("Failed to initialize test records")
	}
	return s
}

func TestSession_Insert(t *testing.T) {
	s := testRecordInit(t)
	if affected, err := s.Insert(user3); err != nil || affected != 1 {
		t.Fatal("Failed to insert a record:", err)
	}
}

func TestSession_Find(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	if err := s.Find(&users); err != nil || len(users) != 2 {
		t.Fatal("Failed to query all records:", err)
	}
}

func TestSession_First(t *testing.T) {
	s := testRecordInit(t)
	user := &User{}
	if err := s.First(user); err != nil || user.Name != "Tom" || user.Age != 18 {
		t.Fatal("Failed to query the first record:", err)
	}
}

func TestSession_Where(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	affected, insertErr := s.Insert(user3)
	whereErr := s.Where("Age = ?", 25).Find(&users)
	if insertErr != nil || whereErr != nil || affected != 1 || len(users) != 2 {
		t.Fatal("Failed to query with where condition")
	}
}

func TestSession_Limit(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	if err := s.Limit(1).Find(&users); err != nil || len(users) != 1 {
		t.Fatal("Failed to query limited number of records:", err)
	}
}

func TestSession_OrderBy(t *testing.T) {
	s := testRecordInit(t)
	user := &User{}
	if err := s.OrderBy("Age DESC").First(user); err != nil || user.Age != 25 {
		t.Fatal("Failed to query with order by")
	}
}

func TestSession_Update(t *testing.T) {
	s := testRecordInit(t)
	affected, updateErr := s.Where("Name = ?", "Tom").Update("Age", 30)
	user := &User{}
	findErr := s.OrderBy("Age DESC").First(user)

	if updateErr != nil || findErr != nil || affected != 1 || user.Name != "Tom" || user.Age != 30 {
		t.Fatal("Failed to update a record")
	}
}

func TestSession_Delete(t *testing.T) {
	s := testRecordInit(t)
	affected, deleteErr := s.Where("Name = ?", "Tom").Delete()
	count, countErr := s.Count()
	if deleteErr != nil || countErr != nil || affected != 1 || count != 1 {
		t.Fatal("Failed to delete records")
	}
}

func TestSession_Count(t *testing.T) {
	s := testRecordInit(t)
	count, err := s.Count()
	if err != nil || count != 2 {
		t.Fatal("Failed to count the number of records:", err)
	}
}
