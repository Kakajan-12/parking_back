package core_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"backend/core"
	"backend/test"
)

func TestMain(m *testing.M) {
	test.Init()
	code := m.Run()
	os.Exit(code)
}

func TestJWTEncodeDecode(t *testing.T) {
	// Prepare payload
	payload := core.JWTPayload(map[string]interface{}{
		"user_id": "12345",
	}, nil, nil)

	// Encode
	tokenStr, err := core.JWTEncode(payload)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)

	// Decode
	claims, err := core.JWTDecode(tokenStr)
	require.NoError(t, err)
	require.Equal(t, "12345", claims["user_id"])
	require.Equal(t, "backend", claims["iss"])
	require.Equal(t, "client", claims["aud"])
}

func TestJWTDecode_ExpiredToken(t *testing.T) {

	// Create a payload that expired 1 second ago
	expDelta := -1 * time.Second
	payload := core.JWTPayload(map[string]interface{}{
		"user_id": "12345",
	}, nil, &expDelta)

	tokenStr, err := core.JWTEncode(payload)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)

	// Decode should fail due to expiration
	_, err = core.JWTDecode(tokenStr)
	require.Error(t, err)

	// Check that it's specifically an expiration error
	require.Contains(t, err.Error(), "token is expired")
}

func TestJWTDecode_ValidToken(t *testing.T) {

	// Token valid for 1 minute
	expDelta := 1 * time.Minute
	payload := core.JWTPayload(map[string]interface{}{
		"user_id": "12345",
	}, nil, &expDelta)

	tokenStr, err := core.JWTEncode(payload)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)

	claims, err := core.JWTDecode(tokenStr)
	require.NoError(t, err)
	require.Equal(t, "12345", claims["user_id"])
}

func TestJWTDecode_InvalidToken(t *testing.T) {

	// invalid token string
	_, err := core.JWTDecode("not.a.token")
	require.Error(t, err)
}
