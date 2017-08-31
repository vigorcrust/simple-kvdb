package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

type Bucket struct {
	Name string `json:"bucket,omitempty"`
}

type KeyValue struct {
	Key    string `json:"key,omitempty"`
	Value  string `json:"value,omitempty"`
	Bucket string `json:"bucket,omitempty"`
}

func main() {
	const version = "v0.1"
	const apiIdent = "/api/"

	log.SetPrefix("[simple-kv] ")
	log.Println("Start simple-kv " + version)

	db, err := bolt.Open("data/skv.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/"+r.URL.Path[1:])
	})
	http.HandleFunc(apiIdent, func(w http.ResponseWriter, r *http.Request) {
		switch dict := strings.Split(strings.Replace(r.URL.Path, apiIdent, "", 1), "/"); len(dict) {
		case 1:
			if dict[0] == "" {
				fmt.Fprint(w, getBuckets(db))
			} else {
				switch r.Method {
				case "GET":
					fmt.Fprint(w, getKeys(db, dict[0]))
				case "POST":
					createBucket(db, dict[0])
				case "DELETE":
					deleteBucket(db, dict[0])
				}
			}
		case 2:
			switch r.Method {
			case "GET":
				fmt.Fprint(w, getValue(db, dict[0], dict[1]))
			case "DELETE":
				deleteKey(db, dict[0], dict[1])
			}
		case 3:
			if r.Method == "POST" {
				createKey(db, dict[0], dict[1], dict[2])
			}
		default:
			fmt.Fprint(w, "unknown amount of parameters")
		}
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
}
func createBucket(db *bolt.DB, bucketname string) {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketname))
		if err != nil {
			log.Printf("create bucket: %s", err)
		}
		return nil
	})
}
func getBuckets(db *bolt.DB) string {
	var buckets []Bucket
	db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			buckets = append(buckets, Bucket{Name: string(name)})
			return nil
		})
	})
	json, _ := json.Marshal(buckets)
	if nil == json || "null" == string(json) {
		return "{}"
	}
	return string(json)
}
func deleteBucket(db *bolt.DB, bucketname string) {
	db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(bucketname))
		if err != nil {
			log.Printf("delete bucket: %s", err)
		}
		return nil
	})
}
func createKey(db *bolt.DB, bucketname string, keyname string, value string) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketname))
		err := b.Put([]byte(keyname), []byte(value))
		return err
	})
}
func getKeys(db *bolt.DB, bucketname string) string {
	var kv []KeyValue
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketname))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			kv = append(kv, KeyValue{Key: string(k), Value: string(v)})
		}
		return nil
	})
	json, _ := json.Marshal(kv)
	if nil == json || "null" == string(json) {
		return "{}"
	}
	return string(json)
}
func deleteKey(db *bolt.DB, bucketname string, keyname string) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketname))
		b.Delete([]byte(keyname))
		return nil
	})
}
func getValue(db *bolt.DB, bucketname string, keyname string) string {
	var kv KeyValue
	var v []byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketname))
		if b == nil {
			return nil
		}
		v = b.Get([]byte(keyname))
		if v == nil {
			return nil
		}
		kv = KeyValue{Value: string(v)}
		return nil
	})
	json, _ := json.Marshal(kv)
	if nil == json || "null" == string(json) {
		return ""
	}
	return string(json)
}
