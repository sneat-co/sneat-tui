package actionstore

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/sneat-co/sneat-ai-backend/actionspec"
)

func TestSaveLoad_RoundTrips(t *testing.T) {
	s := NewStore(t.TempDir())
	d := Draft{ActionID: "act_abc12345", SpaceID: "sp1", Semantic: actionspec.Semantic{Kind: actionspec.KindBuy}}
	if err := s.Save(d); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := s.Load("act_abc12345")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.ActionID != d.ActionID || got.SpaceID != d.SpaceID || got.Semantic.Kind != actionspec.KindBuy {
		t.Fatalf("round trip mismatch: %+v", got)
	}
}

func TestLoad_NotFound(t *testing.T) {
	s := NewStore(t.TempDir())
	if _, err := s.Load("act_missing1"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestList_EmptyDirIsNotAnError(t *testing.T) {
	s := NewStore(filepath.Join(t.TempDir(), "does-not-exist"))
	got, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("got = %+v, want empty", got)
	}
}

func TestList_SortedByID(t *testing.T) {
	s := NewStore(t.TempDir())
	_ = s.Save(Draft{ActionID: "act_zzz00000"})
	_ = s.Save(Draft{ActionID: "act_aaa00000"})
	got, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 || got[0].ActionID != "act_aaa00000" || got[1].ActionID != "act_zzz00000" {
		t.Fatalf("got = %+v", got)
	}
}

func TestDelete_RemovesDraft(t *testing.T) {
	s := NewStore(t.TempDir())
	_ = s.Save(Draft{ActionID: "act_abc12345"})
	if err := s.Delete("act_abc12345"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := s.Load("act_abc12345"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound after delete", err)
	}
}

func TestDelete_MissingIsNotAnError(t *testing.T) {
	s := NewStore(t.TempDir())
	if err := s.Delete("act_missing1"); err != nil {
		t.Fatalf("Delete missing: %v", err)
	}
}

func TestDefaultDir(t *testing.T) {
	dir, err := DefaultDir(func() (string, error) { return "/cfg", nil })
	if err != nil {
		t.Fatalf("DefaultDir: %v", err)
	}
	if dir != filepath.Join("/cfg", "sneat", "actions") {
		t.Fatalf("dir = %q", dir)
	}
}

func TestDefaultDir_Error(t *testing.T) {
	if _, err := DefaultDir(func() (string, error) { return "", errors.New("boom") }); err == nil {
		t.Fatalf("expected propagated error")
	}
}
