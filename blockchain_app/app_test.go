package blockchain_app

import "testing"

func TestDecodePayload(t *testing.T) {
	//malformed payload we should get an error
	payload := []byte("eyJvcGVyYXRpb24iOiAi2hhbmdlQ2FyT3duZXIiLCAiZGF0YSI6IHsiYXNzZXRfaWQiOiAiY2FyMSIsICJuZXdfb3duZXIiiAiUm9iZXJ0In19")

	_, err := decodePayload(payload)

	if err == nil {
		t.Error("We should not be able to decode this payload, we have payload ", payload)
	}

	//ok payload, we should be able to decode it
	payload = []byte("eyJvcGVyYXRpb24iOiAiY3JlYXRlQ2FyIiwgImRhdGEiOiB7IklEIjogImNhcjEiLCAiTWFrZSI6ICJQZXVnZW90IiwgIk1vZGVsIjogIjIwNSIsICJDb2xvdXIiOiAicmVkIiwgIk93bmVyIjogIkphbmUifX0=")

	transaction, err := decodePayload(payload)

	if err != nil {
		t.Error("We should able to decode this payload, we have payload", payload)
	}

	if transaction.Operation != OpCreateCar {
		t.Errorf("The operation should have been %s, we got %s", OpCreateCar, transaction.Operation)
	}
}
