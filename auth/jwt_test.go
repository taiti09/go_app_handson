package auth

import (
	"bytes"
	"context"
	"testing"

	"github.com/taiti09/go_app_handson/clock"
	"github.com/taiti09/go_app_handson/entity"
	"github.com/taiti09/go_app_handson/testutil/fixture"
)

func TestEmbed(t *testing.T) {
	want := []byte("------------BEGIN PUBLIC KEY------------")
	if !bytes.Contains(rawPubKey,want) {
		t.Errorf("want %s, but got %s", want, rawPubKey)
	}
}

func TestJWTer_GenerateToken(t *testing.T) {
	ctx := context.Background()
	moq := &StoreMock{}
	wantID := entity.UserID(20)
	u := fixture.User(&entity.User{ID: wantID})
	moq.SaveFunc = func(ctx context.Context, key string, userID entity.UserID) error {
		if userID != wantID {
			t.Errorf("want %d, but got %d", wantID,userID)
		}
		return nil
	}

	sut, err := NewJWTer(moq,clock.RealClocker{})
	if err != nil {
		t.Fatal(err)
	}
	got, err := sut.GenerateToken(ctx,*u)
	if err != nil {
		t.Fatalf("not want err: %v", err)
	}
	if len(got) == 0 {
		t.Errorf("token is empty")
	}
}