package polygol

import "testing"

func expect(t testing.TB, what bool) {
	t.Helper()
	if !what {
		t.Fatal("expection failure")
	}
}

func terr(t *testing.T, err error) {
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}
