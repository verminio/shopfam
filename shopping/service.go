package shopping

type ItemService struct {
	repository Repository
}

func (s *ItemService) Upsert(id *ItemId, i item) (res *ItemId, err error) {
	if id != nil {
		res = id
		err = s.repository.UpdateItem(*id, i)
	} else {
		res, err = s.repository.SaveItem(i)
	}

	return
}

func (s *ItemService) List() (items, error) {
	return s.repository.ListItems()
}

func (s *ItemService) Delete(id ItemId) error {
	return s.repository.DeleteItem(id)
}

func NewItemService(r Repository) *ItemService {
	return &ItemService{
		repository: r,
	}
}
