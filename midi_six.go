package l_heure_de_catherine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	functions.CloudEvent("MidiSix", MidiSix)
}

// a json file like this:
//
//	{
//	  "identifier": "a bluesky handle",
//	  "password": "(app) password of the account",
//	}
const credentialPath = "/mnt/credential"

func MidiSix(ctx context.Context, _ event.Event) error {
	credential, err := os.ReadFile(credentialPath)
	if err != nil {
		return fmt.Errorf("reading credential: %w", err)
	}

	var session struct {
		AccessJwt string `json:"accessJwt"`
		Handle    string `json:"handle"`
	}
	err = call(ctx, "com.atproto.server.createSession", "", credential, &session)
	if err != nil {
		return fmt.Errorf("creating the session: %w", err)
	}

	defer func() {
		call(ctx, "com.atproto.server.deleteSession", session.AccessJwt, nil, nil)
	}()

	record := struct {
		Repo       string `json:"repo"`
		Collection string `json:"collection"`
		Record     struct {
			Type      string    `json:"$type"`
			Langs     []string  `json:"langs"`
			Text      string    `json:"text"`
			CreatedAt time.Time `json:"createdAt"`
		} `json:"record"`
	}{
		Repo:       session.Handle,
		Collection: "app.bsky.feed.post",
	}
	record.Record.Type = "app.bsky.feed.post"
	record.Record.Langs = []string{"fr"}
	record.Record.Text = "C'est l'heure de Catherine!"
	record.Record.CreatedAt = getTime()

	// could be used to overwrite the default value
	// err = e.DataAs(&record.Record)
	// if err != nil {
	// 	return fmt.Errorf("unmarshaling event data: %w", err)
	// }

	raw, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshaling the record: %w", err)
	}

	err = call(ctx, "com.atproto.repo.createRecord", session.AccessJwt, raw, nil)
	if err != nil {
		return fmt.Errorf("posting the record: %w", err)
	}

	return nil
}

func call(ctx context.Context, endpoint, bearerToken string, payload []byte, body any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://bsky.social/xrpc/"+endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("creating the request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	if bearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+bearerToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("calling endpoint: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errInfo struct {
			Error   string `json:"error"`
			Message string `json:"Message"`
		}

		err = json.Unmarshal(raw, &errInfo)
		if err != nil {
			return fmt.Errorf("unmarshaling errInfo with unexpected status: %d: %w", resp.StatusCode, err)
		}

		return fmt.Errorf("unexpected status: %d, error: %s, message: %s", resp.StatusCode, errInfo.Error, errInfo.Message)
	}

	if body == nil {
		return nil
	}

	err = json.Unmarshal(raw, &body)
	if err != nil {
		return fmt.Errorf("unmarshaling body: %w", err)
	}
	return nil
}

func getTime() time.Time {
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		// should not happen as time/tzdata is imported
		panic(err.Error())
	}
	n := time.Now()
	return time.Date(n.Year(), n.Month(), n.Day(), 12, 6, 0, 0, loc).UTC()
}
