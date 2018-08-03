package nbum

import "testing"

func TestUserSignAuth_Sign(t *testing.T) {
	usa := &UserSignAuth{
		SecretKey: "h8o3T6iVwae4svff99P462mgtbWqCRu3",
		UserId:    1,
	}

	params := map[string]string{
		"operate":   "test",
		"operateId": "3",
		"userId":    "1",
		"ts":        "1531446253",
	}
	rawSign := "e0d4b774667fbf5f400cbc3a17ad5101"
	sign := usa.Sign(params)
	if sign != rawSign {
		t.Fatalf("expect sign %s, but got %s", rawSign, sign)
	}

	params = map[string]string{
		"userId":      "1",
		"ts":          "1531446253",
		"requestbody": `{"operate":"test","operateId":3}`,
	}
	rawSign = "75175495e7a4f1d42dbe0a88674dd14a"
	sign = usa.Sign(params)
	if sign != rawSign {
		t.Fatalf("expect sign %s, but got %s", rawSign, sign)
	}
}
