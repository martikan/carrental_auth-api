package util

import (
	"testing"

	common "github.com/martikan/carrental_common/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {

	pass := common.RandomUtils.RandomString(6)

	hash, err := PasswordUtils.HashPassword(pass)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	err = PasswordUtils.CheckPassword(pass, hash)
	require.NoError(t, err)

	incorrectPassword := common.RandomUtils.RandomString(6)

	err = PasswordUtils.CheckPassword(incorrectPassword, hash)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hash2, err := PasswordUtils.HashPassword(pass)
	require.NoError(t, err)
	require.NotEmpty(t, hash2)
	require.NotEqual(t, hash, hash2)
}
