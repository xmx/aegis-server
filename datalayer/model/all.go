package model

func All() []any {
	return []any{
		Certificate{},
		GridChunk{},
		GridFile{},
		Menu{},
		OAuthConfig{},
		Oplog{},
		Role{},
		RoleMenu{},
		User{},
	}
}
