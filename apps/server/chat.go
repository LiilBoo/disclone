package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	jose "github.com/go-jose/go-jose/v3/jwt"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

type KeyManager interface {
	Validate(key string) bool
	Add() (string, error)
}

type chatServer struct {
	// subscriberMessageBuffer controls the max number
	// of messages that can be queued for a subscriber
	// before it is kicked.
	//
	// Defaults to 16.
	subscriberMessageBuffer int

	// publishLimiter controls the rate limit applied to the publish endpoint.
	//
	// Defaults to one publish every 100ms with a burst of 8.
	publishLimiter *rate.Limiter

	// logf controls where logs are sent.
	// Defaults to log.Printf.
	logf func(f string, v ...interface{})

	// Router routes the various endpoints to the appropriate handler.
	Router *chi.Mux

	subscribersMu sync.Mutex
	subscribers   map[*subscriber]struct{}

	keyManager KeyManager
}

func newChat(s *Server) *chatServer {
	//kr := NewKeyRetention()
	r := chi.NewRouter()

	cs := &chatServer{
		Router:                  r,
		subscriberMessageBuffer: 16,
		logf:                    log.Printf,
		subscribers:             make(map[*subscriber]struct{}),
		publishLimiter:          rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
		//keyManager:              kr,
	}

	cs.Router.Post("/login", WithAuth(cs.login, s.clerk))
	cs.Router.Get("/subscribe", cs.subscribeHandler)
	cs.Router.Post("/publish", cs.publishHandler)

	return cs
}

type subscriber struct {
	msgs      chan []byte
	closeSlow func()
}

func (cs *chatServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cs.Router.ServeHTTP(w, r)
}

type TestClaims struct {
	Subject string
}

func (cs *chatServer) login(w http.ResponseWriter, r *http.Request) {
	sessionClaims := r.Context().Value("session").(jose.Claims)

	fmt.Println(sessionClaims.Subject)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "disclone",
		"sub": sessionClaims.Subject,
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	})
	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "_chat",
		Value:    tokenString,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	w.Write([]byte("cookie set!"))
}

func (cs *chatServer) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_chat")

	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			fmt.Println("cookie not found")
			http.Error(w, "cookie not found", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	fmt.Println(cookie.Value)

	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})

	if err != nil {
		cs.logf("%v", err)
		return
	}

	defer c.Close(websocket.StatusInternalError, "")

	err = cs.subscribe(r.Context(), c)

	if errors.Is(err, context.Canceled) {
		return
	}

	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}

	if err != nil {
		cs.logf("%v", err)
		return
	}
}

func (cs *chatServer) publishHandler(w http.ResponseWriter, r *http.Request) {

	body := http.MaxBytesReader(w, r.Body, 8192)

	msg, err := ioutil.ReadAll(body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}

	cs.publish(msg)

	w.WriteHeader(http.StatusAccepted)
}

func (cs *chatServer) subscribe(ctx context.Context, c *websocket.Conn) error {
	ctx = c.CloseRead(ctx)

	s := &subscriber{
		msgs: make(chan []byte, cs.subscriberMessageBuffer),
		closeSlow: func() {
			c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
		},
	}

	cs.addSubscriber(s)
	defer cs.deleteSubscriber(s)

	for {
		select {
		case msg := <-s.msgs:
			err := writeTimeout(ctx, time.Second*5, c, msg)

			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (cs *chatServer) publish(msg []byte) {
	cs.subscribersMu.Lock()

	defer cs.subscribersMu.Unlock()

	cs.publishLimiter.Wait(context.Background())

	for s := range cs.subscribers {
		select {
		case s.msgs <- msg:
		default:
			go s.closeSlow()
		}
	}
}

func (cs *chatServer) addSubscriber(s *subscriber) {
	cs.subscribersMu.Lock()
	cs.subscribers[s] = struct{}{}
	cs.subscribersMu.Unlock()
}

func (cs *chatServer) deleteSubscriber(s *subscriber) {
	cs.subscribersMu.Lock()
	delete(cs.subscribers, s)
	cs.subscribersMu.Unlock()
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}
