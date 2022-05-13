package junkboy

type Anchor struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type AnchorRepository interface {
	AddAnchor(a Anchor) (int, error)
	UpdateAnchor(a Anchor) error
	GetAnchor(id int) (Anchor, error)
	GetAnchors() ([]Anchor, error)
	DeleteAnchor(id int) error
}

type AnchorService struct {
	Repository AnchorRepository
}

func NewAnchorService(r AnchorRepository) *AnchorService {
	return &AnchorService{
		Repository: r,
	}
}

func (s *AnchorService) AddAnchor(a Anchor) (int, error) {
	id, err := s.Repository.AddAnchor(a)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *AnchorService) UpdateAnchor(a Anchor) error {
	err := s.Repository.UpdateAnchor(a)
	if err != nil {
		return err
	}
	return nil
}

func (s *AnchorService) GetAnchor(id int) (Anchor, error) {
	anchor, err := s.Repository.GetAnchor(id)
	if err != nil {
		return anchor, err
	}
	return anchor, nil
}

func (s *AnchorService) GetAnchors() ([]Anchor, error) {
	anchors, err := s.Repository.GetAnchors()
	if err != nil {
		return nil, err
	}
	return anchors, nil
}

func (s *AnchorService) DeleteAnchor(id int) error {
	err := s.Repository.DeleteAnchor(id)
	if err != nil {
		return err
	}
	return nil
}
