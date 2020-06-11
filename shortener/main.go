package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v2"
)

var (
	cfgFile  string
	redisCfg string
)

var globalLinks map[string]string = map[string]string{
	"/hello": "http://www.example.com",
}

type (
	Redirect struct {
		Path string `yaml:"path"`
		Rdr  string `yaml:"redirect"`
	}

	Redirects []Redirect

	Lookup interface {
		Get(string) (string, error)
		Put(string, string) error
	}

	LookupMap map[string]string

	RedisLookup struct {
		cli *redis.Client
	}
)

func (m LookupMap) Get(path string) (string, error) {
	target, ok := m[path]
	if !ok {
		return "", fmt.Errorf("invalid path: %s", path)
	}
	return target, nil
}

func (m LookupMap) Put(path, val string) error {
	m[path] = val
	return nil
}

func (r *RedisLookup) Get(path string) (string, error) {
	cmd := r.cli.Get(context.Background(), path)
	if cmd.Err() != nil {
		fmt.Println("got err:", cmd.Err().Error())
		return "", cmd.Err()
	}
	fmt.Printf("got from redis: '%s'\n", cmd.Val())
	return cmd.Val(), nil
}

func (r *RedisLookup) Put(path, val string) error {
	cmd := r.cli.Set(context.Background(), path, val, 0)
	if cmd.Err() != nil {
		fmt.Println("got err:", cmd.Err().Error())
		return cmd.Err()
	}
	return nil
}

func putHandler(lk Lookup, w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	val, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = lk.Put(path, string(val))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func LookupHandler(lk Lookup, fallback http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			putHandler(lk, w, r)
			return
		}

		target, err := lk.Get(r.URL.Path)
		if err != nil {
			fallback(w, r)
			return
		}
		http.Redirect(w, r, target, http.StatusTemporaryRedirect)
	}
}

func YAMLHandler(cfg string, fallback http.HandlerFunc) (http.HandlerFunc, error) {
	f, err := os.Open(cfg)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var rdrs Redirects
	err = yaml.NewDecoder(f).Decode(&rdrs)
	if err != nil {
		return nil, err
	}

	links := make(map[string]string)
	for _, v := range rdrs {
		links[v.Path] = v.Rdr
	}

	return LookupHandler(LookupMap(links), fallback), nil
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("fallback: not found!"))
}

func init() {
	flag.StringVar(&cfgFile, "config", "", "path to link configuration yaml")
	flag.StringVar(&redisCfg, "redis", "", "hostname and port to redis")
}

func main() {
	flag.Parse()

	var handler http.HandlerFunc
	var err error
	if cfgFile != "" {
		handler, err = YAMLHandler(cfgFile, notFoundHandler)
		if err != nil {
			log.Fatalf("failed to parse yaml: %s", err)
		}
	} else if redisCfg != "" {
		r := &RedisLookup{
			cli: redis.NewClient(&redis.Options{
				Addr:     redisCfg,
				Password: "", // no password set
				DB:       0,  // use default DB
			}),
		}
		handler = LookupHandler(r, notFoundHandler)
	} else {
		handler = LookupHandler(LookupMap(globalLinks), notFoundHandler)
	}

	http.HandleFunc("/", handler)
	fmt.Println("Server running on 8080...")
	log.Fatalf("server %s:", http.ListenAndServe(":8080", nil))
}
