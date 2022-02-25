package api

// Category
type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
}

// User
type CreateUserRequest struct {
	Name     string `json:"name" binding="required"`
	Address  string `json:"address" binding="required"`
	Pic      string `json:"pic" binding="required"`
	Password string `json:"password" binding="required"`
	Email    string `json:"email" binding="required"`
}

type GenericUserResponse struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Pic     string `json:"pic"`
	Email   string `json:"email"`
}
