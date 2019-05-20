package main


import (
	"fmt"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)


type Person struct {
	Id bson.ObjectId `json:"id" bson:"_id"`
	Name  string
	Phone string
}

func main() {
	session, err := mgo.Dial("114.115.200.40:17017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("oss").C("Person")

	var result []Person
	err = c.Find(bson.M{"_id":bson.ObjectIdHex("5cc80a34c3666e0da9116781")}).All(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Phone:", result)
}