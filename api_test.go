package servrlinkapi

import "testing"

func TestIsRegistered(t *testing.T) {
	mcIsRegistered, err := IsRegistered("c3e4a469-2e9d-4cb1-be1b-80fedf40e71b")
	if err != nil {
		t.Error(err)
	} else if !mcIsRegistered {
		t.Error("Apparently Minecraft account isn't registered.")
	}

	discordIsRegistered, err := IsRegistered("217617036749176833")
	if err != nil {
		t.Error(err)
	} else if !discordIsRegistered {
		t.Error("Apparently Discord account isn't registered.")
	}
}

func TestGet(t *testing.T) {
	id, err := Get("c3e4a469-2e9d-4cb1-be1b-80fedf40e71b")
	if err != nil {
		t.Error(err)
	} else if id != "217617036749176833" {
		t.Error("invalid ID returned: ", id)
	}

	uuid, err := Get("217617036749176833")
	if err != nil {
		t.Error(err)
	} else if uuid != "c3e4a469-2e9d-4cb1-be1b-80fedf40e71b" {
		t.Error("invalid UUID returned: ", uuid)
	}

}
