package graphqlutil

import "github.com/letjoy-club/mida-tool/clienttoken"

type GraphQLPaginator struct {
	Size *int `json:"size,omitempty"`
	Page *int `json:"page,omitempty"`
}

type Paginator struct {
	Size int
	Page int
}

func GetPager(paginator *GraphQLPaginator) Paginator {
	if paginator == nil {
		return Paginator{Size: 10, Page: 1}
	}
	page := 1
	if paginator.Page != nil {
		page = *paginator.Page
		if page <= 0 {
			page = 1
		}
	}
	size := 10
	if paginator.Size != nil {
		size = *paginator.Size
		if size <= 0 {
			size = 10
		}
	}

	return Paginator{Size: size, Page: page}
}

func GetID(token clienttoken.ClientToken, id *string) string {
	if token.IsAnonymous() || token.IsInvalid() {
		return ""
	}
	var ret string
	if token.IsUser() {
		ret = token.String()
	} else {
		if id == nil {
			return ""
		}
		ret = *id
	}
	return ret
}
