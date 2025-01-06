package model

func All() []any {
	return []any{
		ConfigCertificate{},
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
