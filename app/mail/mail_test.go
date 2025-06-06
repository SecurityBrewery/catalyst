package mail

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gm "github.com/wneessen/go-mail"
)

type fakeClient struct {
	msgs []*gm.Msg
	err  error
}

func (f *fakeClient) DialAndSend(msgs ...*gm.Msg) error {
	f.msgs = append(f.msgs, msgs...)

	return f.err
}

func withFakeClient(fc *fakeClient, errNew error) func() {
	original := newSMTPClient
	newSMTPClient = func(_ string, _ ...gm.Option) (smtpClient, error) {
		if errNew != nil {
			return nil, errNew
		}

		return fc, nil
	}

	return func() { newSMTPClient = original }
}

func TestMailer_Send(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	cfg := &Config{SMTPServer: "localhost"}
	m := New(cfg)

	t.Run("invalid from", func(t *testing.T) {
		t.Parallel()

		err := m.Send(ctx, "invalid", "to@example.com", "sub", "body")
		require.Error(t, err)
	})

	t.Run("invalid to", func(t *testing.T) {
		t.Parallel()

		err := m.Send(ctx, "from@example.com", "invalid", "sub", "body")
		require.Error(t, err)
	})

	t.Run("client creation failure", func(t *testing.T) {
		t.Parallel()

		restore := withFakeClient(nil, errors.New("new client error"))
		defer restore()

		err := m.Send(ctx, "from@example.com", "to@example.com", "sub", "body")
		require.Error(t, err)
	})

	t.Run("send failure", func(t *testing.T) {
		t.Parallel()

		fc := &fakeClient{err: errors.New("dial fail")}
		restore := withFakeClient(fc, nil)

		defer restore()

		err := m.Send(ctx, "from@example.com", "to@example.com", "sub", "body")
		require.Error(t, err)
		assert.Len(t, fc.msgs, 1)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		fc := &fakeClient{}
		restore := withFakeClient(fc, nil)

		defer restore()

		err := m.Send(ctx, "from@example.com", "to@example.com", "subject", "hello")
		require.NoError(t, err)
		assert.Len(t, fc.msgs, 1)

		msg := fc.msgs[0]
		assert.Equal(t, []string{"<from@example.com>"}, msg.GetFromString())
		assert.Equal(t, []string{"<to@example.com>"}, msg.GetToString())

		sub := msg.GetGenHeader(gm.HeaderSubject)
		assert.Len(t, sub, 1)
		assert.Equal(t, "subject", sub[0])

		parts := msg.GetParts()
		assert.Len(t, parts, 1)
		b, getErr := parts[0].GetContent()
		require.NoError(t, getErr)
		assert.Equal(t, "hello", string(b))
	})
}
