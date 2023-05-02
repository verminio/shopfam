package shopping

type ItemService struct {
	repository Repository
}

func (s ItemService) Upsert(id *ItemId, i item) (res *ItemId, err error) {
	if id != nil {
		res = id
		err = s.repository.UpdateItem(id, i)
	} else {
		res, err = s.repository.SaveItem(i)
	}

	return
}

func NewItemService(r Repository) *ItemService {
	return &ItemService{
		repository: r,
	}
}
