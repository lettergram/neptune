package user

import(
	"os"
	"fmt"
	"net/http"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"io/ioutil"
	"strings"

	"neptune/pkgs/cookies"
)

type User struct {
	Email     string
	Password  string
	SessionID string
}

// Creates an account and adds it to the Database
func CreateAccount(usr *User) bool {

	session, err := mgo.Dial("127.0.0.1:27017/")
	if err != nil {
		fmt.Println(err)
		return false
	}

	c := session.DB("test").C("users")
	err = c.Insert(*usr)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// Checks if an account exists in the userbase.
func DoesAccountExist(email string) bool {

	session, err := mgo.Dial("127.0.0.1:27017/")
	if err != nil {
		panic(err)
	}

	result := User{}
	c := session.DB("test").C("users")

	err = c.Find(bson.M{"email": email}).One(&result)
	if err != nil {
		return false
	}

	return true
}

// Checks to assure credientials
func CheckCredentials(email string, password string) bool {

	session, err := mgo.Dial("127.0.0.1:27017/")
	if err != nil {
		panic(err)
	}

	result := User{}
	c := session.DB("test").C("users")

	err = c.Find(bson.M{"email": email}).One(&result)
	if err != nil {
		return false
	}

	if result.Password == password && result.Email == email {
		return true
	}

	return false
}

// Updates the user's data in the database (e.g. cookie data).
func UpdateUser(usr *User) bool {
	session, err := mgo.Dial("127.0.0.1:27017/")
	if err != nil {
		fmt.Println(err)
		return false
	}

	c := session.DB("test").C("users")
	err = c.Update(bson.M{"email": usr.Email}, usr)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// Loads users info - or supplys links to login or register 
func LoadUserInfo(title string, r *http.Request)(filename string, option []byte, usr []byte){

	if cookies.IsLoggedIn(r) {
		cookie, _ := r.Cookie("SessionID")
		z := strings.Split(cookie.Value, ":")
		filename = "accounts/" + z[0] + ".txt"
		usr, _ = ioutil.ReadFile(filename)
		option = []byte("<a href='/logout'>logout</a>")
	} else {
		option = []byte("<a href='/login'>login</a> or <a href='/register'>register</a>")
	}

	filename = "web/" + title + ".txt"

	return filename, option, usr

}

func CreateUserFile(usrName string){

	file, err := os.Create("accounts/" + usrName + ".txt")			// creates a file with that usrName
	if err != nil {	fmt.Printf("error createUserFile FIX")  }

	s := usrName + ": "

	file.WriteString(s)

}

func ReadUserFile(usrName string)(file *os.File){

	file, err := os.Open(usrName)
	if err != nil {	fmt.Printf("error readUserFile FIX")  }

	return file

}

func AppendUserFile(usrName string){

	file, err:= os.OpenFile(usrName, os.O_WRONLY, 0666)
	if err != nil {	fmt.Printf("error appendUserFile FIX")  }

	s := "nothing"

	file.WriteString(s)
	
}
