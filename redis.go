package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

//	type User struct {
//		Name     string
//		Password string
//	}
type Store struct {
	data map[string]KeyValue
	//在 map[string]string 中：
	//
	//string：表示键的类型是字符串。
	//string：表示值的类型也是字符串。
	mu   sync.RWMutex
	sets map[string]map[stri ng]bool
	//user map[string]User
}
type KeyValue struct {
	Value string
	ET    time.Time
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]KeyValue),
		sets: make(map[string]map[string]bool),
		//user: make(map[string]User),
	}
}

func (s *Store) Load(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &s.data)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) Save(filename string) error {
	file, err := json.MarshalIndent(s.data, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, file, 0644)
}
func (s *Store) Set(key, value string, Ex int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ET := time.Now().Add(time.Duration(Ex) * time.Second)
	s.data[key] = KeyValue{value, ET}
}
func (s *Store) SetNx(key, Value string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.data[key]; exists {
		return 0
	}
	s.data[key] = KeyValue{Value: Value}
	return 1
}
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	kv, exists := s.data[key]
	if !exists {
		return "", false
	}
	if time.Now().After(kv.ET) {
		delete(s.data, key) // 删除过期的键
		fmt.Println("Key has expired and has been deleted:", key)
		return "", false
	}

	return kv.Value, true // 返回值
}
func (s *Store) Del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}
func (s *Store) SADD(setName string, value string) int {
	if _, exists := s.sets[setName]; !exists {
		s.sets[setName] = make(map[string]bool)
	}
	if s.sets[setName][value] {
		return 0
	}
	s.sets[setName][value] = true
	return 1
}

//	func HashPassWord(password string) (string, bool) {
//		bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//		return string(bytes), err
//	}
//
//	func (s *Store) Register(username, password string) {
//		hash, err := HashPassWord(password)
//		if err != nil {
//			return err
//		}
//		s.user[username] = User{Name: username, Password: hash}
//		return nil
//	}
//
//	func (s *Store) Login(username, password string) bool {
//		user, exists := s.user[username]
//		if !exists {
//			return false
//		}
//		return VerifyPassword(user.Password, password)
//	}
//
// 靠这个软件没有bcrypt先放放
func showMenu() {
	fmt.Println("1.命令行")
	fmt.Println("2.程序说明")
	fmt.Println("3.退出程序")
}
func showhelp() {
	fmt.Println("1.就是让你个ldx输入点东西实现简单redis")
	fmt.Println("2.self")
	fmt.Println("3.看如此屎的代码看不下去了就溜走")
}
func main() {
	store := NewStore()
	const filename = "store.json"
	err := store.Load(filename)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading data:", err)
		return
	}
	for {
		showMenu()
		var choice int
		fmt.Scan(&choice)
		switch choice {
		case 1:
			for {
				var command string
				var key, value string
				fmt.Print(">")
				n, err := fmt.Scanf("%s %s %s", &command, &key, &value)
				if command == "EXIT" {
					break
				}
				if err := store.Save(filename); err != nil {
					fmt.Println("Error saving data:", err)
				}
				if err != nil && err.Error() != "expected newline" {
					fmt.Println("Error reading command:", err)
					continue
				}
				if n == 0 {
					fmt.Println("No input recived.")
					continue
				}
				switch command {
				case "SET":
					store.Set(key, value, 5)
					store.Save(filename)
					fmt.Println("ok了")
				case "SETNX":
					result := store.SetNx(key, value)
					store.Save(filename)
					if result == 0 {
						fmt.Println("0")
					} else {
						fmt.Println("1")
					}
				case "GET":
					time.Sleep(3 * time.Second)
					if val, exists := store.Get(key); exists {
						fmt.Println(val)
					} else {
						fmt.Println("not found")
					}
					time.Sleep(3 * time.Second)
					if val, exists := store.Get(key); exists {
						fmt.Println(val) // 不会输出
					} else {
						fmt.Println("(nil or expired)") // 应该输出此行
					}
				case "DEL":
					store.Del(key)
					store.Save(filename)
					fmt.Println("ok")
				default:
					fmt.Println("Unknown command")
				case "SADD":
					Result1 := store.SADD(key, value)
					fmt.Println(Result1)
				case "EXIT":
					break
				}
			}
		case 2:
			showhelp()
		case 3:
			fmt.Println("Exiting program.")
			return
		}
	}
}
