package junkboy

import (
	"errors"
	"testing"
)

type mockAnchorRepository struct {
	AddAnchorFunc    func(a Anchor) (int, error)
	UpdateAnchorFunc func(a Anchor) error
	GetAnchorFunc    func(id int) (Anchor, error)
	GetAnchorsFunc   func() ([]Anchor, error)
	DeleteAnchorFunc func(id int) error
}

func (ar *mockAnchorRepository) AddAnchor(a Anchor) (int, error)  { return ar.AddAnchorFunc(a) }
func (ar *mockAnchorRepository) UpdateAnchor(a Anchor) error      { return ar.UpdateAnchorFunc(a) }
func (ar *mockAnchorRepository) GetAnchor(id int) (Anchor, error) { return ar.GetAnchorFunc(id) }
func (ar *mockAnchorRepository) GetAnchors() ([]Anchor, error)    { return ar.GetAnchorsFunc() }
func (ar *mockAnchorRepository) DeleteAnchor(id int) error        { return ar.DeleteAnchorFunc(id) }

var (
	testAnchor = Anchor{
		ID:  1,
		URL: "https://example.com",
	}

	testAnchors = []Anchor{
		{
			ID:  2,
			URL: "https://example.com/a",
		},
		{
			ID:  3,
			URL: "https://example.com/b",
		},
		{
			ID:  4,
			URL: "https://example.com/c",
		},
	}
)

func TestAddAnchor(t *testing.T) {
	tests := []struct {
		name        string
		errExpected bool
		method      func(a Anchor) (int, error)
	}{
		{
			name:        "Add anchor ok",
			errExpected: false,
			method: func(a Anchor) (int, error) {
				return 1, nil
			},
		},
		{
			name:        "Add anchor error",
			errExpected: true,
			method: func(a Anchor) (int, error) {
				return 0, errors.New("error adding anchor")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &mockAnchorRepository{AddAnchorFunc: tt.method}
			s := NewAnchorService(r)
			id, err := s.AddAnchor(testAnchor)
			if tt.errExpected {
				assertError(t, err)
			} else {
				assertEqual(t, 1, id)
				assertNoError(t, err)
			}
		})
	}
}

func TestUpdateAnchor(t *testing.T) {
	tests := []struct {
		name        string
		errExpected bool
		method      func(a Anchor) error
	}{
		{
			name:        "Update anchor ok",
			errExpected: false,
			method: func(a Anchor) error {
				return nil
			},
		},
		{
			name:        "Update anchor error",
			errExpected: true,
			method: func(a Anchor) error {
				return errors.New("error updating anchor")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &mockAnchorRepository{UpdateAnchorFunc: tt.method}
			s := NewAnchorService(r)
			err := s.UpdateAnchor(testAnchor)
			if tt.errExpected {
				assertError(t, err)
			} else {
				assertNoError(t, err)
			}
		})
	}
}

func TestGetAnchor(t *testing.T) {
	tests := []struct {
		name        string
		errExpected bool
		method      func(id int) (Anchor, error)
	}{
		{
			name:        "Get anchor ok",
			errExpected: false,
			method: func(id int) (Anchor, error) {
				return testAnchor, nil
			},
		},
		{
			name:        "Get anchor error",
			errExpected: true,
			method: func(id int) (Anchor, error) {
				return Anchor{}, errors.New("error getting anchor")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &mockAnchorRepository{GetAnchorFunc: tt.method}
			s := NewAnchorService(r)
			anchor, err := s.GetAnchor(1)
			if tt.errExpected {
				assertError(t, err)
			} else {
				assertEqual(t, testAnchor, anchor)
				assertNoError(t, err)
			}
		})
	}
}

func TestGetAnchors(t *testing.T) {
	tests := []struct {
		name        string
		errExpected bool
		method      func() ([]Anchor, error)
	}{
		{
			name:        "Get anchors ok",
			errExpected: false,
			method: func() ([]Anchor, error) {
				return testAnchors, nil
			},
		},
		{
			name:        "Get anchors error",
			errExpected: true,
			method: func() ([]Anchor, error) {
				return []Anchor{}, errors.New("error getting anchors")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &mockAnchorRepository{GetAnchorsFunc: tt.method}
			s := NewAnchorService(r)
			anchors, err := s.GetAnchors()
			if tt.errExpected {
				assertError(t, err)
			} else {
				assertEqual(t, len(testAnchors), len(anchors))
				assertNoError(t, err)

				for i, anchor := range anchors {
					assertEqual(t, testAnchors[i].ID, anchor.ID)
					assertEqual(t, testAnchors[i].URL, anchors[i].URL)
				}
			}
		})
	}
}

func TestDeleteAnchor(t *testing.T) {
	tests := []struct {
		name        string
		errExpected bool
		method      func(id int) error
	}{
		{
			name:        "Delete anchor ok",
			errExpected: false,
			method: func(id int) error {
				return nil
			},
		},
		{
			name:        "Delete anchor error",
			errExpected: true,
			method: func(id int) error {
				return errors.New("error deleting anchor")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &mockAnchorRepository{DeleteAnchorFunc: tt.method}
			s := NewAnchorService(r)
			err := s.DeleteAnchor(1)
			if tt.errExpected {
				assertError(t, err)
			} else {
				assertNoError(t, err)
			}
		})
	}
}
