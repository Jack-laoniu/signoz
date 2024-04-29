package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestXxx(t *testing.T) {
	now := time.Now()
	hackjwt := now.AddDate(100, 0, 0).Unix()
	fmt.Println(hackjwt)
}

func TestToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImphY2tsYW9uaXVAaWNsb3VkLmNvbSIsImV4cCI6MTcxMzk0MDE4OSwiZ2lkIjoiY2Q1NDkxMzctZDU2MC00ZGQxLTlkMTUtZjhmZjQ5MDY1NTIzIiwiaWQiOiI2M2RiMjlhYS1mZGEyLTQzY2ItODhjMi00NDhlODQwNjYzMGQifQ.J2CKtT_c72I4u23iPDJXSWXbOmbzlCiRkOwvYpQ_xso"
	// authHeaderParts := strings.Fields(token)
	// fmt.Println(strings.ToLower(authHeaderParts[0]))
	claims, err := ParseJWT(token)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(claims)
	exp, _ := claims["exp"]

	switch expType := exp.(type) {
	case float64:
		fmt.Println(int64(expType))
	case json.Number:
		v, _ := expType.Int64()
		fmt.Println(v)
	}
}

func TestENV(t *testing.T) {
	storage := os.Getenv("STORAGE")
	fmt.Println(storage, "123")
}
